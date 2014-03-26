package magick

// #include <stdio.h>
// #include <string.h>
// #include <stdlib.h>
// #include <magick/api.h>
import "C"

import (
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"sync"
	"unsafe"
)

// Image represents an in-memory decoded image.
// Some images (like e.g. GIF animations) might have
// multiple frames. Unless indicated otherwise, functions
// work in the same way with both single a multi-frame images,
// performing coalescing when needed. To operate on a single
// frame use e.g. im.Frame(i).Resize() instead of im.Resize().
// Frames might be added or removed with list related functions,
// like Append() or Remove().
type Image struct {
	image     *C.Image
	parent    *Image
	coalesced bool
	mu        sync.Mutex
}

// Width returns the image width in pixels.
func (im *Image) Width() int {
	return int(im.image.columns)
}

// Height returns the image height in pixels.
func (im *Image) Height() int {
	return int(im.image.rows)
}

// Rect is a conveniency function which returns a Rect
// at (0, 0) with the image dimensions.
func (im *Image) Rect() Rect {
	return Rect{0, 0, uint(im.Width()), uint(im.Height())}
}

// Format returns the format used to decode
// this image.
func (im *Image) Format() string {
	return C.GoString(&im.image.magick[0])
}

// Depth returns the pixel depth of the image,
// usually 8.
func (im *Image) Depth() int {
	return int(im.image.depth)
}

// Delay returns the time this image stays visible in
// an animation, in 1/100ths of a second. If this image
// is not part of a sequence of animated images, it returns 0.
func (im *Image) Delay() int {
	return int(im.image.delay)
}

// Clone returns a copy of the image. If the image
// has multiple frames, it copies all of them. To
// Clone just one frame use im.Frame(i).Clone().
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

// Dispose frees the memory associated with the image.
// If you try to use a disposed image, you'll get undefined
// behavior. Note that you don't usually need to call
// Dispose manually. Just before an Image is collected by the GC,
// its Dispose method will be called for you. However, if you're
// allocating multiple images in a loop, it's probably better to
// manually Dispose them as soon as you don't need them anymore,
// to avoid the temporary memory usage from getting too high.
// Behind the scenes, Image uses a finalizer to call Dispose. Please,
// see http://golang.org/pkg/runtime/#SetFinalizer for more
// information about finalizers.
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

// Encode writes the image to the given io.Writer, encoding
// it according to the info parameter. Please, see the
// Info type for the available encoding options.
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
	b := goBytes(mem, int(s))
	w.Write(b)
	C.free(mem)
	return nil
}

// Image returns the underlying *C.Image. This is useful for
// calling GM or IM directly and performing operations which
// are not yet supported by magick.
func (im *Image) Image() *C.Image {
	return im.image
}

// Decode tries to decode an image from the given io.Reader.
// If the image can't be decoded or it's corrupt, an error
// will be returned. Depending on the backend and compile
// time options, the number of supported formats might
// vary. Use SupportedFormats() to list the all.
func Decode(r io.Reader) (*Image, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return DecodeData(data)
}

// DecodeData works like Decode, but accepts a
// []byte rather than an io.Reader.
func DecodeData(data []byte) (*Image, error) {
	return decodeData(data, 0)
}

// DecodeFile works like Decode, but accepts a
// filename.
func DecodeFile(filename string) (*Image, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Decode(f)
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

func initializeRefCounts(im *C.Image) {
	for cur := (*C.Image)(im.previous); cur != nil; cur = (*C.Image)(cur.previous) {
		p := (*int32)(unsafe.Pointer(&cur.client_data))
		*p = 0
	}
	for cur := im; cur != nil; cur = (*C.Image)(cur.next) {
		p := (*int32)(unsafe.Pointer(&cur.client_data))
		*p = 0
	}
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
		// Functions which return an image from another image (e.g. crop, resize, etc...)
		// copy the client_data parameter, which is what we're using for reference
		// counting. Since the image has not been initialized yet, only this
		// goroutine can be accessing it, so we may safely just set all the
		// reference counts to 0.
		initializeRefCounts(im)
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
