package magick

// #include <magick/api.h>
import "C"

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
