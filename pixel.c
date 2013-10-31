#include <magick/api.h>

#include "quantum.h"

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
