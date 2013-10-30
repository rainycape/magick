// +build gm

package magick

// #cgo CFLAGS: -I/usr/include/GraphicsMagick
// #cgo LDFLAGS: -lGraphicsMagick -lm
// #include <magick/api.h>
import "C"

import (
	"unsafe"
)

var (
	backend = "GraphicsMagick"
)

func magickBool(flag bool) C.int {
	if flag {
		return 1
	}
	return 0
}

func magickLong(v int) C.long {
	return C.long(v)
}

func magickSize(v uint) C.ulong {
	return C.ulong(v)
}

func imageToBlob(info *Info, im *Image, s *C.size_t, ex *C.ExceptionInfo) unsafe.Pointer {
	return C.ImageToBlob(info.info, im.image, s, ex)
}

func cleanup() {
	C.DestroyMagick()
}

func init() {
	C.InitializeMagick(nil)
}
