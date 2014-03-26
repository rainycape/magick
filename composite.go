package magick

// #include <magick/api.h>
// #include "bridge.h"
// #include "composite.h"
import "C"

// Composite represents a composition operation. Refer
// to the individual constants for more information. Note
// that not all compositions are supported by the
// GraphicsMagick backend.
type Composite int

const (
	CompositeAlpha Composite = iota + 1
	CompositeAtop
	CompositeBlend
	CompositeBlur
	CompositeBumpmap
	CompositeChangeMask
	CompositeClear
	CompositeColorBurn
	CompositeColorDodge
	CompositeColorize
	CompositeCopyBlack
	CompositeCopyBlue
	CompositeCopy
	CompositeCopyCyan
	CompositeCopyGreen
	CompositeCopyMagenta
	CompositeCopyAlpha
	CompositeCopyRed
	CompositeCopyYellow
	CompositeDarken
	CompositeDarkenIntensity
	CompositeDifference
	CompositeDisplace
	CompositeDissolve
	CompositeDistort
	CompositeDivideDst
	CompositeDivideSrc
	CompositeDstAtop
	CompositeDst
	CompositeDstIn
	CompositeDstOut
	CompositeDstOver
	CompositeExclusion
	CompositeHardLight
	CompositeHue
	CompositeIn
	CompositeIntensity
	CompositeLighten
	CompositeLightenIntensity
	CompositeLinearBurn
	CompositeLinearDodge
	CompositeLinearLight
	CompositeLuminize
	CompositeMathematics
	CompositeMinusDst
	CompositeMinusSrc
	CompositeModulate
	CompositeModulusAdd
	CompositeModulusSubtract
	CompositeMultiply
	CompositeNo
	CompositeOut
	CompositeOver
	CompositeOverlay
	CompositePegtopLight
	CompositePinLight
	CompositePlus
	CompositeReplace
	CompositeSaturate
	CompositeScreen
	CompositeSoftLight
	CompositeSrcAtop
	CompositeSrc
	CompositeSrcIn
	CompositeSrcOut
	CompositeSrcOver
	CompositeThreshold
	CompositeUndefined
	CompositeVividLight
	CompositeXor
)

// Composite modifies the image, drawing the draw Image argument at offset
// (x, y) using the c Composite operation.
func (im *Image) Composite(c Composite, draw *Image, x int, y int) error {
	op, err := im.compositeOperator(c)
	if err != nil {
		return err
	}
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	var data C.CompositeData
	data.composite = C.int(op)
	data.draw = draw.image
	data.x = C.int(x)
	data.y = C.int(y)
	res, err := im.applyDataFunc("compositing", C.ImageDataFunc(C.compositeImage), &data)
	// res.image will be != than im.image when im is a non
	// coalesced animation
	if res.image != im.image {
		unrefImages(im.image)
		initializeRefCounts(res.image)
		refImages(res.image)
		im.image = res.image
		dontFree(res)
	}
	return err
}

// CompositeInto is equivalent to canvas.Composite/c, im, x, y). See Image.Composite
// for more information.
func (im *Image) CompositeInto(c Composite, canvas *Image, x int, y int) error {
	return canvas.Composite(c, im, x, y)
}
