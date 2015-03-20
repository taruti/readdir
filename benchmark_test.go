package readdir

import (
	"os"
	"testing"
)

func BenchmarkOsReaddir(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		f, e := os.Open("/tmp")
		if e != nil {
			b.Fatal("Open failed", e)
		}
		b.StartTimer()
		_, e = f.Readdir(0)
		b.StopTimer()
		if e != nil {
			b.Fatal("Readdir failed", e)
		}
		f.Close()
	}
}

func BenchmarkMyReaddir(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		f, e := os.Open("/tmp")
		if e != nil {
			b.Fatal("Open failed", e)
		}
		b.StartTimer()
		_, e = Readdir(f, 0)
		b.StopTimer()
		if e != nil {
			b.Fatal("Readdir failed", e)
		}
		f.Close()
	}
}

func BenchmarkMyReaddirSys(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		f, e := os.Open("/tmp")
		if e != nil {
			b.Fatal("Open failed", e)
		}
		b.StartTimer()
		_, e = ReaddirSys(f, 0)
		b.StopTimer()
		if e != nil {
			b.Fatal("Readdir failed", e)
		}
		f.Close()
	}
}
