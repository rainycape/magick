package magick

// #include <magick/api.h>
import "C"

// Backend returns the name of the library backend
// which was selected at build time. It must be
// either "ImageMagick" or "GraphicsMagick".
func Backend() string {
	return backend
}

// Supported formats returns a list with the names
// of all supported image formats. This varies depending
// on the backend and the compile options that have been
// used while building IM or GM.
func SupportedFormats() ([]string, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	infos, p := supportedFormats(&ex)
	if infos == nil {
		return nil, exError(&ex, "getting supported formats")
	}
	var formats []string
	for _, v := range infos {
		if v == nil {
			break
		}
		formats = append(formats, C.GoString(v.name))
	}
	freeMagickMemory(p)
	return formats, nil
}
