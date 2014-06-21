package magick

// #include <magick/api.h>
import "C"

// ErrorStats represents the difference between
// two compared images.
type ErrorStats struct {
	MeanPerPixel      float64
	NormalizedMean    float64
	NormalizedMaximum float64
}

// IsZero returns true iff the two compared images
// are equal.
func (stats *ErrorStats) IsZero() bool {
	return stats != nil && stats.MeanPerPixel == 0
}

// Compare returns comparison of two images as an ErrorStatus.
// See also IsEqual.
func (im *Image) Compare(src *Image) (*ErrorStats, error) {
	// Note: IM 7.x might return an error from IsImagesEqual,
	// so the error value will not be always nil in the
	// future
	im.mu.Lock()
	defer im.mu.Unlock()
	// Save image.error value, since IsImagesEqual alters it
	e := im.image.error
	C.IsImagesEqual(im.image, src.image)
	stats := &ErrorStats{
		MeanPerPixel:      float64(im.image.error.mean_error_per_pixel),
		NormalizedMean:    float64(im.image.error.normalized_mean_error),
		NormalizedMaximum: float64(im.image.error.normalized_maximum_error),
	}
	im.image.error = e
	return stats, nil
}

// IsEqual returns wheter two images are equal. See also
// Compare.
func (im *Image) IsEqual(src *Image) (bool, error) {
	stats, err := im.Compare(src)
	return stats.IsZero(), err
}
