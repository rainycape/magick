package magick

// #include <magick/api.h>
import "C"

type Rect struct {
	X      int
	Y      int
	Width  uint
	Height uint
}

func (r *Rect) rectangleInfo() *C.RectangleInfo {
	var info C.RectangleInfo
	info.x = magickLong(r.X)
	info.y = magickLong(r.Y)
	info.width = magickSize(r.Width)
	info.height = magickSize(r.Height)
	return &info
}

type Frame struct {
	X          int
	Y          int
	Width      uint
	Height     uint
	InnerBevel int
	OuterBevel int
}

func (f *Frame) frameInfo() *C.FrameInfo {
	var info C.FrameInfo
	info.x = magickLong(f.X)
	info.y = magickLong(f.Y)
	info.width = magickSize(f.Width)
	info.height = magickSize(f.Height)
	info.inner_bevel = magickLong(f.InnerBevel)
	info.outer_bevel = magickLong(f.OuterBevel)
	return &info
}
