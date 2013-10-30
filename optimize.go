package magick

// #include <magick/api.h>
import "C"

import ()

type Optimizer func([]byte) ([]byte, error)
