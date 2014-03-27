package magick

// #include <stdlib.h>
// #include <magick/api.h>
import "C"

import (
	"unsafe"
)

// Properties returns the keys of the properties that
// the image contains.
func (im *Image) Properties() []string {
	return im.properties()
}

// DestroyPropertys removes all properties from the image.
func (im *Image) DestroyProperties() {
	im.destroyProperties()
}

// Property returns the property value assocciated with
// the given key. Note that both non-present keys and
// keys set to "" will return an empty string. To check
// if an image has a property defined use Image.HasProperty.
func (im *Image) Property(key string) string {
	prop := im.property(key)
	if prop != nil {
		return *prop
	}
	return ""
}

// HasProperty returns wheter the image has a property
// with the given key, even if it's set to an empty string.
func (im *Image) HasProperty(key string) bool {
	return im.property(key) != nil
}

// SetProperty adds a new property to the image. If the property
// already exists, it's overwritten. Returns true
// if the property could be added to the list, false otherwise.
// For removing a property, see Image.RemoveProperty.
func (im *Image) SetProperty(key string, value string) bool {
	val := C.CString(value)
	ret := im.setProperty(key, val)
	C.free(unsafe.Pointer(val))
	return ret
}

// RemoveProperty removes the property specified by key.
func (im *Image) RemoveProperty(key string) bool {
	return im.setProperty(key, nil)
}
