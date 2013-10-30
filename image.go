package magick

// #include <stdio.h>
// #include <string.h>
// #include <stdlib.h>
// #include <magick/api.h>
import "C"

import (
	"io"
	"io/ioutil"
	"runtime"
	"sync"
	"unsafe"
)

type Image struct {
	image     *C.Image
	parent    *Image
	coalesced bool
	mu        sync.Mutex
}

/* Image attributes */
func (im *Image) Width() int {
	return int(im.image.columns)
}

func (im *Image) Height() int {
	return int(im.image.rows)
}

func (im *Image) Format() string {
	return C.GoString(&im.image.magick[0])
}

func (im *Image) Depth() int {
	return int(im.image.depth)
}

func (im *Image) Delay() int {
	return int(im.image.delay)
}

func (im *Image) Clone() (*Image, error) {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	var image *C.Image
	if im.parent == nil {
		image = C.CloneImageList(im.image, &ex)
	} else {
		image = C.CloneImage(im.image, magickSize(0), magickSize(0), 1, &ex)
	}
	return checkImage(image, nil, &ex, "cloning")
}

// Dispose frees the resources assocciated with the image.
// If you try to use a disposed image, you'll get undefined
// behavior. Please, note that you don't need to call Dispose
// manually, images will eventually be freed when there are no
// more references to them. However, this function provided in
// case you want to immediately free the memory once you're
// done with an image.
func (im *Image) Dispose() {
	if im.image != nil {
		if im.parent == nil {
			unrefImages(im.image)
		} else {
			unrefImage(im.image)
		}
		im.image = nil
	}
	runtime.SetFinalizer(im, nil)
}

/* Encoding */
func (im *Image) Encode(w io.Writer, info *Info) error {
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	if info == nil {
		info = NewInfo()
	}
	/* ImageToBlob copies the format from the image into
	the image info. Overwrite if required and then restore
	*/
	im.mu.Lock()
	var format *C.char
	copied := false
	if info.info.magick[0] != 0 {
		copied = true
		if im.image.magick[0] != 0 {
			format = C.strdup(&im.image.magick[0])
		}
		C.strncpy(&im.image.magick[0], &info.info.magick[0], C.MaxTextExtent)
	}
	var s C.size_t
	mem := imageToBlob(info, im, &s, &ex)
	if copied {
		/* Restore image format */
		if format != nil {
			C.strncpy(&im.image.magick[0], format, C.MaxTextExtent)
			C.free(unsafe.Pointer(format))
		} else {
			C.memset(unsafe.Pointer(&im.image.magick[0]), 0, C.MaxTextExtent)
		}
	}
	im.mu.Unlock()
	if mem == nil {
		return exError(&ex, "encoding")
	}
	b := C.GoBytes(mem, C.int(s))
	w.Write(b)
	C.free(mem)
	return nil
}

func Decode(r io.Reader) (*Image, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodeData(data)
}

func DecodeData(data []byte) (*Image, error) {
	return decodeData(data, 0)
}

func decodeData(data []byte, try int) (*Image, error) {
	if len(data) == 0 {
		return nil, ErrNoData
	}
	info := C.CloneImageInfo(nil)
	defer C.DestroyImageInfo(info)
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	im := C.BlobToImage(info, unsafe.Pointer(&data[0]), C.size_t(len(data)), &ex)
	if im == nil && try < maxGifTries && ex.severity == C.CorruptImageError && looksLikeGif(data) {
		return fixAndDecodeGif(data, try)
	}
	return checkImage(im, nil, &ex, "decoding")
}

func newImage(im *C.Image, parent *Image) *Image {
	image := new(Image)
	image.image = im
	if parent != nil {
		for parent.parent != nil {
			parent = parent.parent
		}
		image.parent = parent
		refImage(im)
	} else {
		// WARNING: Set the reference count to 0 before calling refImages.
		// Functions hich return an image from another image (e.g. crop, resize, etc...)
		// copy the client_data parameter, which is what we're using for reference
		// counting. Since the image has not been initialized yet, only this
		// goroutine can be accessing it, so we may safely just set all the
		// reference counts to 0.
		for cur := (*C.Image)(im.previous); cur != nil; cur = (*C.Image)(cur.previous) {
			p := (*int32)(unsafe.Pointer(&cur.client_data))
			*p = 0
		}
		for cur := im; cur != nil; cur = (*C.Image)(cur.next) {
			p := (*int32)(unsafe.Pointer(&cur.client_data))
			*p = 0
		}
		refImages(im)
	}
	freeWhenDone(image)
	return image
}

func freeWhenDone(im *Image) {
	runtime.SetFinalizer(im, func(i *Image) {
		i.Dispose()
	})
}

func dontFree(im *Image) {
	runtime.SetFinalizer(im, nil)
}

func checkImage(im *C.Image, parent *Image, ex *C.ExceptionInfo, what string) (*Image, error) {
	if im != nil {
		return newImage(im, parent), nil
	}
	return nil, exError(ex, what)
}
