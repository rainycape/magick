#include <magick/api.h>
#include "resize.h"

Image *
resizeImage(Image *image, void *data, ExceptionInfo *ex)
{
    ResizeData *d = data;
    return ResizeImage(image, d->columns, d->rows, d->filter, d->blur, ex);
}

Image *
sampleImage(Image *image, void *data, ExceptionInfo *ex)
{
    SizeData *d = data;
    return SampleImage(image, d->columns, d->rows, ex);
}

Image *
scaleImage(Image *image, void *data, ExceptionInfo *ex)
{
    SizeData *d = data;
    return ScaleImage(image, d->columns, d->rows, ex);
}

Image *
thumbnailImage(Image *image, void *data, ExceptionInfo *ex)
{
    SizeData *d = data;
    return ThumbnailImage(image, d->columns, d->rows, ex);
}
