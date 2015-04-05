package readdir

import "syscall"

const at_SYMLINK_NOFOLLOW = 0x100
const sys_FSSTATAT = syscall.SYS_NEWFSTATAT
