#include "_cgo_export.h"
#include "event_callback_windows_amd64.h"

extern void eventCallbackGo(GoUint64 e, GoSlice info, GoUint64 tag);

void eventCallback(unsigned long long e, void* info, int length,
		unsigned long long tag) {
	// Create a Go slice with the info data.
	GoSlice info_slice;
	info_slice.data = info;
	info_slice.len = length;
	info_slice.cap = length;

	eventCallbackGo(e, info_slice, tag);
}

