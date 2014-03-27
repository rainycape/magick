#if __GNUC__ >= 4 && __GNUC_MINOR__ >= 6
#define ACQUIRE_IMAGE_PIXELS(im, x, y, w, h, ex) ({ \
    _Pragma("GCC diagnostic push") \
    _Pragma("GCC diagnostic ignored \"-Wdeprecated-declarations\"") \
    AcquireImagePixels(im, x, y, w, h, ex); \
    _Pragma("GCC diagnostic pop") \
})
#else
#define ACQUIRE_IMAGE_PIXELS(im, x, y, w, h, ex) AcquireImagePixels(im, x, y, w, h, ex)
#endif
