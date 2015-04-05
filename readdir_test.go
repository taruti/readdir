package readdir

import (
	"os"
	"sort"
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
	sort.Sort(fis(fis1))
	sort.Sort(fis(fis2))
	for i := 0; i < len(fis1); i++ {
		if fis1[i].Name() != fis2[i].Name() {
			t.Fatal("Sorted names don't match in tmp!")
		}
	}

}

type fis []os.FileInfo

func (a fis) Len() int           { return len(a) }
func (a fis) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a fis) Less(i, j int) bool { return a[i].Name() < a[j].Name() }
