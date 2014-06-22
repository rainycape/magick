#ifndef EFFECT_H
#define EFFECT_H

#include <magick/api.h>

typedef struct {
    int order;
    double *kernel;
} ConvolveData;

Image * convolveImage(Image *image, void *data, ExceptionInfo *ex);

#endif
