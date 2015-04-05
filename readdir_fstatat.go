// +build freebsd linux

package readdir

import (
	"os"
	"syscall"
	"unsafe"
)

func lstatUnder(fd int, name string, st *syscall.Stat_t) error {
	nameptr, err := syscall.BytePtrFromString(name)
	if err != nil {
		return err
	}
	_, _, e1 := syscall.Syscall6(sys_FSSTATAT,
		uintptr(fd),
		uintptr(unsafe.Pointer(nameptr)),
		uintptr(unsafe.Pointer(st)),
		uintptr(at_SYMLINK_NOFOLLOW),
		0, 0)
	if e1 != 0 {
		return e1
	}
	return nil
}

type NamedStat struct {
	FileName string
	FileMode os.FileMode
	syscall.Stat_t
}

func ReaddirSys(f *os.File, max int) ([]NamedStat, error) {
	ns, e := f.Readdirnames(max)
	if e != nil {
		return nil, e
	}
	ss := make([]NamedStat, len(ns))
	fd := int(f.Fd())
	idx, skipped := 0, 0
	for _, n := range ns {
		e = lstatUnder(fd, n, &ss[idx].Stat_t)
		if e != nil {
			if os.IsNotExist(e) {
				skipped++
				continue
			}
			return nil, e
		}
		ss[idx].FileName = n
		ss[idx].FileMode = mode2FM(ss[idx].Stat_t.Mode)
		idx++
	}
	ss = ss[0 : len(ss)-skipped]
	return ss, nil
}

func Readdir(f *os.File, max int) ([]os.FileInfo, error) {
	ss, e := ReaddirSys(f, max)
	if e != nil {
		return nil, e
	}
	fis := make([]os.FileInfo, len(ss))
	for i := range fis {
		fis[i] = &ss[i]
	}
	return fis, nil
}
