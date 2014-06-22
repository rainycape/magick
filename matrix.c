#include <magick/api.h>

#include "macros.h"
#include "quantum.h"

void
image_matrix(const Image *image, double **out, ExceptionInfo *ex) {
    register long y;
    register long x;
    register const PixelPacket *p;
    unsigned int width = image->columns;
    for(y = 0; y < image->rows; ++y) {
        p = ACQUIRE_IMAGE_PIXELS(image, 0, y, width, 1, ex);
        if (!p) {
            continue;
        }
        for (x = 0; x < width; x++, p++) {
            // image is GRAY, so all color channels are the same.
            out[x][y] = ScaleQuantumToChar(p->red) / 255.0;
        }
    }
}
