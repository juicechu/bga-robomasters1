package bridge

import (
	"fmt"
	"log"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper"
)

// EventHandler is the required interface for types that want to listen to Unity
// Events. Implementations are required to call wg.Done() before they return.
type EventHandler interface {
	HandleEvent(event *unity.Event, data []byte, tag uint64,
		wg *sync.WaitGroup)
}

// unityBridge is a frontend to DJI's Unity Bridge code. The type is not
// exported because it is currently implemented as a singleton.
type unityBridge struct {
	m     sync.Mutex
	setup bool

	me              sync.Mutex
	eventHandlerMap map[unity.EventType]map[int]EventHandler
}

var (
	// Singleton instance.
	instance *unityBridge

	wrapperInstance wrapper.Wrapper
)

func init() {
	// Creates the singleton instance.
	instance = &unityBridge{
		sync.Mutex{},
		false,
		sync.Mutex{},
		make(map[unity.EventType]map[int]EventHandler),
	}

	wrapperInstance = wrapper.Instance()
}

// Setup creates and initializes the underlying Unity Bridge. It returns a nil
// error on success and a non-nil error on failure.
func Setup(name string, debuggable bool) error {
	instance.m.Lock()
	defer instance.m.Unlock()

	if instance.setup {
		return fmt.Errorf("bridge already setup")
	}

	// Creates the underlying Unity Bridge.
	wrapperInstance.CreateUnityBridge(name, debuggable)

	// Register the callback to all known events.
	instance.registerCallback()

	if !wrapperInstance.UnityBridgeInitialize() {
		// Something went wrong so we bail out.
		wrapperInstance.DestroyUnityBridge()
		return fmt.Errorf("bridge initialization failed")
	}

	instance.setup = true

	return nil
}

// Teardown uninitializes and destroys the underlying Unity Bridge. It returns
// a nil error on success and a non-nil error on failure.
func Teardown() error {
	instance.m.Lock()
	defer instance.m.Unlock()

	if instance.setup {
		return fmt.Errorf("bridge not setup")
	}

	// Unregister the callback to all known events.
	instance.unregisterCallback()

	wrapperInstance.UnityBridgeUninitialize()
	wrapperInstance.DestroyUnityBridge()

	instance.setup = false

	return nil
}

// IsSetup returns true if the underlying Unity Bridge support was setup and
// false otherwise.
func IsSetup() bool {
	instance.m.Lock()
	defer instance.m.Unlock()

	return instance.setup
}

// Instance returns a pointer to the unityBridge singleton.
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
	for i = 0; ; i++ {
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
	if !ok {
		return fmt.Errorf("no handlers for given event")
	}

	_, ok = handlerMap[index]
	if !ok {
		return fmt.Errorf("no handler for given event at given index")
	}

	delete(handlerMap, index)

	if len(handlerMap) == 0 {
		delete(b.eventHandlerMap, eventType)
	}

	return nil
}

// SendEvent sends a unity event through the underlying Unity Bridge. It can
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
			wrapperInstance.UnitySendEvent(event.Code(), data.([]byte), tag)
		} else {
			wrapperInstance.UnitySendEvent(event.Code(), nil, tag)
		}
	case 1:
		wrapperInstance.UnitySendEventWithString(event.Code(), data.(string), tag)
	case 2:
		wrapperInstance.UnitySendEventWithNumber(event.Code(), data.(uint64), tag)
	}

	return nil
}

func (b *unityBridge) registerCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		wrapperInstance.UnitySetEventCallback(event.Code(),
			b.unityEventCallback)
	}
}

func (b *unityBridge) unregisterCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		wrapperInstance.UnitySetEventCallback(event.Code(), nil)
	}
}

func (b *unityBridge) unityEventCallback(eventCode uint64, data []byte, tag uint64) {
	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d.\n", eventCode)
		return
	}

	eventHandlers, ok := b.eventHandlerMap[event.Type()]
	if !ok {
		log.Printf("No event handlers for %q\n",
			unity.EventTypeName(event.Type()))
		return
	}

	wg := sync.WaitGroup{}

	for _, handler := range eventHandlers {
		wg.Add(1)
		go handler.HandleEvent(event, data, tag, &wg)
	}

	wg.Wait()
}
