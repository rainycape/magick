// +build !gm

package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "shear.h"
import "C"

// Rotate creates a new image that is a rotated copy of an existing one. Positive angles rotate counter-clockwise (right-hand rule), while negative angles rotate clockwise.
func (im *Image) Rotate(degrees float64) (*Image, error) {
	data := C.double(degrees)
	return im.applyDataFunc("rotating", C.ImageDataFunc(C.rotateImage), &data)
}
 