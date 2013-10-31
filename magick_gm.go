// +build gm

package magick

// #cgo CFLAGS: -I/usr/include/GraphicsMagick
// #cgo LDFLAGS: -lGraphicsMagick -lm
// #include <magick/api.h>
import "C"

import (
	"fmt"
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

type notImplementedError string

func (e notImplementedError) Error() string {
	return fmt.Sprintf("magick does not implement %s when built against GraphicsMagick. Build magick against ImageMagick to use it.", string(e))
}

func cleanup() {
	C.DestroyMagick()
}

func init() {
	C.InitializeMagick(nil)
}
