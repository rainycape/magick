// +build !gm

#include <magick/api.h>
#include "shear.h"

Image *
rotateImage(Image *image, void *data, ExceptionInfo *ex)
{
	double *d = data;
	return ShearRotateImage(image, *d, ex);
}
