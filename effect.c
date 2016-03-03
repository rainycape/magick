#include "effect.h"

Image *
convolveImage(Image *image, void *data, ExceptionInfo *ex) {
    ConvolveData *d = data;
    return ConvolveImage(image, d->order, d->kernel, ex);
}

Image *
unsharpMaskImage(Image *image, void *data, ExceptionInfo *ex) {
    UnsharpMaskData *d = data;
    return UnsharpMaskImage(image, d->radius, d->sigma, d->amount, d->threshold ,ex);
}
