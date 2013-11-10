package magick

// #include <magick/api.h>
import "C"

type Channel C.ChannelType

const (
	CRed     Channel = C.RedChannel     /* RGB Red channel */
	CCyan            = C.CyanChannel    /* CMYK Cyan channel */
	CGreen           = C.GreenChannel   /* RGB Green channel */
	CMagenta         = C.MagentaChannel /* CMYK Magenta channel */
	CBlue            = C.BlueChannel    /* RGB Blue channel */
	CYellow          = C.YellowChannel  /* CMYK Yellow channel */
	COpacity         = C.OpacityChannel /* Opacity channel */
	CBlack           = C.BlackChannel   /* CMYK Black (K) channel */
	CAll             = C.AllChannels    /* Color channels */
	CGray            = C.GrayChannel    /* Color channels represent an intensity. */
)
