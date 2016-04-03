// +build gm

package magick

// #include <magick/api.h>
import "C"

func GetMagickCopyright() string{
	p := C.GetMagickCopyright()
	s := C.GoString(p)
	return s
}

func GetMagickWebSite() string {
	p := C.GetMagickWebSite()
	s := C.GoString(p)
	return s
}
