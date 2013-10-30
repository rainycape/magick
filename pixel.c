#include <magick/api.h>

#ifndef ScaleQuantumToChar
#define ScaleQuantumToChar(x) (((unsigned char)((x)/(QuantumRange/255.0))))
#define ScaleCharToQuantum(x) (((Quantum)((x)/(QuantumRange/255.0))))
#endif

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
