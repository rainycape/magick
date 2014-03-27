package magick

// #include <stdlib.h>
// extern void * gif_encode(void *im, int single, int *size);
// #cgo LDFLAGS: -lgif
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"unsafe"
)

var (
	gifsicleCmd          string
	convertCmd           string
	errNoConvert         = errors.New("error decoding GIF image: Corrupt data. Install imagemagick (convert) to try to fix it.")
	errNoGifsicle        = errors.New("error decoding GIF image: Corrupt data. Install gifsicle to try to fix it.")
	errCouldNotEncodeGif = errors.New("error encoding GIF")
	maxGifTries          = 2
)

func looksLikeGif(data []byte) bool {
	return bytes.HasPrefix(data, []byte{'G', 'I', 'F'})
}

func runGifsicle(data []byte, args []string) ([]byte, error) {
	cmd := exec.Command(gifsicleCmd, args...)
	cmd.Stdin = bytes.NewReader(data)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// Workaround bug in gifsicle 1.71
		if strings.Contains(err.Error(), "segmentation fault") {
			return runGifsicle(data, append(args, "--colors=256"))
		}
		return nil, fmt.Errorf("error running gifsicle: %s", err)
	}
	return out.Bytes(), nil
}

func fixAndDecodeGif(data []byte, try int) (*Image, error) {
	if gifsicleCmd == "" {
		return nil, errNoGifsicle
	}
	args := []string{"--careful"}
	if try > 0 {
		args = append(args, "--unoptimize")
	}
	data, err := runGifsicle(data, args)
	if err != nil {
		return nil, err
	}
	if try > 1 {
		if convertCmd == "" {
			return nil, errNoConvert
		}
		cmd := exec.Command(convertCmd, "-", "-")
		cmd.Stdin = bytes.NewReader(data)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return nil, fmt.Errorf("error running convert: %s", err)
		}
		data = out.Bytes()
	}
	return decodeData(data, try+1)
}

// GifEncode encodes the Image as GIF using giflib
// rather than ImageMagick or GraphicksMagick. While this
// will result in a lower quality GIF, encoding is 9-10x
// faster and produces files ~20% smaller.
func (im *Image) GifEncode() ([]byte, error) {
	var size C.int
	var single C.int
	if im.parent != nil {
		single = 1
	}
	res := C.gif_encode(unsafe.Pointer(im.image), single, &size)
	if res != nil {
		p := unsafe.Pointer(res)
		b := C.GoBytes(p, size)
		C.free(p)
		return b, nil
	}
	return nil, errCouldNotEncodeGif
}

func init() {
	gifsicleCmd, _ = exec.LookPath("gifsicle")
	if Backend() == "GraphicsMagick" {
		maxGifTries = 3
		convertCmd, _ = exec.LookPath("convert")
	}
}
