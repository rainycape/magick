#include <magick/api.h>
#include "bridge.h"

Image *
bridge_image_func(ImageFunc f, Image *im, ExceptionInfo *ex)
{
    return f(im, ex);
}

Image *
bridge_image_data_func(ImageDataFunc f, Image *im, void *data, ExceptionInfo *ex)
{
    return f(im, data, ex);
}

Image *
apply_image_func(ImageFunc f, Image *image, void *parent, int is_coalesced, ExceptionInfo *ex)
{
    // XXX: f must either return new images or the same images in all calls. If sometimes
    // it returns a new one and other times it returns the same image, that will cause either
    // a crash or a memory leak. Since currently the C interface is not exposed to package users
    // there's no need to inform them about this quirk, but must be taken into account by people
    // writing magick code in C.
    if (parent || (!image->next && !image->previous)) {
        return f(image, ex);
    }
    Image *coalesced = NULL;
    if (!is_coalesced) {
        coalesced = CoalesceImages(image, ex);
        if (!coalesced) {
            return NULL;
        }
        image = coalesced;
    }
    Image *ret = f(image, ex);
    if (ret) {
        Image *prev = ret;
        Image *cur;
        Image *res;
        for (cur = image->next; cur; cur = cur->next) {
            res = f(cur, ex);
            if (!res) {
                DestroyImageList(ret);
                ret = NULL;
                break;
            }
            prev->next = res;
            res->previous = prev;
            prev = res;
        }
    }
    if (coalesced && ret != image) {
        DestroyImageList(coalesced);
    }
    return ret;
}

Image *
apply_image_data_func(ImageDataFunc f, Image *image, void *data, void *parent, int is_coalesced, ExceptionInfo *ex)
{
    if (parent || (!image->next && !image->previous)) {
        return f(image, data, ex);
    }
    Image *coalesced = NULL;
    if (!is_coalesced) {
        coalesced = CoalesceImages(image, ex);
        if (!coalesced) {
            return NULL;
        }
        image = coalesced;
    }
    Image *ret = f(image, data, ex);
    if (ret) {
        Image *prev = ret;
        Image *cur;
        Image *res;
        for (cur = image->next; cur; cur = cur->next) {
            res = f(cur, data, ex);
            if (!res) {
                DestroyImageList(ret);
                ret = NULL;
                break;
            }
            prev->next = res;
            res->previous = prev;
            prev = res;
        }
    }
    if (coalesced && ret != image) {
        DestroyImageList(coalesced);
    }
    return ret;
}
