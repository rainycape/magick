package magick

// AUTOMATICALLY GENERATED WITH gondola gen -- DO NOT EDIT!
func (c Colorspace) String() string {
	switch c {
	case CMYK:
		return "CMYK"
	case GRAY:
		return "GRAY"
	case HSL:
		return "HSL"
	case HWB:
		return "HWB"
	case LAB:
		return "LAB"
	case OHTA:
		return "OHTA"
	case REC_601_LUMA:
		return "REC_601_LUMA"
	case REC_601_YCBCR:
		return "REC_601_YCBCR"
	case REC_709_LUMA:
		return "REC_709_LUMA"
	case REC_709_YCBCR:
		return "REC_709_YCBCR"
	case RGB:
		return "RGB"
	case SRGB:
		return "SRGB"
	case TRANSPARENT:
		return "TRANSPARENT"
	case XYZ:
		return "XYZ"
	case YCC:
		return "YCC"
	case YIQ:
		return "YIQ"
	case YPBPR:
		return "YPBPR"
	case YUV:
		return "YUV"
	}
	return "unknown colorspace"
}
func (f Filter) String() string {
	switch f {
	case FBessel:
		return "Bessel"
	case FBlackman:
		return "Blackman"
	case FBox:
		return "Box"
	case FCatrom:
		return "Catrom"
	case FCubic:
		return "Cubic"
	case FGaussian:
		return "Gaussian"
	case FHamming:
		return "Hamming"
	case FHanning:
		return "Hanning"
	case FHermite:
		return "Hermite"
	case FLanczos:
		return "Lanczos"
	case FMitchell:
		return "Mitchell"
	case FPoint:
		return "Point"
	case FQuadratic:
		return "Quadratic"
	case FSinc:
		return "Sinc"
	case FTriangle:
		return "Triangle"
	}
	return "unknown filter"
}
func (o Operator) String() string {
	switch o {
	case OpAdd:
		return "Add"
	case OpAnd:
		return "And"
	case OpAssign:
		return "Assign"
	case OpDepth:
		return "Depth"
	case OpDivide:
		return "Divide"
	case OpGamma:
		return "Gamma"
	case OpGaussianNoise:
		return "GaussianNoise"
	case OpImpulseNoise:
		return "ImpulseNoise"
	case OpLShift:
		return "LShift"
	case OpLaplacianNoise:
		return "LaplacianNoise"
	case OpLog:
		return "Log"
	case OpMax:
		return "Max"
	case OpMin:
		return "Min"
	case OpMultiplicativeNoise:
		return "MultiplicativeNoise"
	case OpMultiply:
		return "Multiply"
	case OpNegate:
		return "Negate"
	case OpOr:
		return "Or"
	case OpPoissonNoise:
		return "PoissonNoise"
	case OpPow:
		return "Pow"
	case OpRShift:
		return "RShift"
	case OpRandomNoise:
		return "RandomNoise"
	case OpSubstract:
		return "Substract"
	case OpThreshold:
		return "Threshold"
	case OpThresholdBlack:
		return "ThresholdBlack"
	case OpThresholdWhite:
		return "ThresholdWhite"
	case OpUniformNoise:
		return "UniformNoise"
	case OpXor:
		return "Xor"
	}
	return "unknown operator"
}
