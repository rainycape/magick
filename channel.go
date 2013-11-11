package magick

// #include <magick/api.h>
import "C"

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
