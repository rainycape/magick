package magick

// #include <magick/api.h>
//
// void image_matrix(const Image *image, double **out, ExceptionInfo *ex);
import "C"

import (
	"fmt"
	"math"
	"unsafe"
)

type AffineMatrix struct {
	Sx, Sy, Rx, Ry, Tx, Ty float64
}

func (m *AffineMatrix) matrix() *C.AffineMatrix {
	var mat C.AffineMatrix
	mat.sx = C.double(m.Sx)
	mat.sy = C.double(m.Sy)
	mat.rx = C.double(m.Rx)
	mat.ry = C.double(m.Ry)
	mat.tx = C.double(m.Tx)
	mat.ty = C.double(m.Ty)
	return &mat
}

type FloatMatrix [][]float64

func NewFloatMatrix(rows int, columns int) FloatMatrix {
	r := make([][]float64, rows)
	for ii := 0; ii < rows; ii++ {
		r[ii] = make([]float64, columns)
	}
	return FloatMatrix(r)
}

func (m FloatMatrix) Rows() int {
	return len(m)
}

func (m FloatMatrix) Columns() int {
	if len(m) > 0 {
		return len(m[0])
	}
	return 0
}

func (m FloatMatrix) Transposed() FloatMatrix {
	t := NewFloatMatrix(m.Columns(), m.Rows())
	for ii := 0; ii < m.Rows(); ii++ {
		for jj := 0; jj < m.Columns(); jj++ {
			t[ii][jj] = m[jj][ii]
		}
	}
	return t
}

func (m FloatMatrix) SubMatrix(x int, y int, rows int, columns int) (FloatMatrix, error) {
	if x+rows > m.Rows() {
		return nil, fmt.Errorf("can't extract %d rows starting at %d, matrix has %d rows", rows, x, m.Rows())
	}
	if y+columns > m.Columns() {
		return nil, fmt.Errorf("can't extract %d columns starting at %d, matrix has %d columns", columns, y, m.Columns())
	}
	s := NewFloatMatrix(rows, columns)
	for ii := x; ii < rows; ii++ {
		for jj := 0; jj < columns; jj++ {
			s[ii][jj] = m[ii+x][jj+y]
		}
	}
	return s, nil
}

func (m FloatMatrix) Multiply(n FloatMatrix) (FloatMatrix, error) {
	if m.Columns() != n.Rows() {
		return nil, fmt.Errorf("can't multiply matrix with %d columns by matrix with %d rows", m.Columns(), n.Rows())
	}
	p := NewFloatMatrix(m.Rows(), n.Columns())
	end := m.Columns()
	for ii := 0; ii < p.Rows(); ii++ {
		for jj := 0; jj < p.Columns(); jj++ {
			for kk := 0; kk < end; kk++ {
				p[ii][jj] += m[ii][kk] * n[kk][jj]
			}
		}
	}
	return p, nil
}

// UnrollX returns all values in a slice which contains
// all rows in increasing order.
func (m FloatMatrix) UnrollX() []float64 {
	values := make([]float64, 0, m.Rows()*m.Columns())
	for ii := 0; ii < m.Columns(); ii++ {
		for jj := 0; jj < m.Rows(); jj++ {
			values = append(values, m[jj][ii])
		}
	}
	return values
}

// UnrollX returns all values in a slice which contains
// all columns in increasing order.
func (m FloatMatrix) UnrollY() []float64 {
	values := make([]float64, 0, m.Rows()*m.Columns())
	for _, v := range m {
		values = append(values, v...)
	}
	return values
}

func NewDCTMatrix(order int) FloatMatrix {
	m := NewFloatMatrix(order, order)
	c0 := 1 / math.Sqrt(float64(order))
	c1 := math.Sqrt(2 / float64(order))
	for ii := 0; ii < order; ii++ {
		m[ii][0] = c0
		for jj := 1; jj < order; jj++ {
			m[ii][jj] = c1 * math.Cos((math.Pi/2/float64(order))*float64(jj)*(2*float64(ii)+1))
		}
	}
	return m
}

// FloatMatrix returns an Image in the GRAY colorspace as FloatMatrix,
// where each pixel represents an element of the matrix. Each element
// is normalized to the [0,1] range.
func (im *Image) FloatMatrix() (FloatMatrix, error) {
	if im.Colorspace() != GRAY {
		return nil, fmt.Errorf("FloatMatrix() is available only for GRAY images, this one is in %s", im.Colorspace())
	}
	width := im.Width()
	height := im.Height()
	ptrs := make([]*C.double, width)
	for ii := range ptrs {
		ptrs[ii] = allocDoublePtr(height)
	}
	defer func() {
		for _, p := range ptrs {
			freeDoublePtr(p)
		}
	}()
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	C.image_matrix(im.image, (**C.double)(unsafe.Pointer(&ptrs[0])), &ex)
	if ex.severity != C.UndefinedException {
		return nil, exError(&ex, "generating float matrix")
	}
	m := make([][]float64, width)
	for ii := range m {
		m[ii] = doublePtrToFloat64Slice(ptrs[ii], height)
	}
	return FloatMatrix(m), nil
}
