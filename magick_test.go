package magick

import (
	"bytes"
	"os"
	"testing"
)

type Fataler interface {
	Fatal(...interface{})
}

func decodeFile(fat Fataler, name string) *Image {
	f, err := os.Open("test_data/" + name)
	if err != nil {
		fat.Fatal(err)
	}
	defer f.Close()
	im, err := Decode(f)
	if err != nil {
		fat.Fatal(err)
	}
	if im == nil {
		fat.Fatal("No image")
	}
	return im
}

func encodeFile(t *testing.T, name string, im *Image, info *Info) {
	f, err := os.OpenFile("test_data/out."+name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	err = im.Encode(f, info)
	if err != nil {
		t.Fatal(err)
	}
}

func testImage(t *testing.T, im *Image, frames, width, height, depth int, format string) {
	if im.NFrames() != frames {
		t.Errorf("Invalid number of frames, expected %d, got %d", frames, im.NFrames())
	}
	if im.Width() != width || im.Height() != height {
		t.Errorf("Invalid image dimensions, expected %dx%d, got %dx%d", width, height, im.Width(), im.Height())
	}
	if im.Format() != format {
		t.Errorf("Invalid image format, expected %s got %s", format, im.Format())
	}
	if im.Depth() != depth {
		t.Errorf("Invalid depth format, expected %v got %v", depth, im.Depth())
	}
	if im.NFrames() > 1 {
		_, err := im.Coalesce()
		if err != nil {
			t.Errorf("error coalescing: %s", err)
		}
	}
}

func recodeImage(t *testing.T, im *Image, info *Info) *Image {
	buf := &bytes.Buffer{}
	err := im.Encode(buf, info)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := DecodeData(buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	return decoded
}

func testEntropy(t *testing.T, name string, entropy float32) {
	im := decodeFile(t, name)
	if e := im.Entropy(); e != entropy {
		// IM and GM implementations might return slightly different results
		// due to differences in the precision used during calculations
		delta := e - entropy
		if delta > -0.0001 && delta < 0.0001 {
			t.Logf("Slightly different entropy (due to different precision) for %s. Expected %v, got %v (delta %v)", name, entropy, e, delta)
		} else {
			t.Errorf("Invalid entropy for %s: expecting %v, got %v", name, entropy, e)
		}
	}
}

func TestDecode(t *testing.T) {
	t.Logf("Using backend %s", Backend())
	im := decodeFile(t, "wizard.png")
	testImage(t, im, 1, 1104, 1468, 8, "PNG")
	cloned, err := im.Clone()
	if err != nil {
		t.Fatal(err)
	}
	testImage(t, cloned, 1, 1104, 1468, 8, "PNG")
	anim := decodeFile(t, "Newtons_cradle_animation_book_2.gif")
	testImage(t, anim, 36, 480, 360, 8, "GIF")
	im2 := decodeFile(t, "lenna.jpg")
	testImage(t, im2, 1, 512, 512, 8, "JPEG")
}

func TestResize(t *testing.T) {
	im1 := decodeFile(t, "wizard.png")
	res1, err := im1.Resize(500, 600, FQuadratic)
	if err != nil {
		t.Fatal(err)
	}
	testImage(t, res1, 1, 500, 600, 8, "PNG")
	im2 := decodeFile(t, "Newtons_cradle_animation_book_2.gif")
	testImage(t, im2, 36, 480, 360, 8, "GIF")
	res2, err := im2.Resize(240, 180, FQuadratic)
	if err != nil {
		t.Fatal(err)
	}
	testImage(t, res2, 36, 240, 180, 8, "GIF")
}

func TestGif(t *testing.T) {
	im := decodeFile(t, "Newtons_cradle_animation_book_2.gif")
	testImage(t, im, 36, 480, 360, 8, "GIF")
	encodeFile(t, "newton.gif", im, nil)
	cloned, err := im.Clone()
	if err != nil {
		t.Fatal(err)
	}
	testImage(t, cloned, 36, 480, 360, 8, "GIF")
	encodeFile(t, "newton2.gif", cloned, nil)
	frame1, err := im.Frame(0)
	if err != nil {
		t.Fatal(err)
	}
	testImage(t, frame1, 1, 480, 360, 8, "GIF")
	cframe1, err := frame1.Clone()
	if err != nil {
		t.Fatal(err)
	}
	testImage(t, cframe1, 1, 480, 360, 8, "GIF")
}

func TestEncode(t *testing.T) {
	im1 := decodeFile(t, "wizard.png")
	res1, err := im1.Resize(500, 600, FQuadratic)
	if err != nil {
		t.Fatal(err)
	}
	im2 := recodeImage(t, res1, nil)
	testImage(t, im2, 1, 500, 600, 8, "PNG")
	info := NewInfo()
	info.SetFormat("JPEG")
	im3 := recodeImage(t, res1, info)
	testImage(t, im3, 1, 500, 600, 8, "JPEG")

	gif1 := decodeFile(t, "Newtons_cradle_animation_book_2.gif")
	gif2 := recodeImage(t, gif1, nil)
	if gif1.Duration() != gif2.Duration() {
		t.Errorf("Invalid duration, expected %d, got %d", gif1.Duration(), gif2.Duration())
	}
	testImage(t, gif2, 36, 480, 360, 8, "GIF")
	gif3, err := gif2.Resize(240, 180, FQuadratic)
	if err != nil {
		t.Fatal(err)
	}
	gif4 := recodeImage(t, gif3, nil)
	testImage(t, gif4, 36, 240, 180, 8, "GIF")
	encodeFile(t, "newton-small.gif", gif4, nil)
	if gif1.Duration() != gif4.Duration() {
		t.Errorf("Invalid duration, expected %d, got %d", gif1.Duration(), gif4.Duration())
	}
	info.SetFormat("PNG")
	nongif := recodeImage(t, gif3, info)
	testImage(t, nongif, 1, 240, 180, 8, "PNG")
	encodeFile(t, "newton.png", nongif, nil)
}

func TestList(t *testing.T) {
	list1 := decodeFile(t, "Newtons_cradle_animation_book_2.gif")
	nframes := list1.NFrames()
	list1.Append(list1)
	if list1.NFrames() != 2*nframes {
		t.Errorf("Error appending self, expected %d frames, got %d %p", 2*nframes, list1.NFrames(), list1)
	}
	nframes = list1.NFrames()
	img1 := decodeFile(t, "wizard.png")
	list1.Append(img1)
	if list1.NFrames() != nframes+1 {
		t.Errorf("Error appending single, expected %d frames, got %d", nframes+1, list1.NFrames())
	}
	nframes = list1.NFrames()
	list1.RemoveFirst()
	if list1.NFrames() != nframes-1 {
		t.Errorf("Error removing first: expected %d frames, got %d", nframes-1, list1.NFrames())
	}
	nframes = list1.NFrames()
	list1.RemoveLast()
	if list1.NFrames() != nframes-1 {
		t.Errorf("Error removing last: expected %d frames, got %d", nframes-1, list1.NFrames())
	}
	nframes = list1.NFrames()
	frame, err := list1.Frame(5)
	if err != nil {
		t.Fatal(err)
	}
	frame.Remove()
	if list1.NFrames() != nframes-1 {
		t.Errorf("Error removing specific frame: expected %d frames, got %d", nframes-1, list1.NFrames())
	}
}

func TestEntropy(t *testing.T) {
	testEntropy(t, "wizard.png", 5.073119)
	testEntropy(t, "lenna.jpg", 8.774539)
}

func TestDecodeOptimized(t *testing.T) {
	im := decodeFile(t, "optimized.gif")
	testImage(t, im, 10, 651, 721, 8, "GIF")
	// This one requires gifsicle with --unoptimize
	im2 := decodeFile(t, "math.gif")
	testImage(t, im2, 158, 500, 350, 8, "GIF")
	// This one requires gifsicle with --unoptimize and piping via convert if using GM
	im3 := decodeFile(t, "kick_grandma.gif")
	testImage(t, im3, 25, 240, 180, 8, "GIF")
}

func TestCropResizeAnimated(t *testing.T) {
	im := decodeFile(t, "optimized.gif")
	res, err := im.CropResize(100, 100, FHamming, CSCenter)
	if err != nil {
		t.Fatal(err)
	}
	testImage(t, res, 10, 100, 100, 8, "GIF")
}

func TestAverage(t *testing.T) {
	im := decodeFile(t, "lenna.jpg")
	avg, err := im.AverageColor()
	if err != nil {
		t.Fatal(err)
	}
	if avg.Red != 133 || avg.Green != 80 || avg.Blue != 68 {
		t.Errorf("expected (133, 80, 68), got %+v instead", *avg)
	}
}

func BenchmarkRefUnref(b *testing.B) {
	im := decodeFile(b, "wizard.png")
	img := im.image
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		refImage(img)
		unrefImage(img)
	}
}

func BenchmarkResizeAnimated(b *testing.B) {
	im := decodeFile(b, "Newtons_cradle_animation_book_2.gif")
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		resized, err := im.Resize(240, 180, FQuadratic)
		if err != nil {
			b.Fatal(err)
		}
		resized.Dispose()
	}
}

func BenchmarkMinifyAnimated(b *testing.B) {
	im := decodeFile(b, "Newtons_cradle_animation_book_2.gif")
	coalesced, err := im.Coalesce()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		_, err = coalesced.Minify()
		if err != nil {
			b.Fatal(err)
		}
	}
}
