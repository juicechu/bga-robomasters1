package wrapper

/*
#cgo LDFLAGS: -L${SRCDIR} -lunitybridge

#include <stdlib.h>

#include "event_callback_windows_amd64.h"
#include "unitybridge.h"
*/
import "C"
import (
	"log"
	"unsafe"
)

type Windows struct {
	*generic
}

func Instance() Wrapper {
	once.Do(func() {
		w := &Windows{}

		w.generic = newGeneric(w)

		instance = w
	})

	return instance
}

func (w *Windows) CreateUnityBridge(name string, debuggable bool) {
	intDebuggable := C.int(0)
	if debuggable {
		intDebuggable = C.int(1)
	}

	cName := C.CString(name)

	C.CreateUnityBridge(cName, intDebuggable)

	C.free(unsafe.Pointer(cName))
}

func (w *Windows) DestroyUnityBridge() {
	C.DestroyUnityBridge()
}

func (w *Windows) UnityBridgeInitialize() bool {
	ok := C.UnityBridgeInitialize()

	return ok != 0
}

func (w *Windows) UnityBridgeUninitialize() {
	C.UnityBridgeUninitialze()
}

func (w *Windows) UnitySendEvent(e uint64, info []byte, tag uint64) {
	var infoPtr unsafe.Pointer = nil
	if len(info) != 0 {
		infoPtr = unsafe.Pointer(&info[0])
	}

	C.UnitySendEvent(C.ulonglong(e), infoPtr, C.ulonglong(tag))
}

func (w *Windows) UnitySendEventWithString(e uint64, info string, tag uint64) {
	cInfo := C.CString(info)

	C.UnitySendEventWithString(C.ulonglong(e), cInfo,
		C.ulonglong(tag))

	C.free(unsafe.Pointer(cInfo))
}

func (w *Windows) UnitySendEventWithNumber(e, info, tag uint64) {
	C.UnitySendEventWithNumber(C.ulonglong(e), C.ulonglong(info),
		C.ulonglong(tag))
}

func (w *Windows) unitySetEventCallback(eventCode uint64, add bool) {
	if add {
		C.UnitySetEventCallback(C.ulonglong(eventCode),
			C.UnityEventCallbackFunc(C.eventCallback))
	} else {
		C.UnitySetEventCallback(C.ulonglong(eventCode), nil)
	}
}

//export eventCallbackGo
func eventCallbackGo(eventCode uint64, info []byte, tag uint64) {
	callback := instance.Callback(eventCode)
	if callback == nil {
		log.Printf("No callback for event code %d.\n", eventCode)
		return
	}

	callback(eventCode, info, tag)
}
