package unitybridge

/*
#include "unitybridge.h"
*/
import "C"
import (
	"fmt"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unitybridge/wrapper"
)

type UnityBridge struct {
}

func New(name string, debuggable bool) (*UnityBridge, error) {
	w := wrapper.Instance()
        w.CreateUnityBridge(name, debuggable);
        eventTypes := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 100, 101, 200, 300,
                301, 302, 303, 304, 305,306, 500}
        for _, eventType := range eventTypes {
                w.UnitySetEventCallback(eventType << 32,
			wrapper.UnityEventCallbackFunc(C.UnityEventCallback))
        }
        ok := w.UnityBridgeInitialize()
        if !ok {
                w.DestroyUnityBridge()
                return nil, fmt.Errorf("failed initializing unity bridge")
        }

	return &UnityBridge{}, nil
}

func (b *UnityBridge)SendEvent(params ...interface{}) error {
        if len(params) < 1 || len(params) > 3 {
                return fmt.Errorf("1, 2 or 3 parameters are required")
        }

        event, ok := params[0].(uint64)
        if !ok {
                return fmt.Errorf("event (first) parameter must be uint64")
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

	w := wrapper.Instance()
        switch dataType {
	case 0:
		w.UnitySendEvent(event, data.([]byte), tag)
	case 1:
		w.UnitySendEventWithString(event, data.(string), tag)
	case 2:
		w.UnitySendEventWithNumber(event, data.(uint64), tag)
	}

        return nil
}

//export UnityEventCallbackGo
func UnityEventCallbackGo(e uint64, info []byte, tag uint64) {
	fmt.Printf("Callback: e=%d, info=%#+v, tag=%d\n", e, info, tag)
}

