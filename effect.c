#include "effect.h"

Image *
convolveImage(Image *image, void *data, ExceptionInfo *ex) {
    ConvolveData *d = data;
    return ConvolveImage(image, d->order, d->kernel, ex);
}
