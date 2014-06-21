package magick

// #include <magick/api.h>
import "C"

import (
	"fmt"
	"unsafe"
)

type Storage C.StorageType

const (
	CharPixel    Storage = C.CharPixel
	ShortPixel           = C.ShortPixel
	IntegerPixel         = C.IntegerPixel
	LongPixel            = C.LongPixel
	FloatPixel           = C.FloatPixel
	DoublePixel          = C.DoublePixel
)

func (s Storage) size() int {
	switch s {
	case CharPixel:
		return 1
	case ShortPixel:
		return 2
	case IntegerPixel:
		return 4
	case LongPixel:
		return int(unsafe.Sizeof(int(0)))
	case FloatPixel:
		return 4
	case DoublePixel:
		return 8
	}
	return 0
}

// Constitute returns a new image from in-memory data. channels must be a string formed
// or the following characters, each one representing a channel, in order.
//
//  R = red
//  G = green
//  B = blue
//  A = alpha (same as Transparency)
//  O = Opacity
//  T = Transparency
//  C = cyan
//  Y = yellow
//  M = magenta
//  K = black
//  I = intensity (for grayscale)
//  P = pad, to skip over a chnnel which is intentionally ignored
//
// Creation of an alpha channel for CMYK images is currently not supported. Note that pixels
// must be either empty, resulting in an image with all pixels the zero values, or have the
// appropiate size for the image dimensions, the number of channels and the storage type.
func Constitute(width int, height int, channels string, st Storage, pixels []byte) (*Image, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	m := C.CString(channels)
	var px unsafe.Pointer
	count := width * height * len(channels) * st.size()
	if len(pixels) > 0 {
		if len(pixels) < count {
			return nil, fmt.Errorf("pixels must contain at least %d bytes, it has %d", count, len(pixels))
		}
		px = unsafe.Pointer(&pixels[0])
	} else {
		px = C.calloc(C.size_t(count), 1)
		defer C.free(px)
	}
	im := C.ConstituteImage(magickSize(uint(width)), magickSize(uint(height)), m, C.StorageType(st), px, &ex)
	C.free(unsafe.Pointer(m))
	if ex.severity != C.UndefinedException {
		return nil, exError(&ex, "allocating image")
	}
	refImage(im)
	img := &Image{
		image: im,
	}
	freeWhenDone(img)
	return img, nil
}
