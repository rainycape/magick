#ifndef BRIDGE_H
#define BRIDGE_H

typedef Image* (*ImageFunc) (Image *, ExceptionInfo *);
typedef Image* (*ImageDataFunc) (Image *, void *, ExceptionInfo *);

Image * bridge_image_func(ImageFunc f, Image *im, ExceptionInfo *ex);
Image * bridge_image_data_func(ImageDataFunc f, Image *im, void *data, ExceptionInfo *ex);

Image * apply_image_func(ImageFunc f, Image *image, void *parent, int is_coalesced, ExceptionInfo *ex);
Image * apply_image_data_func(ImageDataFunc f, Image *image, void *data, void *parent, int is_coalesced, ExceptionInfo *ex);

#endif
