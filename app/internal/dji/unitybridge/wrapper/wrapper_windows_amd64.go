package wrapper

/*
#cgo LDFLAGS: -L${SRCDIR} -lunitybridge

#include "unitybridge.h"
*/
import "C"
import (
	"syscall"
	"unsafe"
)

type UnityEventCallbackFunc C.UnityEventCallbackFunc

func CreateUnityBridge(name string, debuggable bool) {
	intDebuggable := C.int(0)
	if debuggable {
		intDebuggable = C.int(1)
	}

	utf16PtrName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		panic(err)
	}

	C.CreateUnityBridge((*C.ushort)(utf16PtrName), intDebuggable)
}

func DestroyUnityBridge() {
	C.DestroyUnityBridge()
}

func UnityBridgeInitialize() bool {
	ok := C.UnityBridgeInitialize()

	return ok == 1
}

func UnityBridgeUninitialize() {
	C.UnityBridgeUninitialze()
}

func UnitySetEventCallback(e uint64, callback UnityEventCallbackFunc) {
	C.UnitySetEventCallback(C.ulonglong(e),
		C.UnityEventCallbackFunc(callback))
}

func UnitySendEvent(e uint64, info []byte, tag uint64) {
	var infoPtr unsafe.Pointer = nil
	if info != nil {
		infoPtr = unsafe.Pointer(&info[0])
	}

	C.UnitySendEvent(C.ulonglong(e), infoPtr, C.ulonglong(tag))
}

func UnitySendEventWithString(e uint64, info string, tag uint64) {
	utf16PtrInfo, err := syscall.UTF16PtrFromString(info)
	if err != nil {
		panic(err)
	}

	C.UnitySendEventWithString(C.ulonglong(e), (*C.ushort)(utf16PtrInfo),
		C.ulonglong(tag))
}

func UnitySendEventWithNumber(e, info, tag uint64) {
	C.UnitySendEventWithNumber(C.ulonglong(e), C.ulonglong(info),
		C.ulonglong(tag))
}
