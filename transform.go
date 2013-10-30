package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "transform.h"
import "C"

import (
	"math"
	"unsafe"
)

// CropStrategy defines the strategy used to
// crop an image to a given ratio.
type CropStrategy int

const (
	// When cropping to a given aspect ratio, discard content
	// from the edges of the image and keep the content at the center.
	CSCenter CropStrategy = iota
	// When cropping to a given aspect ratio, grab the section of the
	// image with higher entropy. This is useful when you want to crop
	// the image to a given ratio and end up with the most 'meaningful'
	// part of the image.
	CSMaxEntropy
)

func (im *Image) Chop(r Rect) (*Image, error) {
	return im.applyRectFunc("chopping", C.ImageDataFunc(C.ChopImage), r)
}

// Coalesce composites a set of images while respecting any page offsets and
// disposal methods. GIF, MIFF, and MNG animation sequences typically start
// with an image background and each subsequent image varies in size and
// offset. Coalesce() returns a new sequence where each image in the sequence
// is the same size as the first and composited with the next image in the
// sequence.
func (im *Image) Coalesce() (*Image, error) {
	if im.coalesced {
		return im, nil
	}
	res, err := im.fn("coalescing", C.ImageFunc(C.CoalesceImages))
	if res != nil {
		res.coalesced = true
	}
	return res, err
}

func (im *Image) Crop(r Rect) (*Image, error) {
	return im.applyRectFunc("cropping", C.ImageDataFunc(C.CropImage), r)
}

func (im *Image) Deconstruct() (*Image, error) {
	return im.fn("deconstructing", C.ImageFunc(C.DeconstructImages))
}

func (im *Image) Extent(r Rect) (*Image, error) {
	return im.applyRectFunc("extenting", C.ImageDataFunc(C.ExtentImage), r)
}

// Flatten merges a sequence of images. This is useful for combining Photoshop
// layers into a single image.
func (im *Image) Flatten() (*Image, error) {
	return im.fn("flattening", C.ImageFunc(C.flattenImages))
}

// Flip creates a vertical mirror image by reflecting the pixels around
// the central x-axis.
func (im *Image) Flip() (*Image, error) {
	return im.applyFunc("flipping", C.ImageFunc(C.FlipImage))
}

// Flop creates a horizontal mirror image by reflecting the pixels around
// the central y-axis.
func (im *Image) Flop() (*Image, error) {
	return im.applyFunc("flopping", C.ImageFunc(C.FlopImage))
}

func (im *Image) Mosaic() (*Image, error) {
	return im.fn("mosaic", C.ImageFunc(C.mosaicImages))
}

func (im *Image) Roll(xoffset, yoffset int) (*Image, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	rolled := C.RollImage(im.image, magickLong(xoffset), magickLong(yoffset), &ex)
	return checkImage(rolled, nil, &ex, "rolling")
}

func (im *Image) Shave(r Rect) (*Image, error) {
	return im.applyRectFunc("shaving", C.ImageDataFunc(C.ShaveImage), r)
}

func (im *Image) Transform(crop, image string) *Image {
	ccrop := C.CString(crop)
	defer C.free(unsafe.Pointer(ccrop))
	cimage := C.CString(image)
	defer C.free(unsafe.Pointer(cimage))
	current := im.image
	C.TransformImage(&current, ccrop, cimage)
	if current == im.image {
		/* An error ocurred, return original image */
		return im
	}
	return newImage(current, nil)
}

func (im *Image) cropToRatioMaxEntropy(ratio, imWidth, imHeight, imRatio float64) (*Image, error) {
	r := Rect{0, 0, uint(im.Width()), uint(im.Height())}
	for math.Abs(ratio-imRatio) >= 0.05 {
		var ra Rect
		var rb Rect
		if imRatio > ratio {
			remaining := imWidth - (imHeight * ratio)
			px := uint(math.Min(10, remaining))
			w := uint(imWidth) - px
			h := uint(imHeight)
			ra = Rect{r.X, 0, w, h}
			rb = Rect{r.X + int(px), 0, w, h}
		} else {
			remaining := imHeight - (imWidth / ratio)
			px := uint(math.Min(10, remaining))
			w := uint(imWidth)
			h := uint(imHeight) - px
			ra = Rect{0, r.Y, w, h}
			rb = Rect{0, r.Y + int(px), w, h}
		}
		if im.EntropyRect(ra) > im.EntropyRect(rb) {
			r = ra
		} else {
			r = rb
		}
		newWidth := float64(r.Width)
		newHeight := float64(r.Height)
		if newWidth == imWidth && newHeight == imHeight {
			// The image can't be resized the requested ratio
			// e.g. image is 1x1 and requested ratio is 3
			break
		}
		imWidth = newWidth
		imHeight = newHeight
		imRatio = imWidth / imHeight
	}
	return im.Crop(r)
}

func (im *Image) cropToRatioCenter(ratio, imWidth, imHeight, imRatio float64) (*Image, error) {
	var r Rect
	if imRatio > ratio {
		// Crop width
		cropWidth := int(imHeight * ratio)
		remaining := im.Width() - cropWidth
		r.X = int(remaining / 2)
		r.Width = uint(cropWidth)
		r.Y = 0
		r.Height = uint(im.Height())
	} else {
		// Crop height
		cropHeight := int(imWidth / ratio)
		remaining := im.Height() - cropHeight
		r.Y = int(remaining / 2)
		r.Height = uint(cropHeight)
		r.X = 0
		r.Width = uint(im.Width())
	}
	return im.Crop(r)
}

func (im *Image) CropToRatio(ratio float64, cs CropStrategy) (*Image, error) {
	imWidth := float64(im.Width())
	imHeight := float64(im.Height())
	imRatio := imWidth / imHeight
	if ratio != imRatio {
		if cs == CSMaxEntropy {
			return im.cropToRatioMaxEntropy(ratio, imWidth, imHeight, imRatio)
		}
		return im.cropToRatioCenter(ratio, imWidth, imHeight, imRatio)
	}
	return im, nil
}

func (im *Image) CropResize(width, height int, filter Filter, cs CropStrategy) (*Image, error) {
	if width == 0 && height == 0 {
		// Return a clone, since the caller might Dispose the original
		// after calling this function to save memory
		return im.Clone()
	}
	imWidth := float64(im.Width())
	imHeight := float64(im.Height())
	ratio := imWidth / imHeight
	if width == 0 {
		width = int(float64(height) * ratio)
	} else if height == 0 {
		height = int(float64(width) / ratio)
	}
	targetRatio := float64(width) / float64(height)
	cropped, err := im.CropToRatio(targetRatio, cs)
	if err != nil {
		return nil, err
	}
	res, err := cropped.Resize(width, height, filter)
	if cropped != im {
		cropped.Dispose()
	}
	return res, err
}
