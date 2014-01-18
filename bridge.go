package magick

// #include <magick/api.h>
// #include "bridge.h"
import "C"

import (
	"reflect"
	"unsafe"
)

func btoi(b bool) C.int {
	if b {
		return C.int(1)
	}
	return C.int(0)
}

func (im *Image) root() *C.Image {
	cur := im.image
	for cur.previous != nil {
		cur = (*C.Image)(cur.previous)
	}
	return cur
}

func (im *Image) coalescedCImage(ex *C.ExceptionInfo) (*C.Image, bool) {
	if im.coalesced {
		return im.root(), false
	}
	coalesced := C.CoalesceImages(im.root(), ex)
	if coalesced == nil {
		return nil, false
	}
	return coalesced, true
}

func (im *Image) applyFunc(what string, f C.ImageFunc) (*Image, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	root := im.root()
	res := C.apply_image_func(f, root, unsafe.Pointer(im.parent), btoi(im.coalesced), &ex)
	if res == nil {
		return nil, exError(&ex, what)
	}
	if res == root {
		return im, nil
	}
	ret := newImage(res, nil)
	ret.coalesced = true
	return ret, nil
}

func (im *Image) applyDataFunc(what string, f C.ImageDataFunc, data interface{}) (*Image, error) {
	p := unsafe.Pointer(reflect.ValueOf(data).Pointer())
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	root := im.root()
	res := C.apply_image_data_func(f, root, p, unsafe.Pointer(im.parent), btoi(im.coalesced), &ex)
	if res == nil {
		return nil, exError(&ex, what)
	}
	if res == root {
		return im, nil
	}
	ret := newImage(res, nil)
	ret.coalesced = true
	return ret, nil
}

func (im *Image) applyRectFunc(what string, f C.ImageDataFunc, r Rect) (*Image, error) {
	return im.applyDataFunc(what, f, r.rectangleInfo())
}

func (im *Image) fn(what string, f C.ImageFunc) (*Image, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	res := C.bridge_image_func(f, im.image, &ex)
	return checkImage(res, nil, &ex, what)
}

// ApplyFunc calls the given C.ImageFunc for all the frames in the given
// image and returns the results. The given function only needs to work
// on individual images. ApplyFunc will coalesce the original image and
// create a new sequence of images if required.
func (im *Image) ApplyFunc(f C.ImageFunc) (*Image, error) {
	return im.applyFunc("", f)
}

// ApplyDataFunc calls the given C.ImageDataFunc with the data argument
// (which must be a pointer) for all the frames in the given image and
// returns the results. The given function only needs to work on individual
// images. ApplyDataFunc will coalesce the original image and create a new
// sequence of images if required.
func (im *Image) ApplyDataFunc(f C.ImageDataFunc, data interface{}) (*Image, error) {
	return im.applyDataFunc("", f, data)
}

// ApplyRectFunc calls the given C.ImageDataFunc with Rect argument converted to
// a *C.RectangleInfo for all the frames in the given image and returns the
// results. The given function only needs to work on individual images. ApplyRectFunc
// will coalesce the original image and create a new sequence of
// images if required.
func (im *Image) ApplyRectFunc(f C.ImageDataFunc, r Rect) (*Image, error) {
	return im.applyRectFunc("", f, r)
}

// Func calls the given C.ImageFunc, which expects a sequence of images
// rather than a single frame and results the result as a new Image.
func (im *Image) Func(f C.ImageFunc) (*Image, error) {
	return im.fn("", f)
}
