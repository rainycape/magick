// +build gm

package magick

// #include <magick/api.h>
// #include "bridge.h"
//
// int
// transformImageColorspace(Image *image, void *cs, ExceptionInfo *ex) {
//  ColorspaceType *c = cs;
//  return TransformColorspace(image, *c);
// }
import "C"

func (im *Image) transformColorspace(cs Colorspace) (*Image, error) {
	return im.applyDataFunc("transforming-colorspace", C.ImageDataFunc(C.transformImageColorspace), &cs)
}
