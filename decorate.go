package magick

// #include <magick/api.h>
import "C"

func (im *Image) AddBorder(r Rect, c *Pixel) (*Image, error) {
	return im.Apply(func(i *Image) (*Image, error) {
		var ex C.ExceptionInfo
		C.GetExceptionInfo(&ex)
		defer C.DestroyExceptionInfo(&ex)
		var b *Pixel
		if c != nil {
			b = i.BorderColor()
			im.SetBorderColor(c)
		}
		bordered := C.BorderImage(i.image, r.rectangleInfo(), &ex)
		if b != nil {
			i.SetBorderColor(b)
		}
		if bordered == nil {
			return nil, exError(&ex, "adding border")
		}
		return newImage(bordered, nil), nil
	})
}

func (im *Image) AddFrame(f Frame, c *Pixel) (*Image, error) {
	return im.Apply(func(i *Image) (*Image, error) {
		var ex C.ExceptionInfo
		C.GetExceptionInfo(&ex)
		defer C.DestroyExceptionInfo(&ex)
		var m *Pixel
		if c != nil {
			m = i.MatteColor()
			i.SetMatteColor(c)
		}
		framed := C.FrameImage(i.image, f.frameInfo(), &ex)
		if m != nil {
			i.SetMatteColor(m)
		}
		if framed == nil {
			return nil, exError(&ex, "adding border")
		}
		return newImage(framed, nil), nil
	})
}

func (im *Image) Raise(r Rect, threedimensional bool) (*Image, error) {
	// RaiseImage modifies the image, use a clone
	cloned, err := im.Clone()
	if err != nil {
		return nil, err
	}
	flag := magickBool(threedimensional)
	rect := r.rectangleInfo()
	return cloned.Apply(func(i *Image) (*Image, error) {
		C.RaiseImage(i.image, rect, flag)
		return i, nil
	})
}
