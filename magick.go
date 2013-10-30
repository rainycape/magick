package magick

// Backend returns the name of the library backend
// which was selected at build time. It must be
// either "ImageMagick" or "GraphicsMagick".
func Backend() string {
	return backend
}
