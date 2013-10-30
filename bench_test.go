package magick

import (
	"github.com/nfnt/resize"
	"image"
	_ "image/png"
	"os"
	"testing"
)

func BenchmarkResizePng(b *testing.B) {
	im := decodeFile(b, "wizard.png")
	b.ResetTimer()
	for ii := 0; ii < b.N; ii++ {
		_, _ = im.Resize(240, 180, FLanczos)
	}
}

func BenchmarkResizePngNative(b *testing.B) {
	f, err := os.Open("test_data/wizard.png")
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
		_ = resize.Resize(240, 180, im, resize.Lanczos2)
	}
}
