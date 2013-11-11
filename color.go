package magick

// #include <magick/api.h>
// extern float calculate_image_entropy(const Image *image);
// extern float calculate_image_entropy_rect(const Image *image, const RectangleInfo *rect);
import "C"

type HistogramItem struct {
	Color *Pixel
	Count int
}

type Histogram struct {
	Items []*HistogramItem
}

func (im *Image) Histogram() (*Histogram, error) {
	return im.histogram()
}

func (im *Image) Entropy() float32 {
	return float32(C.calculate_image_entropy_rect(im.image, nil))
}

func (im *Image) EntropyRect(r Rect) float32 {
	return float32(C.calculate_image_entropy_rect(im.image, r.rectangleInfo()))
}

func (im *Image) NColors() (int, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	val := C.GetNumberColors(im.image, nil, &ex)
	if ex.severity != C.UndefinedException {
		return 0, exError(&ex, "getting colors")
	}
	return int(val), nil
}

func (im *Image) IsPalette() (bool, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	b := C.IsPaletteImage(im.image, &ex)
	if ex.severity != C.UndefinedException {
		return false, exError(&ex, "getting isPalette")
	}
	return int(b) != 0, nil
}

// AverageColor returns the image's average color
// as a *Pixel.
func (im *Image) AverageColor() (*Pixel, error) {
	stats, err := im.statistics()
	if err != nil {
		return nil, err
	}
	return &Pixel{
		Red:     uint8(stats.Red.Mean * 255),
		Green:   uint8(stats.Green.Mean * 255),
		Blue:    uint8(stats.Blue.Mean * 255),
		Opacity: uint8(stats.Opacity.Mean * 255),
	}, nil
}
