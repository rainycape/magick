// +build !gm

package magick

// #cgo pkg-config: MagickCore
// #cgo CFLAGS: -D_MAGICK_USES_IM
// #cgo LDFLAGS: -lm
// #include <magick/MagickCore.h>
import "C"

import (
	"fmt"
	"reflect"
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

func magickUint(v uint) C.size_t {
	return magickSize(v)
}

func imageToBlob(info *Info, im *Image, s *C.size_t, ex *C.ExceptionInfo) unsafe.Pointer {
	if info.info.quality != 0 {
		quality := im.image.quality
		im.image.quality = info.info.quality
		defer func() {
			im.image.quality = quality
		}()
	}
	return unsafe.Pointer(C.ImagesToBlob(info.info, im.image, s, ex))
}

func freeMagickMemory(p unsafe.Pointer) {
	C.RelinquishMagickMemory(p)
}

func supportedFormats(ex *C.ExceptionInfo) ([]*C.MagickInfo, unsafe.Pointer) {
	var count C.size_t
	patt := C.CString("")
	info := C.GetMagickInfoList(patt, &count, ex)
	if info == nil {
		return nil, nil
	}
	var infos []*C.MagickInfo
	header := (*reflect.SliceHeader)(unsafe.Pointer(&infos))
	header.Len = int(count)
	header.Cap = header.Len
	p := unsafe.Pointer(info)
	header.Data = uintptr(p)
	return infos, p
}

type notImplementedError string

func (e notImplementedError) Error() string {
	return fmt.Sprintf("magick does not implement %s when built against ImageMagick. Build magick against GraphicsMagick to use it.", string(e))
}
