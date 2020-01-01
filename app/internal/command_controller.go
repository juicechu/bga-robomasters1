package internal

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
	"log"
	"sync"
)

type EventHandler func(result *dji.Result)

type CommandController struct {
	*GenericController

	mslm              sync.Mutex
	startListeningMap map[dji.Key]map[int]EventHandler
}

func NewCommandController() (*CommandController, error) {
	cc := &CommandController{
		startListeningMap: make(map[dji.Key]map[int]EventHandler),
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

func (c *CommandController) StartListening(key dji.Key, eventHandler EventHandler) (int, error) {
	if key < 1 || key >= dji.KeysCount {
		return -1, fmt.Errorf("invalid key")
	}
	if eventHandler == nil {
		return -1, fmt.Errorf("eventHandler must not be nil")
	}
	if (key.AccessType() & dji.KeyAccessTypeRead) == 0 {
		return -1, fmt.Errorf("key is not readable")
	}

	c.mslm.Lock()
	defer c.mslm.Unlock()

	handlerMap, ok := c.startListeningMap[key]
	if !ok {
		handlerMap = make(map[int]EventHandler)
	}

	if len(handlerMap) == 0 {
		err := bridge.Instance().SendEvent(unity.NewEventWithSubType(
			unity.EventTypeStartListening, uint64(key.Value())))
		if err != nil {
			return -1, err
		}
	}

	var i int
	for i = 0; ; i++ {
		_, ok := handlerMap[i]
		if !ok {
			handlerMap[i] = eventHandler
			break
		}
	}

	c.startListeningMap[key] = handlerMap

	return i, nil
}

func (c *CommandController) StopListening(key dji.Key, index int) error {
	if key < 1 || key >= dji.KeysCount {
		return fmt.Errorf("invalid key")
	}
	if index < 0 {
		return fmt.Errorf("index must be non-negative")
	}

	c.mslm.Lock()
	defer c.mslm.Unlock()

	handlerMap, ok := c.startListeningMap[key]
	if !ok {
		return fmt.Errorf("no handlers for given key")
	}

	_, ok = handlerMap[index]
	if !ok {
		return fmt.Errorf("no listening handler at given index")
	}

	delete(handlerMap, index)

	if len(handlerMap) == 0 {
		delete(c.startListeningMap, key)

		err := bridge.Instance().SendEvent(unity.NewEventWithSubType(
			unity.EventTypeStopListening, uint64(key.Value())))
		if err != nil {
			return err
		}
	}

	return nil
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

func (c *CommandController) HandleEvent(event *unity.Event, info []byte,
	tag uint64, wg *sync.WaitGroup) {
	var value interface{}

	infoType := (tag >> 56) & 0xff
	switch infoType {
	case 0:
		value = string(info)
	case 1:
		value = binary.LittleEndian.Uint64(info)
	default:
		value = info
	}

	switch event.Type() {
	case unity.EventTypeStartListening:
		c.handleStartListening(value, tag)
	default:
		log.Printf("Unsupported event %s. Value:%v. Tag:%d\n",
			unity.EventTypeName(event.Type()), value, tag)
	}

	wg.Done()
}

func (c *CommandController) handleStartListening(value interface{}, tag uint64) {
	result := dji.NewResultFromJSON([]byte(value.(string)))

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
