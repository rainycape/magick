#ifndef QUANTUM_H
#define QUANTUM_H

#include <magick/api.h>

#ifndef ScaleQuantumToChar
#define ScaleQuantumToChar(x) (((unsigned char)((x)/(QuantumRange/255.0))))
#define ScaleCharToQuantum(x) (((Quantum)((x)*(QuantumRange/255.0))))
#endif

#endif
