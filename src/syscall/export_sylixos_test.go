// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package syscall

import (
	"internal/abi"
)

func Ioctl(fd, req, arg uintptr) (err Errno) {
	_, _, err = rawSyscall(abi.FuncPCABI0(libc_ioctl_trampoline), fd, req, arg)
	return err
}
