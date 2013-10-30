// +build !debug

package magick

// #include <magick/api.h>
import "C"

import (
	"sync/atomic"
	"unsafe"
)

func refImage(image *C.Image) {
	atomic.AddInt32((*int32)(unsafe.Pointer(&image.client_data)), 1)
}

func unrefImage(image *C.Image) {
	if atomic.AddInt32((*int32)(unsafe.Pointer(&image.client_data)), -1) == 0 {
		C.DestroyImage(image)
	}
}
