// +build !gm

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
	case CompositeBlend:
		cc = C.BlendCompositeOp
	case CompositeBlur:
		cc = C.BlurCompositeOp
	case CompositeBumpmap:
		cc = C.BumpmapCompositeOp
	case CompositeChangeMask:
		cc = C.ChangeMaskCompositeOp
	case CompositeClear:
		cc = C.ClearCompositeOp
	case CompositeColorBurn:
		cc = C.ColorBurnCompositeOp
	case CompositeColorDodge:
		cc = C.ColorDodgeCompositeOp
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
	case CompositeDarkenIntensity:
		cc = C.DarkenIntensityCompositeOp
	case CompositeDifference:
		cc = C.DifferenceCompositeOp
	case CompositeDisplace:
		cc = C.DisplaceCompositeOp
	case CompositeDissolve:
		cc = C.DissolveCompositeOp
	case CompositeDistort:
		cc = C.DistortCompositeOp
	case CompositeDivideDst:
		cc = C.DivideDstCompositeOp
	case CompositeDivideSrc:
		cc = C.DivideSrcCompositeOp
	case CompositeDstAtop:
		cc = C.DstAtopCompositeOp
	case CompositeDst:
		cc = C.DstCompositeOp
	case CompositeDstIn:
		cc = C.DstInCompositeOp
	case CompositeDstOut:
		cc = C.DstOutCompositeOp
	case CompositeDstOver:
		cc = C.DstOverCompositeOp
	case CompositeExclusion:
		cc = C.ExclusionCompositeOp
	case CompositeHardLight:
		cc = C.HardLightCompositeOp
	case CompositeHue:
		cc = C.HueCompositeOp
	case CompositeIn:
		cc = C.InCompositeOp
	case CompositeLighten:
		cc = C.LightenCompositeOp
	case CompositeLightenIntensity:
		cc = C.LightenIntensityCompositeOp
	case CompositeLinearBurn:
		cc = C.LinearBurnCompositeOp
	case CompositeLinearDodge:
		cc = C.LinearDodgeCompositeOp
	case CompositeLinearLight:
		cc = C.LinearLightCompositeOp
	case CompositeLuminize:
		cc = C.LuminizeCompositeOp
	case CompositeMathematics:
		cc = C.MathematicsCompositeOp
	case CompositeMinusDst:
		cc = C.MinusDstCompositeOp
	case CompositeMinusSrc:
		cc = C.MinusSrcCompositeOp
	case CompositeModulate:
		cc = C.ModulateCompositeOp
	case CompositeModulusAdd:
		cc = C.ModulusAddCompositeOp
	case CompositeModulusSubtract:
		cc = C.ModulusSubtractCompositeOp
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
	case CompositePegtopLight:
		cc = C.PegtopLightCompositeOp
	case CompositePinLight:
		cc = C.PinLightCompositeOp
	case CompositePlus:
		cc = C.PlusCompositeOp
	case CompositeReplace:
		cc = C.ReplaceCompositeOp
	case CompositeSaturate:
		cc = C.SaturateCompositeOp
	case CompositeScreen:
		cc = C.ScreenCompositeOp
	case CompositeSoftLight:
		cc = C.SoftLightCompositeOp
	case CompositeSrcAtop:
		cc = C.SrcAtopCompositeOp
	case CompositeSrc:
		cc = C.SrcCompositeOp
	case CompositeSrcIn:
		cc = C.SrcInCompositeOp
	case CompositeSrcOut:
		cc = C.SrcOutCompositeOp
	case CompositeSrcOver:
		cc = C.SrcOverCompositeOp
	case CompositeThreshold:
		cc = C.ThresholdCompositeOp
	case CompositeUndefined:
		cc = C.UndefinedCompositeOp
	case CompositeVividLight:
		cc = C.VividLightCompositeOp
	case CompositeXor:
		cc = C.XorCompositeOp
	default:
		return 0, notImplementedError(fmt.Sprintf("composite %s", c))
	}
	return int(cc), nil
}
