package magick

import (
	"reflect"
	"unsafe"
)

func goBytes(p unsafe.Pointer, s int) []byte {
	var b []byte
	header := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	header.Len = s
	header.Cap = s
	header.Data = uintptr(p)
	return b
}
