package readdir

import (
	"os"
	"testing"
)

func TestReaddir(t *testing.T) {
	f, e := os.Open("/tmp")
	if e != nil {
		t.Fatal("Open failed", e)
	}
	fis1, e := f.Readdir(0)
	if e != nil {
		t.Fatal("os Readdir failed", e)
	}
	f.Close()
	f, e = os.Open("/tmp")
	if e != nil {
		t.Fatal("Open failed", e)
	}
	fis2, e := Readdir(f, 0)
	if e != nil {
		t.Fatal("os Readdir failed", e)
	}
	f.Close()
	if len(fis1) != len(fis2) {
		t.Fatal("Readdir lengths differ for /tmp:", len(fis1), len(fis2))
	}

}
