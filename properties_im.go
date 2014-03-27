// +build !gm

package magick

// #include <stdlib.h>
// #include <magick/api.h>
import "C"

import (
	"unsafe"
)

func (im *Image) properties() []string {
	C.ResetImagePropertyIterator(im.image)
	var props []string
	for {
		prop := C.GetNextImageProperty(im.image)
		if prop == nil {
			break
		}
		props = append(props, C.GoString(prop))
	}
	C.ResetImagePropertyIterator(im.image)
	return props
}

func (im *Image) destroyProperties() {
	C.DestroyImageProperties(im.image)
}

func (im *Image) property(key string) *string {
	k := C.CString(key)
	prop := C.GetImageProperty(im.image, k)
	C.free(unsafe.Pointer(k))
	if prop != nil {
		s := C.GoString(prop)
		return &s
	}
	return nil
}

func (im *Image) setProperty(key string, value *C.char) bool {
	k := C.CString(key)
	var ret bool
	if value == nil {
		ret = C.DeleteImageProperty(im.image, k) != 0
	} else {
		C.DeleteImageProperty(im.image, k)
		ret = C.SetImageProperty(im.image, k, value) != 0
	}
	C.free(unsafe.Pointer(k))
	return ret
}
