package magick

import (
	"bytes"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	wizard = "test_data/wizard.png"
	newton = "test_data/Newtons_cradle_animation_book_2.gif"
	lenna  = "test_data/lenna.jpg"
)

func BenchmarkResizePng(b *testing.B) {
	im := decodeFile(b, "wizard.png")
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		_, _ = im.Resize(240, 180, FLanczos)
	}
}

func BenchmarkResizePngGo(b *testing.B) {
	f, err := os.Open(wizard)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	im, _, err := image.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		_ = resize.Resize(240, 180, im, resize.Lanczos2Lut)
	}
}

func benchmarkDecode(b *testing.B, file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		if _, err := DecodeData(data); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkDecodeGo(b *testing.B, file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		b.Fatal(err)
	}
	r := bytes.NewReader(data)
	if filepath.Ext(file) == ".gif" {
		b.ResetTimer()
		for ii := 0; ii < b.N; ii++ {
			if _, err := gif.DecodeAll(r); err != nil {
				b.Fatal(err)
			}
			r.Seek(0, 0)
		}
	} else {
		b.ResetTimer()
		for ii := 0; ii < b.N; ii++ {
			if _, _, err := image.Decode(r); err != nil {
				b.Fatal(err)
			}
			r.Seek(0, 0)
		}
	}
}

func benchmarkEncode(b *testing.B, file string, format string, quality uint) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		b.Fatal(err)
	}
	im, err := DecodeData(data)
	if err != nil {
		b.Fatal(err)
	}
	info := NewInfo()
	info.SetFormat(format)
	info.SetQuality(quality)
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		if err := im.Encode(ioutil.Discard, info); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkEncodeGo(b *testing.B, file string, format string, quality uint) {
	f, err := os.Open(file)
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	im, _, err := image.Decode(f)
	if err != nil {
		b.Fatal(err)
	}
	switch format {
	case "jpeg":
		opts := &jpeg.Options{
			Quality: int(quality),
		}
		b.ResetTimer()
		for ii := 0; ii < b.N; ii++ {
			if err := jpeg.Encode(ioutil.Discard, im, opts); err != nil {
				b.Fatal(err)
			}
		}
	case "png":
		b.ResetTimer()
		for ii := 0; ii < b.N; ii++ {
			if err := png.Encode(ioutil.Discard, im); err != nil {
				b.Fatal(err)
			}
		}
	default:
		b.Fatalf("format %s is not supported", format)
	}
}

func BenchmarkDecodePng(b *testing.B) {
	benchmarkDecode(b, wizard)
}

func BenchmarkDecodePngGo(b *testing.B) {
	benchmarkDecodeGo(b, wizard)
}

func BenchmarkDecodeGif(b *testing.B) {
	benchmarkDecode(b, newton)
}

func BenchmarkDecodeGifGo(b *testing.B) {
	benchmarkDecodeGo(b, newton)
}

func BenchmarkDecodeJpeg(b *testing.B) {
	benchmarkDecode(b, lenna)
}

func BenchmarkDecodeJpegGo(b *testing.B) {
	benchmarkDecodeGo(b, lenna)
}

func BenchmarkEncodePng(b *testing.B) {
	benchmarkEncode(b, lenna, "png", 0)
}

func BenchmarkEncodePngGo(b *testing.B) {
	benchmarkEncodeGo(b, lenna, "png", 0)
}

func BenchmarkEncodeJpeg(b *testing.B) {
	benchmarkEncode(b, wizard, "jpeg", 60)
}

func BenchmarkEncodeJpegGo(b *testing.B) {
	benchmarkEncodeGo(b, wizard, "jpeg", 60)
}

func BenchmarkEncodeJpegHQ(b *testing.B) {
	benchmarkEncode(b, wizard, "jpeg", 100)
}

func BenchmarkEncodeJpegHQGo(b *testing.B) {
	benchmarkEncodeGo(b, wizard, "jpeg", 100)
}

func BenchmarkEncodeJpegLQ(b *testing.B) {
	benchmarkEncode(b, wizard, "jpeg", 10)
}

func BenchmarkEncodeJpegLQGo(b *testing.B) {
	benchmarkEncodeGo(b, wizard, "jpeg", 10)
}
