// +build !gm

package magick

// #include <magick/api.h>
import "C"

import (
	"fmt"
)

func (im *Image) operateChannel(op Operator, ch Channel, value float64) error {
	var cop C.MagickEvaluateOperator
	switch op {
	case OpAdd:
		cop = C.AddEvaluateOperator
	case OpAnd:
		cop = C.AndEvaluateOperator
	case OpAssign:
		cop = C.SetEvaluateOperator
	case OpDivide:
		cop = C.DivideEvaluateOperator
	case OpLog:
		cop = C.LogEvaluateOperator
	case OpLShift:
		cop = C.LeftShiftEvaluateOperator
	case OpMultiply:
		cop = C.MultiplyEvaluateOperator
	case OpGaussianNoise:
		cop = C.GaussianNoiseEvaluateOperator
	case OpImpulseNoise:
		cop = C.ImpulseNoiseEvaluateOperator
	case OpLaplacianNoise:
		cop = C.LaplacianNoiseEvaluateOperator
	case OpMultiplicativeNoise:
		cop = C.MultiplicativeNoiseEvaluateOperator
	case OpPoissonNoise:
		cop = C.PoissonNoiseEvaluateOperator
	case OpUniformNoise:
		cop = C.UniformNoiseEvaluateOperator
	case OpOr:
		cop = C.OrEvaluateOperator
	case OpPow:
		cop = C.PowEvaluateOperator
	case OpRShift:
		cop = C.RightShiftEvaluateOperator
	case OpSubstract:
		cop = C.SubtractEvaluateOperator
	case OpThresholdBlack:
		cop = C.ThresholdBlackEvaluateOperator
	case OpThreshold:
		cop = C.ThresholdEvaluateOperator
	case OpThresholdWhite:
		cop = C.ThresholdWhiteEvaluateOperator
	case OpXor:
		cop = C.XorEvaluateOperator
	default:
		return notImplementedError(fmt.Sprintf("operator %s", op))
	}
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	if C.EvaluateImageChannel(im.image, C.ChannelType(ch), cop, C.double(value), &ex) != C.MagickTrue {
		return exError(&ex, "operating")
	}
	return nil
}
