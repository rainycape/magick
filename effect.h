#ifndef EFFECT_H
#define EFFECT_H

#include <magick/api.h>

typedef struct {
    int order;
    double *kernel;
} ConvolveData;

typedef struct {
	double radius;
	double sigma;
	double amount;
	double threshold;
} UnsharpMaskData;

Image * convolveImage(Image *image, void *data, ExceptionInfo *ex);
Image * unsharpMaskImage(Image *image, void *data, ExceptionInfo *ex);

#endif
