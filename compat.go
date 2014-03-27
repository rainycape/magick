package magick

// #include <magick/api.h>
// extern int copy_rgba_pixels(const Image *image, unsigned char *dest);
import "C"

import (
	"errors"
	"image"
	"unsafe"
)

var (
	errCantCopyPixels = errors.New("can't copy pixels")
)

// GoImage returns the image converted to an RGBA image.Image.
// Note that since the in-memory representation is different, the
// conversion could introduce small errors due to rounding
// in semitransparent pixels, because Go uses premultiplied
// alpha while GM/IM does not.
func (im *Image) GoImage() (image.Image, error) {
	w := im.Width()
	h := im.Height()
	res := image.NewRGBA(image.Rect(0, 0, w, h))
	dest := (*C.uchar)(unsafe.Pointer(&res.Pix[0]))
	if C.copy_rgba_pixels(im.image, dest) == 0 {
		return nil, errCantCopyPixels
	}
	return res, nil
}
