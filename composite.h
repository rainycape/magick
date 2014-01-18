#ifndef COMPOSITE_H
#define COMPOSITE_H

#include <magick/api.h>

typedef struct {
    int composite;
    Image *draw;
    int x;
    int y;
} CompositeData;

Image * compositeImage(Image *canvas, void *data, ExceptionInfo *ex);

#endif
