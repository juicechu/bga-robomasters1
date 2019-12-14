package unitybridge

/*
#include "unitybridge_event_callback.h"
*/
import "C"
import (
	"fmt"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unitybridge/wrapper"
)

type UnityBridge struct {
}

func New(name string, debuggable bool) (*UnityBridge, error) {
	wrapper.CreateUnityBridge(name, debuggable)

	ub := &UnityBridge{}

	ub.RegisterCallback()

	ok := wrapper.UnityBridgeInitialize()
	if !ok {
		ub.UnregisterCallback()
		wrapper.DestroyUnityBridge()
		return nil, fmt.Errorf("failed initializing unity bridge")
	}

	return ub, nil
}

func (b *UnityBridge) SendEvent(params ...interface{}) error {
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

func (b *UnityBridge) RegisterCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		wrapper.UnitySetEventCallback(event.Code(),
			wrapper.UnityEventCallbackFunc(C.UnityEventCallback))
	}
}

func (b *UnityBridge) UnregisterCallback() {
	for _, eventType := range unity.AllEventTypes() {
		event := unity.NewEvent(eventType)
		wrapper.UnitySetEventCallback(event.Code(), nil)
	}
}

//export UnityEventCallbackGo
func UnityEventCallbackGo(e uint64, info []byte, tag uint64) {
	fmt.Printf("Callback: e=%d, info=%#+v, tag=%d\n", e, info, tag)
}
