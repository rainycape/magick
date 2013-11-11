package magick

// AUTOMATICALLY GENERATED WITH gondola gen -- DO NOT EDIT!
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
