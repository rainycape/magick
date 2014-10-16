// +build gm

package magick

// #include <magick/api.h>
import "C"

import (
	"fmt"
)

func (im *Image) compositeOperator(c Composite) (int, error) {
	var cc C.CompositeOperator
	switch c {
	case CompositeAtop:
		cc = C.AtopCompositeOp
	case CompositeBumpmap:
		cc = C.BumpmapCompositeOp
	case CompositeClear:
		cc = C.ClearCompositeOp
	case CompositeColorize:
		cc = C.ColorizeCompositeOp
	case CompositeCopyBlack:
		cc = C.CopyBlackCompositeOp
	case CompositeCopyBlue:
		cc = C.CopyBlueCompositeOp
	case CompositeCopy:
		cc = C.CopyCompositeOp
	case CompositeCopyCyan:
		cc = C.CopyCyanCompositeOp
	case CompositeCopyGreen:
		cc = C.CopyGreenCompositeOp
	case CompositeCopyMagenta:
		cc = C.CopyMagentaCompositeOp
	case CompositeCopyRed:
		cc = C.CopyRedCompositeOp
	case CompositeCopyYellow:
		cc = C.CopyYellowCompositeOp
	case CompositeDarken:
		cc = C.DarkenCompositeOp
	case CompositeDifference:
		cc = C.DifferenceCompositeOp
	case CompositeDisplace:
		cc = C.DisplaceCompositeOp
	case CompositeDissolve:
		cc = C.DissolveCompositeOp
	case CompositeHue:
		cc = C.HueCompositeOp
	case CompositeIn:
		cc = C.InCompositeOp
	case CompositeLuminize:
		cc = C.LuminizeCompositeOp
	case CompositeModulate:
		cc = C.ModulateCompositeOp
	case CompositeMultiply:
		cc = C.MultiplyCompositeOp
	case CompositeNo:
		cc = C.NoCompositeOp
	case CompositeOut:
		cc = C.OutCompositeOp
	case CompositeOver:
		cc = C.OverCompositeOp
	case CompositeOverlay:
		cc = C.OverlayCompositeOp
	case CompositePlus:
		cc = C.PlusCompositeOp
	case CompositeSaturate:
		cc = C.SaturateCompositeOp
	case CompositeScreen:
		cc = C.ScreenCompositeOp
	case CompositeThreshold:
		cc = C.ThresholdCompositeOp
	case CompositeUndefined:
		cc = C.UndefinedCompositeOp
	case CompositeXor:
		cc = C.XorCompositeOp
	case CompositeAlpha:
		cc = C.CopyOpacityCompositeOp
	default:
		return 0, notImplementedError(fmt.Sprintf("composite %s", c))
	}
	return int(cc), nil
}
