// +build !gm

package magick

// #include <magick/api.h>
import "C"

import (
	"errors"
)

func (im *Image) toChannel(ch Channel) error {
	if C.SeparateImageChannel(im.image, C.ChannelType(ch)) == 0 {
		return errors.New("could no extract channel")
	}
	return nil
}

func (im *Image) importChannel(src *Image, ch Channel) error {
	return notImplementedError("ImportChannel")
}
