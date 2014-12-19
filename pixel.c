#include <magick/api.h>

#include "quantum.h"

#include "pixel.h"

unsigned char
quantum_to_char(Quantum q)
{
    return ScaleQuantumToChar(q);
}

Quantum
char_to_quantum(unsigned char c)
{
    return ScaleCharToQuantum(c);
}

void
copy_pixel_packets(PixelPacket *src, PixelPacket *dst, int count)
{
    int ii;
    for (ii = 0; ii < count; ii++, src++, dst++) {
        *dst = *src;
    }
}
