package magick

// #include "pixel.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// Pixel represents a color identified by its
// red, green, blue and alpha components. Their
// value goes from 0 to 255.
type Pixel struct {
	Red     uint8
	Green   uint8
	Blue    uint8
	Opacity uint8 // 0 means fully opaque, 255 means fully transparent
}

func (p *Pixel) String() string {
	return fmt.Sprintf("(R:%d G:%d B:%d O:%d)", p.Red, p.Green, p.Blue, p.Opacity)
}

func newPixel(px *C.PixelPacket) *Pixel {
	return &Pixel{
		uint8(C.quantum_to_char(px.red)),
		uint8(C.quantum_to_char(px.green)),
		uint8(C.quantum_to_char(px.blue)),
		uint8(C.quantum_to_char(px.opacity)),
	}
}

func copyPixel(px *Pixel, pkt *C.PixelPacket) {
	pkt.red = C.char_to_quantum(C.uchar(px.Red))
	pkt.green = C.char_to_quantum(C.uchar(px.Green))
	pkt.blue = C.char_to_quantum(C.uchar(px.Blue))
	pkt.opacity = C.char_to_quantum(C.uchar(px.Opacity))
}

// BackgroundColor returns the image background color.
func (im *Image) BackgroundColor() *Pixel {
	return newPixel(&im.image.background_color)
}

// SetBackgroundColor changes the image background color.
func (im *Image) SetBackgroundColor(px *Pixel) {
	copyPixel(px, &im.image.background_color)
}

// BorderColore returns the image border color.
func (im *Image) BorderColor() *Pixel {
	return newPixel(&im.image.border_color)
}

// SetBorderColor changes the image border color.
func (im *Image) SetBorderColor(px *Pixel) {
	copyPixel(px, &im.image.border_color)
}

// MatteColor returns the image matter color.
func (im *Image) MatteColor() *Pixel {
	return newPixel(&im.image.matte_color)
}

// SetMatteColor changes the image matte color.
func (im *Image) SetMatteColor(px *Pixel) {
	copyPixel(px, &im.image.matte_color)
}

// Pixel returns a pixel from the image at the given (x, y).
// coordinates. Note that modifications to the returned Pixel
// won't alter the image. Use SetPixel to change a pixel in the
// image.
func (im *Image) Pixel(x int, y int) (*Pixel, error) {
	px, err := im.Pixels(Rect{X: x, Y: y, Width: 1, Height: 1})
	if err != nil {
		return nil, err
	}
	return px[0], nil
}

// SetPixel changes the pixel at the given (x, y) coordinates.
func (im *Image) SetPixel(x int, y int, p *Pixel) error {
	return im.SetPixels(Rect{X: x, Y: y, Width: 1, Height: 1}, []*Pixel{p})
}

// Pixels returns the image pixels contained in the given rect,
// in row major order. Note that modifications to the returned pixels
// won't alter the image. Use SetPixels to change pixels in the
// image.
func (im *Image) Pixels(r Rect) ([]*Pixel, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	ptr := im.pixels(&r, &ex)
	if ptr == nil {
		return nil, exError(&ex, "getting pixels")
	}
	count := int(r.Width * r.Height)
	var pkts []C.PixelPacket
	header := (*reflect.SliceHeader)((unsafe.Pointer(&pkts)))
	header.Cap = count
	header.Len = count
	header.Data = uintptr(unsafe.Pointer(ptr))
	px := make([]*Pixel, count)
	for ii, v := range pkts {
		px[ii] = newPixel(&v)
	}
	return px, nil
}

// SetPixels changes the pixels in the given rect, in row major order.
// Note that the number of elements in p must match r.Width * r.Height,
// otherwise an error is returned.
func (im *Image) SetPixels(r Rect, p []*Pixel) error {
	if int(r.Width*r.Height) != len(p) {
		return fmt.Errorf("rect contains %d pixels, but %d were provided", int(r.Width*r.Height), len(p))
	}
	if len(p) == 0 {
		return nil
	}
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	packets := make([]C.PixelPacket, len(p))
	for ii, v := range p {
		copyPixel(v, &packets[ii])
	}
	if !im.setPixels(&r, (*C.PixelPacket)(unsafe.Pointer(&packets[0])), &ex) {
		return exError(&ex, "setting pixels")
	}
	return nil
}
