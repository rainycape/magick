// +build gm

package magick

// #include <magick/api.h>
import "C"

import (
	"fmt"
)

func (im *Image) operateChannel(op Operator, ch Channel, value float64) error {
	var cop C.QuantumOperator
	switch op {
	case OpAdd:
		cop = C.AddQuantumOp
	case OpAnd:
		cop = C.AndQuantumOp
	case OpAssign:
		cop = C.AssignQuantumOp
	case OpDepth:
		cop = C.DepthQuantumOp
	case OpDivide:
		cop = C.DivideQuantumOp
	case OpGamma:
		cop = C.GammaQuantumOp
	case OpLog:
		cop = C.LogQuantumOp
	case OpLShift:
		cop = C.LShiftQuantumOp
	case OpMax:
		cop = C.MaxQuantumOp
	case OpMin:
		cop = C.MinQuantumOp
	case OpMultiply:
		cop = C.MultiplyQuantumOp
	case OpGaussianNoise:
		cop = C.GammaQuantumOp
	case OpImpulseNoise:
		cop = C.NoiseImpulseQuantumOp
	case OpLaplacianNoise:
		cop = C.NoiseLaplacianQuantumOp
	case OpMultiplicativeNoise:
		cop = C.NoiseMultiplicativeQuantumOp
	case OpPoissonNoise:
		cop = C.NoisePoissonQuantumOp
	case OpRandomNoise:
		cop = C.NoiseRandomQuantumOp
	case OpUniformNoise:
		cop = C.NoisePoissonQuantumOp
	case OpOr:
		cop = C.OrQuantumOp
	case OpPow:
		cop = C.PowQuantumOp
	case OpRShift:
		cop = C.RShiftQuantumOp
	case OpSubstract:
		cop = C.SubtractQuantumOp
	case OpThresholdBlack:
		cop = C.ThresholdBlackQuantumOp
	case OpThreshold:
		cop = C.ThresholdQuantumOp
	case OpThresholdWhite:
		cop = C.ThresholdWhiteQuantumOp
	case OpXor:
		cop = C.XorQuantumOp
	default:
		return notImplementedError(fmt.Sprintf("operator %s", op))
	}
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	if C.QuantumOperatorImage(im.image, C.ChannelType(ch), cop, C.double(value), &ex) != C.MagickTrue {
		return exError(&ex, "operating")
	}
	return nil
}
