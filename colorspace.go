package magick

// #include <magick/api.h>
import "C"

// When encoding an output image, the colorspaces RGB, CMYK, and GRAY
// may be specified. The CMYK option is only applicable when writing
// TIFF, JPEG, and Adobe Photoshop bitmap (PSD) files.
type Colorspace C.ColorspaceType

const (
	RGB           Colorspace = C.RGBColorspace         // Red, Green, Blue colorspace.
	GRAY          Colorspace = C.GRAYColorspace        // Similar to Luma (Y) according to ITU-R 601
	TRANSPARENT   Colorspace = C.TransparentColorspace // RGB which preserves the matte while quantizing colors.
	OHTA          Colorspace = C.OHTAColorspace
	XYZ           Colorspace = C.XYZColorspace // CIE XYZ
	YCC           Colorspace = C.YCCColorspace // Kodak PhotoCD PhotoYCC
	YIQ           Colorspace = C.YIQColorspace
	YPBPR         Colorspace = C.YPbPrColorspace
	YUV           Colorspace = C.YUVColorspace         // YUV colorspace as used for computer video.
	CMYK          Colorspace = C.CMYKColorspace        // Cyan, Magenta, Yellow, Black colorspace.
	SRGB          Colorspace = C.sRGBColorspace        // Kodak PhotoCD sRGB
	HSL           Colorspace = C.HSLColorspace         // Hue, saturation, luminosity
	HWB           Colorspace = C.HWBColorspace         // Hue, whiteness, blackness
	LAB           Colorspace = C.LABColorspace         // ITU LAB
	REC_601_LUMA  Colorspace = C.Rec601LumaColorspace  // Luma (Y) according to ITU-R 601
	REC_601_YCBCR Colorspace = C.Rec601YCbCrColorspace // YCbCr according to ITU-R 601
	REC_709_LUMA  Colorspace = C.Rec709LumaColorspace  // Luma (Y) according to ITU-R 709
	REC_709_YCBCR Colorspace = C.Rec709YCbCrColorspace // YCbCr according to ITU-R 709
)

func (c Colorspace) String() string {
	switch c {
	case RGB:
		return "RGB"
	case GRAY:
		return "GRAY"
	case TRANSPARENT:
		return "TRANSPARENT"
	case OHTA:
		return "OHTA"
	case XYZ:
		return "XYZ"
	case YCC:
		return "YCC"
	case YIQ:
		return "YIQ"
	case YPBPR:
		return "YPBPR"
	case YUV:
		return "YUV"
	case CMYK:
		return "CMYK"
	case SRGB:
		return "SRGB"
	case HSL:
		return "HSL"
	case HWB:
		return "HWB"
	case LAB:
		return "LAB"
	case REC_601_LUMA:
		return "REC_601_LUMA"
	case REC_601_YCBCR:
		return "REC_601_YCBCR"
	case REC_709_LUMA:
		return "REC_709_LUMA"
	case REC_709_YCBCR:
		return "REC_709_YCBCR"
	}
	return "UNDEFINED"
}

// Colorspace returns the image colorspace.
func (im *Image) Colorspace() Colorspace {
	return Colorspace(im.image.colorspace)
}

// SetColorspace changes the image colorspace. Note
// that this only changes how the pixels are interpreted.
// If you want to transform the image to another colorspace
// use TransformColorspace().
func (im *Image) SetColorspace(cs Colorspace) {
	im.image.colorspace = C.ColorspaceType(cs)
}

// TransformColorspace changes the image colorspace while also
// changing the pixels to represent the same image in the
// new colorspace.
func (im *Image) TransformColorspace(cs Colorspace) (*Image, error) {
	return im.transformColorspace(cs)
}
