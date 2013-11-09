package magick

// #include <magick/api.h>
import "C"

// ChannelStatistics includes several statistics about
// a color channel. Magick normalizes All fields in this
// structure to the interval [0, 1].
type ChannelStatistics struct {
	Minimum  float64
	Maximum  float64
	Mean     float64
	StdDev   float64
	Variance float64
	// The following fields are only filled in when building
	// against ImageMagick
	Kurtosis float64
	Skewness float64
}

type Statistics struct {
	Red     *ChannelStatistics
	Green   *ChannelStatistics
	Blue    *ChannelStatistics
	Opacity *ChannelStatistics
}

func (im *Image) Statistics() (*Statistics, error) {
	return im.statistics()
}
