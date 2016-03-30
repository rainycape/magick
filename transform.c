#include <stdlib.h>
#include <math.h>
#include <string.h>

#include <magick/api.h>

#include "macros.h"

#define HISTOGRAM_SIZE 768

Image *flattenImages(const Image *image, ExceptionInfo *exception)
{
#ifdef _MAGICKCORE_LAYER_H
    return MergeImageLayers((Image *)image, FlattenLayer, exception);
#else
    return FlattenImages(image, exception);
#endif
}

Image *mosaicImages(const Image *image, ExceptionInfo *exception)
{
#ifdef _MAGICKCORE_LAYER_H
    return MergeImageLayers((Image *)image, MosaicLayer, exception);
#else
    return MosaicImages(image, exception);
#endif
}

float
calculate_entropy_with_histogram(const unsigned int histogram[], const int total) {
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

void
calculate_image_histogram_in_rect(const Image *image, const RectangleInfo *rect, unsigned int histogram[], int o)
{
    register long y;
    register long x;
    register const PixelPacket *p;
    ExceptionInfo ex;
    long sx;
    long sy;
    unsigned int width;
    unsigned int height;
    if (o > 0) {
        memset(histogram, 0, sizeof(unsigned int) * HISTOGRAM_SIZE);
    }
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
            histogram[ScaleQuantumToChar(p->red)] += o;
            histogram[ScaleQuantumToChar(p->green) + 256] += o;
            histogram[ScaleQuantumToChar(p->blue) + 512] += o;
        }
    }
}

RectangleInfo *
calculate_max_entropy_rect(double ratio, const Image *image, double imWidth, double imHeight, double imRatio) {
    unsigned int remaining;
    unsigned int px;
    unsigned int w;
    unsigned int h;
    RectangleInfo ra;
    RectangleInfo rb;
    RectangleInfo *r = (RectangleInfo*)malloc(sizeof(RectangleInfo));
    unsigned int histogram[HISTOGRAM_SIZE];
    unsigned int histogram_a[HISTOGRAM_SIZE];
    unsigned int histogram_b[HISTOGRAM_SIZE];

    r->x = 0; r->y = 0; r->width = image->columns; r->height = image->rows;
    calculate_image_histogram_in_rect(image, r, histogram, 1);
    memcpy(histogram_a, histogram, sizeof(unsigned int) * HISTOGRAM_SIZE);
    memcpy(histogram_b, histogram, sizeof(unsigned int) * HISTOGRAM_SIZE);

    while (ratio - imRatio >= 0.001 || ratio - imRatio <= -0.001) {
        if (imRatio > ratio) {
            remaining = (unsigned int)(imWidth - (imHeight * ratio));
            px = remaining < 10 ? remaining : 10;
            w = (unsigned int)(imWidth - px);
            h = (unsigned int)(imHeight);

            ra.x = r->x + w; ra.y = 0; ra.width = px; ra.height = h;
            rb.x = r->x; rb.y = 0; rb.width = px; rb.height = h;
            calculate_image_histogram_in_rect(image, &ra, histogram_a, -1);
            calculate_image_histogram_in_rect(image, &rb, histogram_b, -1);

            ra.x = r->x; ra.y = 0; ra.width = w; ra.height = h;
            rb.x = r->x + px; rb.y = 0; rb.width = w; rb.height = h;
        } else {
            remaining = (unsigned int)(imHeight - (imWidth / ratio));
            px = remaining < 10 ? remaining : 10;
            w = (unsigned int)(imWidth);
            h = (unsigned int)(imHeight - px);

            ra.x = 0; ra.y = r->y + h; ra.width = w; ra.height = px;
            rb.x = 0; rb.y = r->y; rb.width = w; rb.height = px;
            calculate_image_histogram_in_rect(image, &ra, histogram_a, -1);
            calculate_image_histogram_in_rect(image, &rb, histogram_b, -1);

            ra.x = 0; ra.y = r->y; ra.width = w; ra.height = h;
            rb.x = 0; rb.y = r->y + px; rb.width = w; rb.height = h;
        }

        if (calculate_entropy_with_histogram(histogram_a, ra.width * ra.height) > calculate_entropy_with_histogram(histogram_b, rb.width * rb.height)) {
            r->x = ra.x; r->y = ra.y; r->width = ra.width; r->height = ra.height;
            memcpy(histogram, histogram_a, sizeof(unsigned int) * HISTOGRAM_SIZE);
            memcpy(histogram_b, histogram, sizeof(unsigned int) * HISTOGRAM_SIZE);
        } else {
            r->x = rb.x; r->y = rb.y; r->width = rb.width; r->height = rb.height;
            memcpy(histogram, histogram_b, sizeof(unsigned int) * HISTOGRAM_SIZE);
            memcpy(histogram_a, histogram, sizeof(unsigned int) * HISTOGRAM_SIZE);
        }

        if ((int)(imWidth + 0.5) == r->width && (int)(imHeight + 0.5) == r->height) {
            return r;
        }

        imWidth = r->width;
        imHeight = r->height;
        imRatio = imWidth / imHeight;
    }
    return r;
}
