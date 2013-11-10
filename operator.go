package magick

// Operator represents an operation to be applied on all the pixels
// of the selected channel(s) of an image.
type Operator int

const (
	OpAdd                 = iota + 1 /* Add value */
	OpAnd                            /* Bitwise AND value */
	OpAssign                         /* Direct value assignment */
	OpDepth                          /* Divide by value */
	OpDivide                         /* Bitwise left-shift value N bits */
	OpGamma                          /* Adjust image gamma */
	OpLog                            /* log(quantum*value+1)/log(value+1) */
	OpLShift                         /* Bitwise left-shift value N bits */
	OpMax                            /* Assign value if > quantum */
	OpMin                            /* Assign value if < quantum */
	OpMultiply                       /* Multiply by value */
	OpNegate                         /* Negate channel, ignore value */
	OpGaussianNoise                  /* Gaussian noise */
	OpImpulseNoise                   /* Impulse noise */
	OpLaplacianNoise                 /* Laplacian noise */
	OpMultiplicativeNoise            /* Multiplicative gaussian noise */
	OpPoissonNoise                   /* Poisson noise */
	OpRandomNoise                    /* Random noise */
	OpUniformNoise                   /* Uniform noise */
	OpOr                             /* Bitwise OR value */
	OpPow                            /* Power function: pow(quantum,value) */
	OpRShift                         /* Bitwise right shift value */
	OpSubstract                      /* Subtract value */
	OpThresholdBlack                 /* Below threshold is black */
	OpThreshold                      /* Above threshold white, otherwise black */
	OpThresholdWhite                 /* Above threshold is white */
	OpXor                            /* Bitwise XOR value */
)

// Operate is a shorthand for OperateChannel(op, CAll, value).
func (im *Image) Operate(op Operator, value float64) error {
	return im.OperateChannel(op, CAll, value)
}

// Operate applies the selected operator to the selected channel(s) of the
// image. The image is modified and only errors are returned. This method
// works only on single images, so if you're working with animations you
// must perform the coalescing and then apply the operator to each frame.
// See the Operator type for the available operations.
func (im *Image) OperateChannel(op Operator, ch Channel, value float64) error {
	return im.operateChannel(op, ch, value)
}
