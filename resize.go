package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "resize.h"
import "C"

import (
	"fmt"
)

type Filter C.FilterTypes

const (
	FPoint     Filter = C.PointFilter
	FBox       Filter = C.BoxFilter
	FTriangle  Filter = C.TriangleFilter
	FHermite   Filter = C.HermiteFilter
	FHanning   Filter = C.HanningFilter
	FHamming   Filter = C.HammingFilter
	FBlackman  Filter = C.BlackmanFilter
	FGaussian  Filter = C.GaussianFilter
	FQuadratic Filter = C.QuadraticFilter
	FCubic     Filter = C.CubicFilter
	FCatrom    Filter = C.CatromFilter
	FMitchell  Filter = C.MitchellFilter
	FLanczos   Filter = C.LanczosFilter
	FBessel    Filter = C.BesselFilter
	FSinc      Filter = C.SincFilter
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
// filter and blur factor (where > 1 is blurry, < 1 is sharp). If width or height is
// < 0, it's calculated proportionally to the other dimension. If both of them are < 0,
// an error is returned.
func (im *Image) ResizeBlur(width, height int, filter Filter, blur float64) (*Image, error) {
	var data C.ResizeData
	if width < 0 {
		if height < 0 {
			return nil, fmt.Errorf("invalid resize %dx%d", width, height)
		}
		h := float64(im.Height())
		var ratio float64
		if h != 0 {
			ratio = float64(im.Width()) / h
		}
		width = int(float64(height) * ratio)
	}
	if height < 0 {
		if width < 0 {
			return nil, fmt.Errorf("invalid resize %dx%d", width, height)
		}
		var ratio float64
		w := float64(im.Width())
		if w != 0 {
			ratio = float64(im.Height()) / w
		}
		height = int(float64(width) * ratio)
	}
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
