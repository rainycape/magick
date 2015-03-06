package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "shear.h"
import "C"

// AffineTransform returns a new image created by transforming
// the original image as dictated by the affine matrix.
func (im *Image) AffineTransform(m *AffineMatrix) (*Image, error) {
	return im.applyDataFunc("transforming (affine)", C.ImageDataFunc(C.AffineTransformImage), m.matrix())
}

// Rotate creates a new image that is a rotated copy of an existing one. Positive angles rotate counter-clockwise (right-hand rule), while negative angles rotate clockwise.
func (im *Image) Rotate(degrees float64) (*Image, error) {
	data := C.double(degrees)
	return im.applyDataFunc("rotating", C.ImageDataFunc(C.rotateImage), &data)
}
 