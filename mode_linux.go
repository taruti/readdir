package readdir

import (
	"os"
	"syscall"
	"time"
)

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
