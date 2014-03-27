#include <string.h>
#include <stdlib.h>
#include <stdint.h>

#include <magick/api.h>
// This is required because both GM and
// giflib define a function named
// DrawRectangle
#define DrawRectangle _DrawRectangle
#include <gif_lib.h>
#undef DrawRectangle

#include "macros.h"

#define NCOLORS 256
#define GIF_APP "NETSCAPE2.0"


typedef struct {
    GifPixelType *data;
    int x;
    int y;
    int width;
    int height;
    int duration;
} Frame;

void
gif_frames_free(Frame *frames, int count)
{
    int ii;
    for (ii = 0; ii < count; ii++) {
        if (frames[ii].data) {
            free(frames[ii].data);
        }
    }
    free(frames);
}

typedef struct {
    GifByteType *data;
    int size;
    int alloc;
} GifBuffer;

int
gif_buffer_write(GifFileType *ft, const GifByteType *bytes, int len)
{
    GifBuffer *buf = ft->UserData;
    if (buf->size + len > buf->alloc) {
        int new_size = buf->alloc * 1.25;
        if (new_size < buf->size + len) {
            new_size = buf->size + len;
        }
        buf->data = realloc(buf->data, new_size);
        buf->alloc = new_size;
    }
    memcpy(buf->data + buf->size, bytes, len);
    buf->size += len;
    return len;
}

void *
gif_save(const Image *image, const ColorMapObject *color_map, Frame *frames, int count, int *size)
{
    GifBuffer buf = {0,};
    int estimated = count * (image->columns * image->rows);
    buf.alloc = estimated;
    buf.data = malloc(estimated);
    GifFileType *gif_file = EGifOpen(&buf, gif_buffer_write);
    if (!gif_file) {
        return NULL;
    }
    if (EGifPutScreenDesc(gif_file, image->columns, image->rows, NCOLORS, 0, color_map) == GIF_ERROR) {
        EGifCloseFile(gif_file);
        return NULL;
    }
    if (EGifPutExtensionFirst(gif_file, APPLICATION_EXT_FUNC_CODE, strlen(GIF_APP), GIF_APP) == GIF_ERROR) {
        EGifCloseFile(gif_file);
        return NULL;
    }
    unsigned char meta[] = {
        0x01, //  data sub-block index (always 1)
        0xFF, 0xFF // 65535 repetitions - unsigned
    };
    if (EGifPutExtensionLast(gif_file, APPLICATION_EXT_FUNC_CODE, sizeof(meta), meta) == GIF_ERROR) {
        EGifCloseFile(gif_file);
        return NULL;
    }
    int ii;
    for (ii = 0; ii < count; ii++) {
        Frame *frame = &frames[ii];
        // GCE
        unsigned char gce[] = {
            0x08, // no transparency
            frame->duration % 256, // LSB of delay
            frame->duration / 256, // MSB of delay in millisecs
            0x00, // no transparent color
        };
        if (EGifPutExtension(gif_file, GRAPHICS_EXT_FUNC_CODE, sizeof(gce), gce) == GIF_ERROR) {
            EGifCloseFile(gif_file);
            return NULL;
        }
        if (EGifPutImageDesc(gif_file, frame->x, frame->y, frame->width, frame->height, 0, NULL) == GIF_ERROR) {
            EGifCloseFile(gif_file);
            return NULL;
        }
        int yy;
        GifPixelType *p = frame->data;
        for (yy = 0; yy < frame->height; yy++, p += frame->width) {
            if (EGifPutLine(gif_file, p, frame->width) == GIF_ERROR) {
                EGifCloseFile(gif_file);
                return NULL;
            }
        }
    } 
    EGifCloseFile(gif_file);
    *size = buf.size;
    return buf.data;
}

int
acquire_image_pixels(const Image *image, GifByteType *red, GifByteType *green, GifByteType *blue)
{
    register long y;
    register long x;
    register const PixelPacket *p;
    ExceptionInfo ex;
    int width = image->columns;
    int height = image->rows;
    int ii = 0;
    for(y = 0; y < height; ++y) {
        p = ACQUIRE_IMAGE_PIXELS(image, 0, y, width, 1, &ex);
        if (!p) {
            return 0;
        }
        for (x = 0; x < width; x++, ii++, p++) {
            red[ii] = ScaleQuantumToChar(p->red);
            green[ii] = ScaleQuantumToChar(p->green);
            blue[ii] = ScaleQuantumToChar(p->blue);
        }
    }
    return 1;
}

int
aprox_image_pixels(const Image *image, GifColorType *colors, int color_count, GifPixelType *out)
{
    int width = image->columns;
    int height = image->rows;
    ExceptionInfo ex;
    int ii;
    int jj;
    register const PixelPacket *p = ACQUIRE_IMAGE_PIXELS(image, 0, 0, width, height, &ex);
    if (!p) {
        return 0;
    }
    int end = width * height;
    for (ii = 0; ii < end; ii++, p++) {
        GifByteType r = ScaleQuantumToChar(p->red);
        GifByteType g = ScaleQuantumToChar(p->green);
        GifByteType b = ScaleQuantumToChar(p->blue);
        uint32_t key = (r << 16) + (g << 8) + b;
        int min_delta = 3 * (NCOLORS * NCOLORS);
        int min_pos = 0;
        GifColorType *c = colors;
        for (jj = 0; jj < color_count; jj++, c++) {
            int rd = c->Red - r;
            int gd = c->Green - g;
            int bd = c->Blue - b;
            int delta = (rd * rd) + (gd * gd) + (bd * bd);
            if (delta < min_delta) {
                min_delta = delta;
                min_pos = jj;
                if (min_delta == 0) {
                    break;
                }
            }
        }
        out[ii] = min_pos;
    }
    return 1;
}

void *
gif_encode(Image *image, int single, int *size)
{
    int width = image->columns;
    int height = image->rows;
    int total = width * height;
    GifByteType output[total];
    GifByteType red[total];
    GifByteType green[total];
    GifByteType blue[total];

    if (!acquire_image_pixels(image, red, green, blue)) {
        return NULL;
    }

    int count = GetImageListLength(image);
    Frame *frames = calloc(count, sizeof(*frames));

    int palette_size = NCOLORS;
    ColorMapObject *palette = MakeMapObject(palette_size, NULL);
    if (QuantizeBuffer(width, height, &palette_size, red, green, blue, output, palette->Colors) == GIF_ERROR) {
        FreeMapObject(palette);
        gif_frames_free(frames, count);
        return NULL;
    }

    frames[0].data = malloc(total);
    memcpy(frames[0].data, output, total);
    frames[0].width = width;
    frames[0].height = height;
    frames[0].duration = image->delay;
    GifColorType *colors = palette->Colors;
    int color_count = palette->ColorCount;
    int ii;
    Image *cur = image->next;
    for (ii = 1; ii < count; ii++, cur = cur->next) {
        frames[ii].width = width;
        frames[ii].height = height;
        frames[ii].duration = cur->delay;
        GifPixelType *data = malloc(total);
        frames[ii].data = data;
        
        if (!aprox_image_pixels(cur, colors, color_count, data)) {
            FreeMapObject(palette);
            gif_frames_free(frames, count);
            return NULL;
        }
    }
    void *ret = gif_save(image, palette, frames, count, size);
    FreeMapObject(palette);
    gif_frames_free(frames, count);
    return ret;
}
