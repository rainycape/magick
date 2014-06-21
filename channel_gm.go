// +build gm

package magick

// #include <magick/api.h>
import "C"

import (
	"errors"
)

func (im *Image) toChannel(ch Channel) error {
	if C.ChannelImage(im.image, C.ChannelType(ch)) == 0 {
		return errors.New("error extracting channel")
	}
	return nil
}

func (im *Image) importChannel(src *Image, ch Channel) error {
	C.ImportImageChannel(src.image, im.image, C.ChannelType(ch))
	return nil
}
