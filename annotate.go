// +build gm

package magick

// #include <string.h>
// #include <magick/api.h>
import "C"
import "unsafe"

func (im *Image) AnnotateImage(x, y int, text string, fontSize int, color string){
	image_info := C.CloneImageInfo(nil)
	defer C.DestroyImageInfo(image_info);

	draw_info := C.CloneDrawInfo(image_info, nil)
	defer C.DestroyDrawInfo(draw_info)

	textBytes := []byte(text)
	bytesPointer := unsafe.Pointer(&textBytes[0])

	textArray := (*C.uchar)(bytesPointer)

	colorArray := C.CString(color)
	defer C.free(unsafe.Pointer(colorArray))

	draw_context := C.DrawAllocateContext(draw_info, im.image)
	defer C.DrawDestroyContext(draw_context)

	C.DrawSetFontSize(draw_context, C.double(fontSize))
	C.DrawSetStrokeColorString(draw_context, colorArray)
	C.DrawSetFillColorString(draw_context, colorArray)
	C.DrawAnnotation(draw_context, C.double(x), C.double(y), textArray)
	C.DrawRender(draw_context)
}
