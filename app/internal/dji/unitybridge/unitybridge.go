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

//export UnityEventCallbackGo
func UnityEventCallbackGo(e uint64, info []byte, tag uint64) {
	fmt.Printf("Callback: e=%d, info=%#+v, tag=%d\n", e, info, tag)
}

