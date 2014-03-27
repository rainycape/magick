// +build gm

package magick

// #cgo pkg-config: GraphicsMagick
// #cgo CFLAGS: -D_MAGICK_USES_GM
// #cgo LDFLAGS: -lm
// #include <magick/api.h>
import "C"

import (
	"fmt"
	"reflect"
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

func freeMagickMemory(p unsafe.Pointer) {
	C.MagickFree(p)
}

func supportedFormats(ex *C.ExceptionInfo) ([]*C.MagickInfo, unsafe.Pointer) {
	info := C.GetMagickInfoArray(ex)
	if info == nil {
		return nil, nil
	}
	var infos []*C.MagickInfo
	header := (*reflect.SliceHeader)(unsafe.Pointer(&infos))
	header.Len = 10000
	header.Cap = header.Len
	p := unsafe.Pointer(info)
	header.Data = uintptr(p)
	return infos, p
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
