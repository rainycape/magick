#ifndef RESIZE_H
#define RESIZE_H

struct ResizeData {
    unsigned long columns;
    unsigned long rows;
    FilterTypes filter;
    double blur;
};

typedef struct ResizeData ResizeData;

struct SizeData {
    unsigned long columns;
    unsigned long rows;
};

typedef struct SizeData SizeData;

Image * resizeImage(Image *image, void *data, ExceptionInfo *ex);
Image * sampleImage(Image *image, void *data, ExceptionInfo *ex);
Image * scaleImage(Image *image, void *data, ExceptionInfo *ex);
Image * thumbnailImage(Image *image, void *data, ExceptionInfo *ex);

#endif
