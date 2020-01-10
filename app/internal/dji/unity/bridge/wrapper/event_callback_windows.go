package wrapper

import "C"
import (
	"fmt"
	"unsafe"

	"github.com/mattn/go-pointer"
)

//export eventCallbackGo
func eventCallbackGo(context unsafe.Pointer, eventCode uint64, data []byte,
	tag uint64) {
	fmt.Printf("eventCallbackGo(context=%p, eventCode=%d, data=%v, "+
		"tag=%d)\n", context, eventCode, data, tag)
	cb := pointer.Restore(context).(EventCallback)
	cb(eventCode, data, tag)
}
