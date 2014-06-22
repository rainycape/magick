package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "effect.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// Convolve applies the given convolution kernel to the image. The
// order parameter must be a non-negative odd number, while the kernel
// parameter must have either order*order elements or just one element
// (in the latter case, it's interpreted as having all elements set to
// that first value).
func (im *Image) Convolve(order int, kernel []float64) (*Image, error) {
	count := order * order
	if len(kernel) < count {
		if len(kernel) != 1 {
			return nil, fmt.Errorf("kernel for order %d must have %d or %d elements, not %d", order, count, 1, len(kernel))
		}
		newKernel := make([]float64, count)
		for ii := range newKernel {
			newKernel[ii] = kernel[0]
		}
		kernel = newKernel
	}
	var data C.ConvolveData
	data.order = C.int(order)
	data.kernel = (*C.double)(unsafe.Pointer(&kernel[0]))
	return im.applyDataFunc("convolving", C.ImageDataFunc(C.convolveImage), &data)
}
