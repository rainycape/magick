package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "resize.h"
import "C"

type Filter C.FilterTypes

const (
	FPoint     = C.PointFilter
	FBox       = C.BoxFilter
	FTriangle  = C.TriangleFilter
	FHermite   = C.HermiteFilter
	FHanning   = C.HanningFilter
	FHamming   = C.HammingFilter
	FBlackman  = C.BlackmanFilter
	FGaussian  = C.GaussianFilter
	FQuadratic = C.QuadraticFilter
	FCubic     = C.CubicFilter
	FCatrom    = C.CatromFilter
	FMitchell  = C.MitchellFilter
	FLanczos   = C.LanczosFilter
	FBessel    = C.BesselFilter
	FSinc      = C.SincFilter
)

// Magnify is a convenience method that scales an image proportionally to twice its size.
func (im *Image) Magnify() (*Image, error) {
	return im.applyFunc("magnifying", C.ImageFunc(C.MagnifyImage))
}

// Minify is a convenience method that scales an image proportionally to half its size.
func (im *Image) Minify() (*Image, error) {
	return im.applyFunc("minifying", C.ImageFunc(C.MinifyImage))
}

// ResizeBlur returns a new image resized to the given dimensions using the provided
// filter and blur factor (where > 1 is blurry, < 1 is sharp).
func (im *Image) ResizeBlur(width, height int, filter Filter, blur float64) (*Image, error) {
	var data C.ResizeData
	data.columns = C.ulong(width)
	data.rows = C.ulong(height)
	data.filter = C.FilterTypes(filter)
	data.blur = C.double(blur)
	return im.applyDataFunc("resizing", C.ImageDataFunc(C.resizeImage), &data)
}

// Resize works like ResizeBlur, but sets the blur to 1
func (im *Image) Resize(width, height int, filter Filter) (*Image, error) {
	return im.ResizeBlur(width, height, filter, 1)
}

func (im *Image) sizeFunc(what string, width, height int, f C.ImageDataFunc) (*Image, error) {
	var s C.SizeData
	s.columns = C.ulong(width)
	s.rows = C.ulong(height)
	return im.applyDataFunc(what, f, &s)
}

// Sample scales an image to the desired dimensions with pixel sampling.
// Unlike other scaling methods, this method does not introduce any
// additional color into the scaled image.
func (im *Image) Sample(width, height int) (*Image, error) {
	return im.sizeFunc("sampling", width, height, C.ImageDataFunc(C.sampleImage))
}

// Scale changes the size of an image to the given dimensions.
func (im *Image) Scale(width, height int) (*Image, error) {
	return im.sizeFunc("scaling", width, height, C.ImageDataFunc(C.scaleImage))
}

// Thumbnail changes the size of an image to the given dimensions. This
// method was designed by Bob Friesenhahn as a low cost thumbnail generator.
func (im *Image) Thumbnail(width, height int) (*Image, error) {
	return im.sizeFunc("thumbnailing", width, height, C.ImageDataFunc(C.thumbnailImage))
}
