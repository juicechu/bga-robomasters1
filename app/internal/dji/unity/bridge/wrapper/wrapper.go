package wrapper

import (
	"sync"
)

var (
	once     sync.Once
	instance Wrapper
)

type Wrapper interface {
	CreateUnityBridge(string, bool)
	DestroyUnityBridge()
	UnityBridgeInitialize() bool
	UnityBridgeUninitialize()
	UnitySetEventCallback(uint64, EventCallback)
	UnitySendEvent(uint64, []byte, uint64)
	UnitySendEventWithString(uint64, string, uint64)
	UnitySendEventWithNumber(uint64, uint64, uint64)

	Callback(uint64) EventCallback

	unitySetEventCallback(uint64, bool)
}
