// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package runtime

import (
	"internal/abi"
	"unsafe"
)

type pollfd struct {
	fd      int32
	events  int16
	revents int16
}

const (
	_POLLIN  = 0x0003
	_POLLOUT = 0x0008
	_POLLHUP = 0x0040
	_POLLERR = 0x0020

	_LW_OPTION_VUTEX_FLAG_WAKEALL  = 0x0001
	_LW_OPTION_VUTEX_FLAG_DONTSET  = 0x0002
	_LW_OPTION_VUTEX_FLAG_DEEPWAKE = 0x0004

	_LW_VPROC_EXIT_NORMAL = 0
	_LW_VPROC_EXIT_FORCE  = 1
)

//go:nosplit
//go:cgo_unsafe_args
func poll(pfds *pollfd, npfds uintptr, timeout uintptr) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(poll_trampoline)), unsafe.Pointer(&pfds))
	KeepAlive(pfds)
	return ret
}
func poll_trampoline()

func EpollCreate1(flags int32) (fd int32, errno uintptr) {
	r1, _, e := syscall_rawSyscall(abi.FuncPCABI0(libc_epoll_create1_trampoline), uintptr(flags), 0, 0)
	return int32(r1), e
}
func libc_epoll_create1_trampoline()

var _zero uintptr

func EpollWait(epfd int32, events []EpollEvent, maxev, waitms int32) (n int32, errno uintptr) {
	var ev unsafe.Pointer
	if len(events) > 0 {
		ev = unsafe.Pointer(&events[0])
	} else {
		ev = unsafe.Pointer(&_zero)
	}
	r1, _, e := syscall_rawSyscall6(abi.FuncPCABI0(libc_epoll_wait_trampoline), uintptr(epfd), uintptr(ev), uintptr(len(events)), uintptr(waitms), 0, 0)
	KeepAlive(events)
	return int32(r1), e
}
func libc_epoll_wait_trampoline()

func EpollCtl(epfd, op, fd int32, event *EpollEvent) (errno uintptr) {
	_, _, e := syscall_rawSyscall6(abi.FuncPCABI0(libc_epoll_ctl_trampoline), uintptr(epfd), uintptr(op), uintptr(fd), uintptr(unsafe.Pointer(event)), 0, 0)
	KeepAlive(event)
	return e
}
func libc_epoll_ctl_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func API_VutexPend(addr *uint32, desired uint32, timeout uintptr) {
	libcCall(unsafe.Pointer(abi.FuncPCABI0(API_VutexPend_trampoline)), unsafe.Pointer(&addr))
}
func API_VutexPend_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func API_VutexPostEx(addr *uint32, value uint32, flags uint32) {
	libcCall(unsafe.Pointer(abi.FuncPCABI0(API_VutexPostEx_trampoline)), unsafe.Pointer(&addr))
}
func API_VutexPostEx_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func vprocExitModeSet(mode uint32) {
	libcCall(unsafe.Pointer(abi.FuncPCABI0(vprocExitModeSet_trampoline)), unsafe.Pointer(&mode))
}
func vprocExitModeSet_trampoline()

//go:cgo_import_dynamic libc_epoll_create1 epoll_create1 "libvpmpdm.so"
//go:cgo_import_dynamic libc_epoll_ctl epoll_ctl "libvpmpdm.so"
//go:cgo_import_dynamic libc_epoll_wait epoll_wait "libvpmpdm.so"

//go:cgo_import_dynamic libc_poll poll "libvpmpdm.so"

//go:cgo_import_dynamic libc_API_VutexPend API_VutexPend "libvpmpdm.so"
//go:cgo_import_dynamic libc_API_VutexPostEx API_VutexPostEx "libvpmpdm.so"

//go:cgo_import_dynamic libc_vprocExitModeSet vprocExitModeSet "libvpmpdm.so"
