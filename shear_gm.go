// +build gm

package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "shear.h"
import "C"

// Rotate creates a new image that is a rotated copy of an existing one. Positive angles rotate counter-clockwise (right-hand rule), while negative angles rotate clockwise.
func (im *Image) Rotate(degrees float64) (*Image, error) {
	data := C.double(degrees)
	frames := im.NFrames()
	newImage, err := im.First().applyDataFunc("rotating", C.ImageDataFunc(C.rotateImage), &data)
	if err != nil {
		return nil, err
	}
	for i := 1; i < frames; i++ {
		frame, err := im.Frame(i)
		if err != nil {
			return nil, err
		}
		newFrame, err := frame.applyDataFunc("rotating", C.ImageDataFunc(C.rotateImage), &data)
		if err != nil {
			return nil, err
		}
		err = newImage.Append(newFrame)
		if err != nil {
			return nil, err
		}
	}
	return newImage.Coalesce()
}
 