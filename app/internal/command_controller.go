package internal

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
	"git.bug-br.org.br/bga/robomasters1/app/internal/support/callbacks"
)

type EventHandler func(result *dji.Result)

type CommandController struct {
	*GenericController

	starListeningCallbacks *callbacks.Callbacks
}

func NewCommandController() (*CommandController, error) {
	startListeningCallbacks := callbacks.New(
		"CommandController/StartListening",
		func() error {
			return bridge.Instance().SendEvent(
				unity.NewEventWithSubType(
					unity.EventTypeStartListening,
					uint64(key.Value())))
		},
		func() error {
			return bridge.Instance().SendEvent(
				unity.NewEventWithSubType(
					unity.EventTypeStopListening,
					uint64(key.Value())))
		},
	)

	cc := &CommandController{
		startListeningCallbacks: startListeningCallbacks,
	}

	cc.GenericController = NewGenericController(cc)

	var err error
	err = cc.StartControllingEvent(unity.EventTypeGetValue)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(unity.EventTypeSetValue)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(unity.EventTypePerformAction)
	if err != nil {
		return nil, err
	}
	err = cc.StartControllingEvent(unity.EventTypeStartListening)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (c *CommandController) StartListening(key dji.Key,
	eventHandler EventHandler) (uint64, error) {
	if key < 1 || key >= dji.KeysCount {
		return 0, fmt.Errorf("invalid key")
	}
	if eventHandler == nil {
		return 0, fmt.Errorf("eventHandler must not be nil")
	}
	if (key.AccessType() & dji.KeyAccessTypeRead) == 0 {
		return 0, fmt.Errorf("key is not readable")
	}

	return c.startListeningCallbacks.AddContinuous(callbacks.Key(key),
		eventHandler)
}

func (c *CommandController) StopListening(key dji.Key, tag uint64) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}

	return c.startListeningCallbacks.Remove(callbacks.Key(key),
		callbacks.Tag(tag))
}

func (c *CommandController) PerformAction(key dji.Key, param interface{},
	eventHandler EventHandler) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}
	if (key.AccessType() & dji.KeyAccessTypeAction) == 0 {
		return fmt.Errorf("key can not be acted upon")
	}

	if eventHandler != nil {
		// TODO(bga): Fix this.
		panic("No event handler support in PerformAction.")
	}

	var data []byte
	if param != nil {
		var err error
		data, err = json.Marshal(param)
		if err != nil {
			return err
		}
	}

	bridge.Instance().SendEvent(unity.NewEventWithSubType(
		unity.EventTypePerformAction, uint64(key.Value())), data,
		uint64(key.Value()))

	return nil
}

func (c *CommandController) HandleEvent(event *unity.Event, data []byte,
	tag uint64, wg *sync.WaitGroup) {
	var value interface{}

	// TODO(bga): Apparently, the unity bridge reserves the upper 8 bits
	//  for reporting back type information. Double check this.
	infoType := (tag >> 56) & 0xff
	switch infoType {
	case 0:
		value = string(info)
	case 1:
		value = binary.LittleEndian.Uint64(info)
	default:
		// Apparently only string and uint64 types are supported
		// currently.
		panic(fmt.Sprintf("Unexpected data type: %d.\n", infoType))
	}

	// See above.
	adjustedTag := tag & 0xffffffffffffff

	switch event.Type() {
	case unity.EventTypeStartListening:
		c.handleStartListening(event.Subtype(), value, adjustedTag)
	default:
		log.Printf("Unsupported event %s. Value:%v. Tag:%d\n",
			unity.EventTypeName(event.Type()), value, tag)
	}

	wg.Done()
}

func (c *CommandController) handleStartListening(value interface{},
	tag uint64) {
	stringValue, ok := value.(string)
	if !ok {
		panic("unexpected non-string value")
	}

	result, ok := dji.NewResultFromJSON([]byte(stringValue))

	c.startListeningCallbacks.Callback(

	c.mslm.Lock()
	defer c.mslm.Unlock()

	for _, handlerMap := range c.startListeningMap {
		for _, handler := range handlerMap {
			go handler(result)
		}
	}
}

func (c *CommandController) Teardown() error {
	var err error
	err = c.StopControllingEvent(unity.EventTypeGetValue)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(unity.EventTypeSetValue)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(unity.EventTypePerformAction)
	if err != nil {
		return err
	}
	err = c.StopControllingEvent(unity.EventTypeStartListening)
	if err != nil {
		return err
	}

	return nil
}
