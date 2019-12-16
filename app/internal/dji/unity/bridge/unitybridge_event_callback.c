#include "_cgo_export.h"
#include "unitybridge_event_callback.h"

extern void unityEventCallbackGo(GoUint64 e, GoSlice info, GoUint64 tag);

void UnityEventCallback(unsigned long long e, void* info, int length,
		unsigned long long tag) {
	// Create a Go slice with the info data.
	GoSlice info_slice;
	info_slice.data = info;
	info_slice.len = length;
	info_slice.cap = length;

	unityEventCallbackGo(e, info_slice, tag);
}

