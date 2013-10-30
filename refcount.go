package magick

// #include <magick/api.h>
import "C"

import (
	"sync/atomic"
	"unsafe"
)

func refImages(image *C.Image) {
	for image.previous != nil {
		image = (*C.Image)(image.previous)
	}
	for cur := image; cur != nil; cur = (*C.Image)(cur.next) {
		refImage(cur)
	}
}

func unrefImages(image *C.Image) {
	for image.previous != nil {
		image = (*C.Image)(image.previous)
	}
	var next *C.Image
	for cur := image; cur != nil; cur = next {
		next = (*C.Image)(cur.next)
		unrefImage(cur)
	}
}

func refCount(image *C.Image) int {
	return int(atomic.LoadInt32((*int32)(unsafe.Pointer(&image.client_data))))
}
