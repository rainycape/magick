#ifndef TRANSFORM_H
#define TRANSFORM_H

Image *flattenImages(const Image *image, ExceptionInfo *exception);
Image *mosaicImages(const Image *image, ExceptionInfo *exception);
RectangleInfo *calculate_max_entropy_rect(double ratio, const Image *image, double imWidth, double imHeight, double imRatio);


#endif
