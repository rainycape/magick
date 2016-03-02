package magick

import (
	"sort"
	"strconv"
)

// PHash represents a perceptual hash. Use Image.PHash()
// to obtain the perceptual hash of an image.
type PHash uint64

func (p PHash) Uint64() uint64 {
	return uint64(p)
}

func (p PHash) String() string {
	return strconv.FormatUint(uint64(p), 16)
}

// Compare returns a number between in the [0,1] interval.
// The difference is calculated as number of different bits
// / 64. This means 0 indicates that the hashes are equal,
// while 1 indicates that the hashes have no bits in common.
func (p PHash) Compare(q PHash) float64 {
	var c uint
	v := uint64(p) ^ uint64(q)
	// Brian Kernighan style, O(number_of_bits_set)
	for c = 0; v != 0; c++ {
		v &= v - 1
	}
	return float64(c) / 64.0
}

// PHash returns a perceptual hash of the image. The algorithm is based
// on the pHash library and works as follows:
//
//  1 - Transform the image to grayscale.
//  2 - Apply convolution using a (7, []float64{1}) kernel.
//  3 - Scale to 32x32.
//  4 - Compute the DCT for the image.
//  5 - Grab the 8x8 submatrix starting at (1, 1).
//  6 - Unroll the matrix by the x axis into an array.
//  7 - Compute the median for the previous array.
//  8 - Generate a 64bit hash, where bit x = 1 iff array[x] > median.
//
// Note that due to differences in scaling and convolution algorithms, this
// function will generate different values for the same image depending on
// the backend library (IM or GM).
func (im *Image) PHash() (PHash, error) {
	c, err := im.Clone()
	if err != nil {
		return 0, err
	}
	defer c.Dispose()
	if err := c.ToColorspace(GRAY); err != nil {
		return 0, err
	}
	convolved, err := c.Convolve(7, []float64{1})
	if err != nil {
		return 0, err
	}
	defer convolved.Dispose()
	scaled, err := convolved.Scale(32, 32)
	if err != nil {
		return 0, err
	}
	defer scaled.Dispose()
	imMatrix, err := scaled.FloatMatrix()
	if err != nil {
		return 0, err
	}
	dctMatrix := NewDCTMatrix(32)
	// We can safely ignore errors here, since the sizes
	// always match.
	dct, _ := dctMatrix.Multiply(imMatrix)
	dct, _ = dct.Multiply(dctMatrix.Transposed())
	sub, _ := dct.SubMatrix(1, 1, 8, 8)
	values := sub.UnrollX()
	// We need the original values, so we must sort a copy
	cpy := make([]float64, len(values))
	copy(cpy, values)
	sort.Float64s(cpy)
	median := (cpy[64/2-1] + cpy[64/2]) / 2
	bit := uint64(1)
	hash := uint64(0)
	for _, v := range values {
		if v > median {
			hash |= bit
		}
		bit <<= 1
	}
	return PHash(hash), nil
}
