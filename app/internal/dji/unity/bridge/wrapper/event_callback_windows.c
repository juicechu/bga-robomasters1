#include "event_callback_windows.h"

#include "_cgo_export.h"

extern void eventCallbackGo(void* context, GoUint64 e, GoSlice info,
		GoUint64 tag);

void event_callback(void* context, va_alist alist) {
	va_start_void(alist);
	unsigned long long event_code = va_arg_ulonglong(alist);
        void* data = va_arg_ptr(alist, void*);
        int length = va_arg_int(alist);
        unsigned long long tag = va_arg_ulonglong(alist);
	
	// Create a Go slice with the info data.
	GoSlice data_slice;
	data_slice.data = data;
	data_slice.len = length;
	data_slice.cap = length;

	eventCallbackGo(context, event_code, data_slice, tag);
}

