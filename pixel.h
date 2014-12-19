#ifndef PIXEL_H
#define PIXEL_H

#include <magick/api.h>

unsigned char quantum_to_char(Quantum q);
Quantum char_to_quantum(unsigned char c);

void copy_pixel_packets(PixelPacket *src, PixelPacket *dst, int count);

#endif
