// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package unix

import (
	"syscall"
	"unsafe"
)

func Ioctl(fd int, cmd int, args unsafe.Pointer) (err error) {
	_, _, err = syscall.RawSyscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(cmd), uintptr(args))
	return
}
