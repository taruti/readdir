// +build dragonfly freebsd netbsd openbsd solaris

package readdir

import "syscall"

const sys_FSSTATAT = syscall.SYS_FSTATAT
