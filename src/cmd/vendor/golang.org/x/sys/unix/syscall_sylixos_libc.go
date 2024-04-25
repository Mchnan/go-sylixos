// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package unix

import _ "unsafe"

// Implemented in the runtime package (runtime/sys_sylixos3.go)
func syscall_syscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscall_syscallX(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscall_syscallXnull(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscall_syscallXerrno(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscall_syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func syscall_syscall6X(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func syscall_syscall6Xerrno(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func syscall_rawSyscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscall_rawSyscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)

//go:linkname syscall_syscall syscall.syscall
//go:linkname syscall_syscallX syscall.syscallX
//go:linkname syscall_syscallXnull syscall.syscallXnull
//go:linkname syscall_syscallXerrno syscall.syscallXerrno
//go:linkname syscall_syscall6 syscall.syscall6
//go:linkname syscall_syscall6X syscall.rawSyscall6X
//go:linkname syscall_syscall6Xerrno syscall.syscall6Xerrno
//go:linkname syscall_rawSyscall syscall.rawSyscall
//go:linkname syscall_rawSyscall6 syscall.rawSyscall6
