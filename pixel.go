package magick

// #include <magick/api.h>
// extern unsigned char quantum_to_char(Quantum q);
// extern Quantum char_to_quantum(unsigned char c);
import "C"

// Pixel represents a color identified by its
// red, green, blue and alpha components. Their
// value goes from 0 to 255.
type Pixel struct {
	Red     uint8
	Green   uint8
	Blue    uint8
	Opacity uint8
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
