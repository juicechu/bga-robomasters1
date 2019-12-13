// +build !windows !amd64

package wrapper

/*
#include "unitybridge.h"
*/
import "C"

func CreateUnityBridge(name string, debuggable bool) {
	panic("Windows amd64 is the only supported platform.")
}

func DestroyUnityBridge() {
	panic("Windows amd64 is the only supported platform.")
}

func UnityBridgeInitialize() bool {
	panic("Windows amd64 is the only supported platform.")
}

func UnityBridgeUninitialize() {
	panic("Windows amd64 is the only supported platform.")
}

func UnitySetEventCallback(e uint64, callback C.UnityEventCallbackFunc) {
	panic("Windows amd64 is the only supported platform.")
}

func UnitySendEvent(e uint64, info []byte, tag uint64) {
	panic("Windows amd64 is the only supported platform.")
}

func UnitySendEventWithString(e uint64, info string, tag uint64) {
	panic("Windows amd64 is the only supported platform.")
}

func UnitySendEventWithNumber(e, info, tag uint64) {
	panic("Windows amd64 is the only supported platform.")
}
