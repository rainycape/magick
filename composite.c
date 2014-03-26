#include "composite.h"

Image * compositeImage(Image *canvas, void *data, ExceptionInfo *ex)
{
    CompositeData *d = data;
    if (!CompositeImage(canvas, d->composite, d->draw, d->x, d->y)) {
        ex->severity = DrawError;
        return NULL;
    }
    return canvas;
}
