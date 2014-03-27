#include <stdlib.h>
#include <math.h>

#include <magick/api.h>

#include "macros.h"

#define HISTOGRAM_SIZE 768

float
calculate_image_entropy_rect(const Image *image, const RectangleInfo *rect)
{
    register long y;
    register long x;
    register const PixelPacket *p;
    ExceptionInfo ex;
    unsigned int histogram[HISTOGRAM_SIZE] = {0,};
    long sx;
    long sy;
    unsigned int width;
    unsigned int height;
    if (rect && rect->width > 0 && rect->height > 0) {
        sx = rect->x;
        sy = rect->y;
        width = rect->width;
        height = rect->height;
    } else {
        sx = 0;
        sy = 0;
        width = image->columns;
        height = image->rows;
    }
    unsigned int ey = sy + height;
    int total = 0;
    for(y = sy; y < ey; ++y) {
        p = ACQUIRE_IMAGE_PIXELS(image, sx, y, width, 1, &ex);
        if (!p) {
            continue;
        }
        total += width;
        for (x=0; x < width; x++, p++) {
            histogram[ScaleQuantumToChar(p->red)]++;
            histogram[ScaleQuantumToChar(p->green) + 256]++;
            histogram[ScaleQuantumToChar(p->blue) + 512]++;
        }
    }
    float t = (float)(total * 3);
    float entropy = 0.0f;
    int ii;
    for (ii = 0; ii < HISTOGRAM_SIZE; ++ii) {
        float p = histogram[ii] / t;
        if (p) {
            entropy += p * log2f(p);
        }
    }
    return -entropy;
}
