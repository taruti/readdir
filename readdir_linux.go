package readdir

import (
	"os"
	"syscall"
	"time"
	"unsafe"
)

const at_SYMLINK_NOFOLLOW = 0x100

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

// modelled after Golang src/os/stat_linux.go
func mode2FM(mode uint32) os.FileMode {
	var fs os.FileMode = os.FileMode(mode & 0777)
	switch mode & syscall.S_IFMT {
	case syscall.S_IFBLK:
		fs |= os.ModeDevice
	case syscall.S_IFCHR:
		fs |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFDIR:
		fs |= os.ModeDir
	case syscall.S_IFIFO:
		fs |= os.ModeNamedPipe
	case syscall.S_IFLNK:
		fs |= os.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		fs |= os.ModeSocket
	}
	if mode&syscall.S_ISGID != 0 {
		fs |= os.ModeSetgid
	}
	if mode&syscall.S_ISUID != 0 {
		fs |= os.ModeSetuid
	}
	if mode&syscall.S_ISVTX != 0 {
		fs |= os.ModeSticky
	}
	return fs
}

func (fs *NamedStat) Size() int64        { return int64(fs.Stat_t.Size) }
func (fs *NamedStat) Mode() os.FileMode  { return fs.FileMode }
func (fs *NamedStat) ModTime() time.Time { return time.Unix(int64(fs.Mtim.Sec), int64(fs.Mtim.Nsec)) }
func (fs *NamedStat) Sys() interface{}   { return fs }
func (fs *NamedStat) Name() string       { return fs.FileName }
func (fs *NamedStat) IsDir() bool        { return fs.FileMode.IsDir() }
