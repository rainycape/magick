package magick

// #include <stdlib.h>
// #include <string.h>
import "C"

import (
	"unsafe"
)

func allocDoublePtr(size int) *C.double {
	return (*C.double)(C.malloc(C.size_t(unsafe.Sizeof(C.double(0)) * uintptr(size))))
}

func doublePtrToFloat64Slice(p *C.double, n int) []float64 {
	s := make([]float64, n)
	C.memcpy(unsafe.Pointer(&s[0]), unsafe.Pointer(p), C.size_t(unsafe.Sizeof(C.double(0))*uintptr(n)))
	return s
}

func float64SliceToDoublePtr(v []float64) *C.double {
	sz := C.size_t(unsafe.Sizeof(C.double(0)) * uintptr(len(v)))
	p := C.malloc(sz)
	C.memcpy(p, unsafe.Pointer(&v[0]), sz)
	return (*C.double)(p)
}

func freeDoublePtr(p *C.double) {
	C.free(unsafe.Pointer(p))
}
