package magick

// #include <magick/api.h>
import "C"

func (im *Image) IsOrphan() bool {
	return im.image.next == nil && im.image.previous == nil
}

func (im *Image) NFrames() int {
	if im.parent != nil {
		return 1
	}
	return int(C.GetImageListLength(im.image))
}

func (im *Image) FrameIndex() int {
	return int(C.GetImageIndexInList(im.image))
}

func (im *Image) Duration() int {
	if im.parent != nil {
		return im.Delay()
	}
	duration := 0
	for cur := C.GetFirstImageInList(im.image); cur != nil; cur = (*C.Image)(cur.next) {
		duration += int(cur.delay)
	}
	return duration
}

func (im *Image) Frame(idx int) (*Image, error) {
	if im.IsOrphan() && idx == 0 {
		return im, nil
	}
	image := C.GetImageFromList(im.image, magickLong(idx))
	if image == nil {
		return nil, ErrNoSuchFrame
	}
	return newImage(image, im), nil
}

func (im *Image) First() *Image {
	first, _ := im.Frame(0)
	return first
}

func (im *Image) Last() *Image {
	last, _ := im.Frame(im.NFrames() - 1)
	return last
}

func (im *Image) Prev() *Image {
	prev := C.GetPreviousImageInList(im.image)
	if prev == nil {
		return nil
	}
	return newImage(prev, im)
}

func (im *Image) Next() *Image {
	next := C.GetNextImageInList(im.image)
	if next == nil {
		return nil
	}
	return newImage(next, im)
}

func (im *Image) Reverse() {
	C.ReverseImageList(&im.image)
}

/*func (im *Image) InsertAfter(image *Image) error {
	cloned, err := image.Clone()
	if err != nil {
		return err
	}
	first := cloned.First()
	first.prev = im
	if im.next != nil {
		im.next.prev = cloned.Last()
	}
	im.next = first
	return nil
}

func (im *Image) InsertBefore(image *Image) error {
	cloned, err := image.CloneAll()
	if err != nil {
		return err
	}
	last := cloned.Last()
	last.next = im
	if im.prev != nil {
		im.prev.next = cloned.First()
	}
	image.prev = last
	return nil
}
*/

func (im *Image) Prepend(image *Image) error {
	// TODO: Don't clone, refcount
	cloned, err := image.Clone()
	if err != nil {
		return err
	}
	dontFree(cloned)
	C.PrependImageToList(&im.image, cloned.image)
	return nil
}

func (im *Image) Append(image *Image) error {
	// TODO: Don't clone, refcount
	cloned, err := image.Clone()
	if err != nil {
		return err
	}
	dontFree(cloned)
	C.AppendImageToList(&im.image, cloned.image)
	return nil
}

func (im *Image) RemoveIndex(idx int) bool {
	if frame, _ := im.Frame(idx); frame != nil {
		return frame.Remove()
	}
	return false
}

func (im *Image) RemoveFirst() bool {
	return im.RemoveIndex(0)
}

func (im *Image) RemoveLast() bool {
	return im.RemoveIndex(im.NFrames() - 1)
}

func (im *Image) Remove() bool {
	if im.parent != nil && !im.parent.IsOrphan() {
		// Don't use DeleteImageFromList, since it calls DestroyImage,
		// bypassing our refcounting.
		if im.image == im.parent.image {
			im.parent.image = (*C.Image)(im.parent.image.next)
		}
		if p := (*C.Image)(im.image.previous); p != nil {
			p.next = im.image.next
		}
		if n := (*C.Image)(im.image.next); n != nil {
			n.previous = im.image.previous
		}
		im.image.previous = nil
		im.image.next = nil
		im.parent = nil
		// im already holds a reference, remove the reference that
		// the parent had.
		unrefImage(im.image)
		return true
	}
	return false
}

func (im *Image) Apply(f func(*Image) (*Image, error)) (*Image, error) {
	if im.IsOrphan() || im.parent != nil {
		return f(im)
	}
	var ex C.ExceptionInfo
	C.GetExceptionInfo(&ex)
	defer C.DestroyExceptionInfo(&ex)
	cur, free := im.coalescedCImage(&ex)
	if cur == nil {
		return nil, exError(&ex, "coalescing")
	}
	if free {
		defer C.DestroyImageList(cur)
	}
	buf := &Image{
		parent: im,
		image:  cur,
	}
	first, err := f(buf)
	if err != nil {
		return nil, err
	}
	prev := first.image
	for cur := (*C.Image)(im.image.next); cur != nil; cur = (*C.Image)(cur.next) {
		buf.image = cur
		res, err := f(buf)
		if err != nil {
			return nil, err
		}
		dontFree(res)
		img := res.image
		prev.next = (*C.struct__Image)(img)
		img.previous = (*C.struct__Image)(img)
		prev = img
	}
	return first, nil
}
