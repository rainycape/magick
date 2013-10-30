// +build !gm

package magick

// #cgo pkg-config: MagickCore
// #cgo LDFLAGS: -lm
// #include <magick/MagickCore.h>
import "C"

import (
	"fmt"
	"unsafe"
)

var (
	backend = "ImageMagick"
)

func magickBool(flag bool) C.MagickBooleanType {
	if flag {
		return C.MagickTrue
	}
	return C.MagickFalse
}

func magickLong(v int) C.ssize_t {
	return C.ssize_t(v)
}

func magickSize(v uint) C.size_t {
	return C.size_t(v)
}

func imageToBlob(info *Info, im *Image, s *C.size_t, ex *C.ExceptionInfo) unsafe.Pointer {
	return unsafe.Pointer(C.ImagesToBlob(info.info, im.image, s, ex))
}

type notImplementedError string

func (e notImplementedError) Error() string {
	return fmt.Sprintf("magick does not implement %s when built against ImageMagick. Build magick against GraphicsMagick to use it.", string(e))
}

func cleanup() {
	C.MagickCoreTerminus()
}

func init() {
	C.MagickCoreGenesis(nil, magickBool(false))
}
