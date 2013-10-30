// +build debug

package magick

// #include <magick/api.h>
import "C"

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	mu         sync.Mutex
	liveImages = make(map[*C.Image]struct{})
	live       int32
)

func refImage(image *C.Image) {
	if atomic.AddInt32((*int32)(unsafe.Pointer(&image.client_data)), 1) == 1 {
		mu.Lock()
		liveImages[image] = struct{}{}
		mu.Unlock()
		atomic.AddInt32(&live, 1)
	}
}

func unrefImage(image *C.Image) {
	if atomic.AddInt32((*int32)(unsafe.Pointer(&image.client_data)), -1) == 0 {
		mu.Lock()
		delete(liveImages, image)
		mu.Unlock()
		atomic.AddInt32(&live, -1)
		C.DestroyImage(image)
	}
}
