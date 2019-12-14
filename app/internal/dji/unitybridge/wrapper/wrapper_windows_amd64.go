package wrapper

/*
#cgo LDFLAGS: -L${SRCDIR} -lunitybridge

#include <stdlib.h>

#include "unitybridge.h"
*/
import "C"
import (
	"unsafe"
)

type UnityEventCallbackFunc C.UnityEventCallbackFunc

func CreateUnityBridge(name string, debuggable bool) {
	intDebuggable := C.int(0)
	if debuggable {
		intDebuggable = C.int(1)
	}

	cName := C.CString(name)

	C.CreateUnityBridge(cName, intDebuggable)

	C.free(unsafe.Pointer(cName))
}

func DestroyUnityBridge() {
	C.DestroyUnityBridge()
}

func UnityBridgeInitialize() bool {
	ok := C.UnityBridgeInitialize()

	return ok != 0
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
	cInfo := C.CString(info)

	C.UnitySendEventWithString(C.ulonglong(e), cInfo,
		C.ulonglong(tag))

	C.free(unsafe.Pointer(cInfo))
}

func UnitySendEventWithNumber(e, info, tag uint64) {
	C.UnitySendEventWithNumber(C.ulonglong(e), C.ulonglong(info),
		C.ulonglong(tag))
}
