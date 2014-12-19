// +build !gm

package magick

// #include "pixel.h"
import "C"

func (im *Image) pixels(r *Rect, ex *C.ExceptionInfo) *C.PixelPacket {
	return C.GetVirtualPixels(im.image, C.ssize_t(r.X), C.ssize_t(r.Y), C.size_t(r.Width), C.size_t(r.Height), ex)
}

func (im *Image) setPixels(r *Rect, src *C.PixelPacket, ex *C.ExceptionInfo) bool {
	dst := C.GetAuthenticPixels(im.image, C.ssize_t(r.X), C.ssize_t(r.Y), C.size_t(r.Width), C.size_t(r.Height), ex)
	if dst == nil {
		return false
	}
	C.copy_pixel_packets(src, dst, C.int(r.Width*r.Height))
	if C.SyncAuthenticPixels(im.image, ex) != C.MagickTrue {
		return false
	}
	return true
}
