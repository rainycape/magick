// +build gm

package magick

// #include <magick/api.h>
import "C"

func (im *Image) statistics() (*Statistics, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	var stats C.ImageStatistics
	C.GetImageStatistics(im.image, &stats, &ex)
	if ex.severity != C.UndefinedException {
		return nil, exError(&ex, "getting statistics")
	}
	return newStatistics(&stats), nil
}

func newChannelStatistics(ch *C.ImageChannelStatistics) *ChannelStatistics {
	return &ChannelStatistics{
		Minimum:  float64(ch.minimum),
		Maximum:  float64(ch.maximum),
		Mean:     float64(ch.mean),
		StdDev:   float64(ch.standard_deviation),
		Variance: float64(ch.variance),
	}
}

func newStatistics(stats *C.ImageStatistics) *Statistics {
	red := newChannelStatistics(&stats.red)
	green := newChannelStatistics(&stats.green)
	blue := newChannelStatistics(&stats.blue)
	opacity := newChannelStatistics(&stats.opacity)
	return &Statistics{red, green, blue, opacity}
}
