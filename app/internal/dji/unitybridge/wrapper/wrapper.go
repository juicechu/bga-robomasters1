// +build !windows !amd64

package wrapper

type wrapper struct {}

type UnityEventCallbackFunc func()

func Instance() *wrapper {
	panic("Only windows amd64 is supported.")
}

func (w *wrapper) CreateUnityBridge(name string, debuggable bool) {
}

func (w *wrapper) DestroyUnityBridge() {
}

func (w *wrapper) UnityBridgeInitialize() bool {
	return false
}

func (w *wrapper) UnityBridgeUninitialize() {
}

func (w *wrapper) UnitySetEventCallback(e int64,
	callback UnityEventCallbackFunc) {
}

func (w *wrapper) UnitySendEvent(e uint64, info []byte, tag uint64) {
}

func (w *wrapper) UnitySendEventWithString(e uint64, info string, tag uint64) {
}

func (w *wrapper) UnitySendEventWithNumber(e, info, tag uint64) {
}

