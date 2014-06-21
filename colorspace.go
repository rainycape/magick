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

// TransformColorspace returns a new image by covnerting the original to
// the given colorspace while also changing the pixels to represent the
// same image in the new colorspace.
func (im *Image) TransformColorspace(cs Colorspace) (*Image, error) {
	clone, err := im.Clone()
	if err != nil {
		return nil, err
	}
	if err := clone.ToColorspace(cs); err != nil {
		return nil, err
	}
	return clone, nil
}

// ToColorspace changes the image colorspace in place.
func (im *Image) ToColorspace(cs Colorspace) error {
	_, err := im.transformColorspace(cs)
	return err
}
