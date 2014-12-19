// +build gm

package magick

// #include "pixel.h"
import "C"

func (im *Image) pixels(r *Rect, ex *C.ExceptionInfo) *C.PixelPacket {
	return C.GetImagePixelsEx(im.image, C.long(r.X), C.long(r.Y), C.ulong(r.Width), C.ulong(r.Height), ex)
}

func (im *Image) setPixels(r *Rect, src *C.PixelPacket, ex *C.ExceptionInfo) bool {
	dst := C.SetImagePixelsEx(im.image, C.long(r.X), C.long(r.Y), C.ulong(r.Width), C.ulong(r.Height), ex)
	if dst == nil {
		return false
	}
	C.copy_pixel_packets(src, dst, C.int(r.Width*r.Height))
	if C.SyncImagePixelsEx(im.image, ex) != C.MagickPass {
		return false
	}
	return true
}
