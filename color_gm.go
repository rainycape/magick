// +build gm

package magick

// #include <magick/api.h>
import "C"

import (
	"reflect"
	"unsafe"
)

func (im *Image) histogram() (*Histogram, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	var colors C.ulong
	hist := C.GetColorHistogram(im.image, &colors, &ex)
	if hist == nil {
		return nil, exError(&ex, "getting histogram")
	}
	defer C.free(unsafe.Pointer(hist))
	count := int(colors)
	items := make([]*HistogramItem, count)
	var hitems []C.HistogramColorPacket
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&hitems)))
	sliceHeader.Cap = count
	sliceHeader.Len = count
	sliceHeader.Data = uintptr(unsafe.Pointer(hist))
	for ii, v := range hitems {
		px := newPixel(&v.pixel)
		items[ii] = &HistogramItem{px, int(v.count)}
	}
	return &Histogram{items}, nil
}
