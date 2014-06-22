// +build gm

package magick

// #include <magick/api.h>
// #include "bridge.h"
//
// Image *
// transformImageColorspace(Image *image, void *cs, ExceptionInfo *ex) {
//  ColorspaceType *c = cs;
//  if (!TransformColorspace(image, *c)) {
//	return NULL;
//  }
//  return image;
// }
import "C"

func (im *Image) transformColorspace(cs Colorspace) (*Image, error) {
	return im.applyDataFunc("transforming-colorspace", C.ImageDataFunc(C.transformImageColorspace), &cs)
}
