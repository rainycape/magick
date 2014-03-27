#include <math.h>

#include <magick/api.h>

#include "quantum.h"
#include "macros.h"

int
copy_rgba_pixels(const Image *image, unsigned char *dest)
{
    register long y;
    register long x;
    register const PixelPacket *p;
    ExceptionInfo ex;
    unsigned int width = image->columns;
    unsigned int height = image->rows;
    for(y = 0; y < height; ++y) {
        p = ACQUIRE_IMAGE_PIXELS(image, 0, y, width, 1, &ex);
        if (!p) {
            continue;
        }
        unsigned char *d = dest + y * width * 4;
        for (x = 0; x < width; x++, p++) {
            unsigned char opacity = ScaleQuantumToChar(p->opacity);
            if (opacity == 0) {
                *d++ = ScaleQuantumToChar(p->red);
                *d++ = ScaleQuantumToChar(p->green);
                *d++ = ScaleQuantumToChar(p->blue);
                *d++ = 255;
            } else {
                // image.Image wants the alpha premultiplied
                unsigned char alpha = 255 - opacity;
                double factor = alpha / 255.0;
                *d++ = round(ScaleQuantumToChar(p->red) * factor);
                *d++ = round(ScaleQuantumToChar(p->green) * factor);
                *d++ = round(ScaleQuantumToChar(p->blue) * factor);
                *d++ = alpha;
            }
        }
    }
    return 1;
}
