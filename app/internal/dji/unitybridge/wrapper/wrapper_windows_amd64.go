package wrapper

/*
#include <stdlib.h>

typedef void (*UnityEventCallbackFunc)(unsigned long long e, void* info,
	unsigned long long tag);
*/
import "C"
import (
	"sync"
	"syscall"
	"unsafe"
)

type wrapper struct {
	unityBridgeDLL *syscall.DLL
}

type unityEventCallbackFunc C.UnityEventCallbackFunc

var (
	m         sync.Mutex = sync.Mutex{}
	singleton *wrapper   = nil
)

var dllFunctionMap = map[string]*syscall.Proc{
	"CreateUnityBridge":        nil,
	"DestroyUnityBridge":       nil,
	"UnityBridgeInitialize":    nil,
	"UnityBridgeUninitialze":   nil, // Typo is present in the DLL.
	"UnitySetEventCallback":    nil,
	"UnitySendEvent":           nil,
	"UnitySendEventWithString": nil,
	"UnitySendEventWithNumber": nil,
}

func Instance() *wrapper {
	m.Lock()
	defer m.Unlock()

	if singleton == nil {
		unityBridgeDLL, err := syscall.LoadDLL("unitybridge.dll")
		if err != nil {
			panic(err)
		}

		for k := range dllFunctionMap {
			proc, err := unityBridgeDLL.FindProc(k)
			if err != nil {
				panic(err)
			}

			dllFunctionMap[k] = proc
		}

		singleton = &wrapper{
			unityBridgeDLL,
		}
	}

	return singleton
}

func (w *wrapper) CreateUnityBridge(name string, debuggable bool) {
	intDebuggable := int(0)
	if debuggable {
		intDebuggable = int(1)
	}

	utf16PtrName, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		panic(err)
	}

	dllFunctionMap["CreateUnityBridge"].Call(
		uintptr(unsafe.Pointer(utf16PtrName)), uintptr(intDebuggable))
}

func (w *wrapper) DestroyUnityBridge() {
	dllFunctionMap["DestroyUnityBridge"].Call()
}

func (w *wrapper) UnityBridgeInitialize() bool {
	r1, _, _ := dllFunctionMap["UnityBridgeInitialize"].Call()

	return r1 != 0
}

func (w *wrapper) UnityBridgeUninitialize() {
	dllFunctionMap["UnityBridgeUninitialze"].Call()
}

func (w *wrapper) UnitySetEventCallback(e int64,
	callback UnityEventCallbackFunc) {
	dllFunctionMap["UnitySetEventCallback"].Call(uintptr(e),
		uintptr(unsafe.Pointer(callback)))
}

func (w *wrapper) UnitySendEvent(e uint64, info []byte, tag uint64) {
	var infoPtr unsafe.Pointer = nil
	if info != nil {
		infoPtr = unsafe.Pointer(&info[0])
	}

	dllFunctionMap["UnitySendEvent"].Call(uintptr(e), uintptr(infoPtr),
		uintptr(tag))
}

func (w *wrapper) UnitySendEventWithString(e uint64, info string, tag uint64) {
	infoPtr := C.CString(info)
	defer C.free(unsafe.Pointer(infoPtr))

	dllFunctionMap["UnitySendEventWithString"].Call(uintptr(e),
		uintptr(unsafe.Pointer(infoPtr)), uintptr(tag))
}

func (w *wrapper) UnitySendEventWithNumber(e, info, tag uint64) {
	dllFunctionMap["UnitySendEventWithNumber"].Call(uintptr(e),
		uintptr(info), uintptr(tag))
}

