package bridge

/*
#include "unitybridge_event_callback.h"
*/
import "C"
import (
	"fmt"
	"log"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/internal/wrapper"
)

type EventHandler interface{
	HandleEvent(event *unity.Event, info []byte, tag uint64)
}

type unityBridge struct {
	m sync.Mutex
	setup bool

	me sync.Mutex
	eventHandlerMap map[unity.EventType]map[int]EventHandler
}

var (
	instance *unityBridge
)

func init() {
	instance = &unityBridge{
		sync.Mutex{},
		false,
		sync.Mutex{},
		make(map[unity.EventType]map[int]EventHandler),
	}
}

// Setup creates and initializes the unity bridge. It returns a nil error on
// success and a non-nil error on failure.
func Setup(name string, debuggable bool) error {
	instance.m.Lock()
	defer instance.m.Unlock()

	if instance.setup {
		return fmt.Errorf("unity bridge manager already setup")
	}

	wrapper.CreateUnityBridge(name, debuggable)

	instance.registerCallback()

	if !wrapper.UnityBridgeInitialize() {
		wrapper.DestroyUnityBridge()
		return fmt.Errorf("unity bridge initialization failed")
	}

	instance.setup = true

	return nil
}

// Teardown uninitializes and destroys the unity bridge. It returns a nil error
// on success and a non-nil error on failure.
func Teardown() error {
	instance.m.Lock()
	defer instance.m.Unlock()

	if instance.setup {
		return fmt.Errorf("unity bridge manager not setup")
	}

	instance.unregisterCallback()

	wrapper.UnityBridgeUninitialize()
	wrapper.DestroyUnityBridge()

	instance.setup = false

	return nil
}

// IsSetup returns true if the unity bridge support was setup and false
// otherwise.
func IsSetup() bool {
	return instance.setup
}

// Instance returns a pointer to to unity bridge.
func Instance() *unityBridge {
	return instance
}

// AddEventHandler adds an event handler for the given event.
func (b *unityBridge) AddEventHandler(eventType unity.EventType,
		eventHandler EventHandler) (int, error) {
	if !unity.IsValidEventType(eventType) {
		return -1, fmt.Errorf("invalid event type")
	}
	if eventHandler == nil {
		return -1, fmt.Errorf("eventHandler must not be nil")
	}

	b.me.Lock()
	defer b.me.Unlock()

	handlerMap, ok := b.eventHandlerMap[eventType]
	if !ok {
		handlerMap = make(map[int]EventHandler)
	}

	var i int
	for i = 0;; i++ {
		_, ok := handlerMap[i]
		if !ok {
			handlerMap[i] = eventHandler
			break
		}
	}

	b.eventHandlerMap[eventType] = handlerMap

	return i, nil
}

// RemoveEventHandler removes the event handler at the given index for the
// given event.
func (b *unityBridge) RemoveEventHandler(eventType unity.EventType, index int) error {
	if !unity.IsValidEventType(eventType) {
		return fmt.Errorf("invalid event type")
	}
	if index < 0 {
		return fmt.Errorf("index must be non-negative")
	}

	b.me.Lock()
	defer b.me.Unlock()

	handlerMap, ok := b.eventHandlerMap[eventType]
	if ! ok {
		return fmt.Errorf("no handlers for given event")
	}

	_, ok = handlerMap[index]
	if !ok {
		return fmt.Errorf("no handler for given event at given index")
	}

	delete(handlerMap, index)

	return nil
}

// SendEvent sends a unity event through the underlying unity bridge. It can
// accept one, two or three parameters. The first one is the event itself and
// must be a *unity.Event. The second one is the data to send associated with
// the event and can be a []byte, a string or a uint64. The third one is the
// tag number associated with the event (which is used to disambiguate events)
// and must be a uint64.
func (b *unityBridge) SendEvent(params ...interface{}) error {
	if len(params) < 1 || len(params) > 3 {
		return fmt.Errorf("1, 2 or 3 parameters are required")
	}

	event, ok := params[0].(*unity.Event)
	if !ok {
		return fmt.Errorf("event (first) parameter must be a *unity.Event")
	}

	dataType := 0
	var data interface{} = nil
	if len(params) > 1 {
		switch params[1].(type) {
		case []byte:
			// Do nothing.
		case string:
			dataType = 1
		case uint64:
			dataType = 2
		default:
			return fmt.Errorf("data (second) parameter must be " +
				"[]byte, string or uint64")
		}
		data = params[1]
	}

	var tag uint64 = 0
	if len(params) > 2 {
		tag, ok = params[2].(uint64)
		if !ok {
			return fmt.Errorf("tag (third) parameter must be uint64")
		}
	}

	switch dataType {
	case 0:
		if data != nil {
			wrapper.UnitySendEvent(event.Code(), data.([]byte), tag)
		} else {
			wrapper.UnitySendEvent(event.Code(), nil, tag)
		}
	case 1:
		wrapper.UnitySendEventWithString(event.Code(), data.(string), tag)
	case 2:
		wrapper.UnitySendEventWithNumber(event.Code(), data.(uint64), tag)
	}

	return nil
}

func (b *unityBridge) registerCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		wrapper.UnitySetEventCallback(event.Code(),
			wrapper.UnityEventCallbackFunc(C.UnityEventCallback))
	}
}

func (b *unityBridge) unregisterCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		wrapper.UnitySetEventCallback(event.Code(), nil)
	}
}

//export unityEventCallbackGo
func unityEventCallbackGo(eventCode uint64, info []byte, tag uint64) {
	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d.\n", eventCode)
		return
	}

	eventHandlers, ok := instance.eventHandlerMap[event.Type()]
	if !ok {
		log.Printf("No event handlers for %q\n",
			unity.EventTypeName(event.Type()))
	}

	for _, handler := range eventHandlers {
		go handler.HandleEvent(event, info, tag)
	}
}
