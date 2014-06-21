package magick

// #include <magick/api.h>
import "C"

import (
	"errors"
)

type Channel C.ChannelType

const (
	CRed     Channel = C.RedChannel     /* RGB Red channel */
	CCyan    Channel = C.CyanChannel    /* CMYK Cyan channel */
	CGreen   Channel = C.GreenChannel   /* RGB Green channel */
	CMagenta Channel = C.MagentaChannel /* CMYK Magenta channel */
	CBlue    Channel = C.BlueChannel    /* RGB Blue channel */
	CYellow  Channel = C.YellowChannel  /* CMYK Yellow channel */
	COpacity Channel = C.OpacityChannel /* Opacity channel */
	CBlack   Channel = C.BlackChannel   /* CMYK Black (K) channel */
	CAll     Channel = C.AllChannels    /* Color channels */
	CGray    Channel = C.GrayChannel    /* Color channels represent an intensity. */
)

// ChannelImage returns a new image with all the channels equal to the
// given channel. e.g. on an image with all pixels with red = 255,
// im.ChannelImage(magick.CRed) will result in a white image, since green and blue
// will be set to the red value. See also ToChannelImage.
func (im *Image) ChannelImage(ch Channel) (*Image, error) {
	clone, err := im.Clone()
	if err != nil {
		return nil, err
	}
	if err := clone.ToChannelImage(ch); err != nil {
		return nil, err
	}
	return clone, nil
}

// ToChannelImage works like ChannelImage, but modifies the image in place.
func (im *Image) ToChannelImage(ch Channel) error {
	return im.toChannel(ch)
}

// ImportChannel imports the given channel from the src image.
func (im *Image) ImportChannel(src *Image, ch Channel) error {
	return im.importChannel(src, ch)
}

// ChannelDepth returns the depth of the given channel.
func (im *Image) ChannelDepth(ch Channel) (uint, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	ret := C.GetImageChannelDepth(im.image, C.ChannelType(ch), &ex)
	if ex.severity != C.UndefinedException {
		return 0, exError(&ex, "getting channel info")
	}
	return uint(ret), nil
}

// SetChannelDepth sets the depth of the channel. The range
// of depth is 1 to QuantumDepth.
func (im *Image) SetChannelDepth(ch Channel, depth uint) error {
	if C.SetImageChannelDepth(im.image, C.ChannelType(ch), magickUint(depth)) == 0 {
		return errors.New("error setting channel")
	}
	return nil
}
