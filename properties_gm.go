// +build gm

package magick

// #include <stdlib.h>
// #include <magick/api.h>
import "C"

import (
	"unsafe"
)

func (im *Image) properties() []string {
	prop := C.GetImageAttribute(im.image, nil)
	var props []string
	for prop != nil {
		if prop.key != nil {
			s := C.GoString(prop.key)
			props = append(props, s)
		}
		prop = (*C.ImageAttribute)(prop.next)
	}
	return props
}

func (im *Image) destroyProperties() {
	C.DestroyImageAttributes(im.image)
}

func (im *Image) property(key string) *string {
	k := C.CString(key)
	prop := C.GetImageAttribute(im.image, k)
	C.free(unsafe.Pointer(k))
	if prop != nil && prop.value != nil {
		s := C.GoString(prop.value)
		return &s
	}
	return nil
}

func (im *Image) setProperty(key string, value *C.char) bool {
	if value != nil {
		// Clear the property first, otherwise GM concatenates
		// the values
		im.setProperty(key, nil)
	}
	k := C.CString(key)
	ret := C.SetImageAttribute(im.image, k, value)
	C.free(unsafe.Pointer(k))
	return ret != 0
}
