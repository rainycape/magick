// +build !gm

package magick

// #include <magick/api.h>
import "C"

import (
	"reflect"
	"unsafe"
)

func (im *Image) statistics() (*Statistics, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	stats := C.GetImageChannelStatistics(im.image, &ex)
	if stats != nil {
		// Don't use MagickRelinquishMemory. There's no point
		// in requiring the wand API only for this function and
		// it's only a wrapper of RelinquishMagickMemory adding
		// logging.
		defer C.RelinquishMagickMemory(unsafe.Pointer(stats))
	}
	if stats == nil || ex.severity != C.UndefinedException {
		return nil, exError(&ex, "getting statistics")
	}
	return newStatistics(stats), nil
}

func newChannelStatistics(ch C.ChannelStatistics) *ChannelStatistics {
	return &ChannelStatistics{
		Minimum:  float64(ch.minima / C.QuantumRange),
		Maximum:  float64(ch.maxima / C.QuantumRange),
		Mean:     float64(ch.mean / C.QuantumRange),
		StdDev:   float64(ch.standard_deviation / C.QuantumRange),
		Variance: float64(ch.variance / C.QuantumRange),
		Kurtosis: float64(ch.kurtosis / C.QuantumRange),
		Skewness: float64(ch.skewness / C.QuantumRange),
	}
}

func newStatistics(stats *C.ChannelStatistics) *Statistics {
	count := C.OpacityChannel + 1
	var channels []C.ChannelStatistics
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&channels)))
	sliceHeader.Cap = count
	sliceHeader.Len = count
	sliceHeader.Data = uintptr(unsafe.Pointer(stats))
	red := newChannelStatistics(channels[C.RedChannel])
	green := newChannelStatistics(channels[C.GreenChannel])
	blue := newChannelStatistics(channels[C.BlueChannel])
	opacity := newChannelStatistics(channels[C.OpacityChannel])
	return &Statistics{red, green, blue, opacity}
}
