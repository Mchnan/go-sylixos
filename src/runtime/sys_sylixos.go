// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package runtime

import (
	"internal/abi"
	"unsafe"
)

// The *_trampoline functions convert from the Go calling convention to the C calling convention
// and then call the underlying libc function. These are defined in sys_sylixos_$ARCH.s.

//go:nosplit
//go:cgo_unsafe_args
func pthread_attr_init(attr *pthreadattr) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(pthread_attr_init_trampoline)), unsafe.Pointer(&attr))
	KeepAlive(attr)
	return ret
}
func pthread_attr_init_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_attr_destroy(attr *pthreadattr) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(pthread_attr_destroy_trampoline)), unsafe.Pointer(&attr))
	KeepAlive(attr)
	return ret
}
func pthread_attr_destroy_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_attr_setdetachstate(attr *pthreadattr, state int) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(pthread_attr_setdetachstate_trampoline)), unsafe.Pointer(&attr))
	KeepAlive(attr)
	return ret
}
func pthread_attr_setdetachstate_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_create(attr *pthreadattr, start uintptr, arg unsafe.Pointer) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(pthread_create_trampoline)), unsafe.Pointer(&attr))
	KeepAlive(attr)
	KeepAlive(arg) // Just for consistency. Arg of course needs to be kept alive for the start function.
	return ret
}
func pthread_create_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_kill(thread pthread, signo int) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(pthread_kill_trampoline)), unsafe.Pointer(&thread))
	return ret
}
func pthread_kill_trampoline()

//go:nosplit
func pthread_self() int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(pthread_self_trampoline)), unsafe.Pointer(nil))
	return ret
}
func pthread_self_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func pthread_attr_setstack(attr *pthreadattr, addr unsafe.Pointer, size uintptr) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(pthread_attr_setstack_trampoline)), unsafe.Pointer(&attr))
	KeepAlive(attr)
	KeepAlive(addr) // Just for consistency. addr of course needs to be kept alive for the start function.
	return ret
}
func pthread_attr_setstack_trampoline()

// Tell the linker that the libc_* functions are to be found
// in a system library, with the libc_ prefix missing.

//go:cgo_import_dynamic libc_pthread_attr_init pthread_attr_init "libvpmpdm.so"
//go:cgo_import_dynamic libc_pthread_attr_destroy pthread_attr_destroy "libvpmpdm.so"
//go:cgo_import_dynamic libc_pthread_attr_setdetachstate pthread_attr_setdetachstate "libvpmpdm.so"
//go:cgo_import_dynamic libc_pthread_create pthread_create "libvpmpdm.so"
//go:cgo_import_dynamic libc_pthread_sigmask pthread_sigmask "libvpmpdm.so"
//go:cgo_import_dynamic libc_pthread_kill pthread_kill "libvpmpdm.so"
//go:cgo_import_dynamic libc_pthread_self pthread_self "libvpmpdm.so"
//go:cgo_import_dynamic libc_pthread_attr_setstack pthread_attr_setstack "libvpmpdm.so"
