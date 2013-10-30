// +build debug

package magick

import (
	"runtime"
	"testing"
	"time"
)

func TestGC(t *testing.T) {
	t.Logf("%v live images before GC()", live)
	for ii := 0; ii < 1000; ii++ {
		runtime.GC()
	}
	time.Sleep(5 * time.Second)
	for ii := 0; ii < 1000; ii++ {
		runtime.GC()
	}
	if live != 0 {
		t.Errorf("%v images still alive", live)
		for k := range liveImages {
			t.Errorf("image %p has reference count %d", k, refCount(k))
		}
	}
}
