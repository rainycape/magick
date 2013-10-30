package magick

// #include <magick/api.h>
import "C"

import (
	"errors"
	"fmt"
)

var (
	ErrNoSuchFrame = errors.New("No such frame")
	ErrNoData      = errors.New("No image data")
)

func exError(ex *C.ExceptionInfo, what string) error {
	if what != "" {
		return fmt.Errorf("Error %s image: %s (%s)", what, C.GoString(ex.reason), C.GoString(ex.description))
	}
	return fmt.Errorf("%s: %s", C.GoString(ex.reason), C.GoString(ex.description))
}
