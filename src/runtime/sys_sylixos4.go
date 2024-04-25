// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package runtime

import (
	"internal/abi"
	"unsafe"
)

//go:nosplit
//go:cgo_unsafe_args
func sem_init(sem *semt, pshared int32, value uint32) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(sem_init_trampoline)), unsafe.Pointer(&sem))
	KeepAlive(sem)
	return ret
}
func sem_init_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func sem_post(sem *semt) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(sem_post_trampoline)), unsafe.Pointer(&sem))
	KeepAlive(sem)
	return ret
}
func sem_post_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func sem_reltimedwait_np(sem *semt, timeout *timespec) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(sem_reltimedwait_np_trampoline)), unsafe.Pointer(&sem))
	KeepAlive(sem)
	return ret
}
func sem_reltimedwait_np_trampoline()

//go:nosplit
//go:cgo_unsafe_args
func sem_wait(sem *semt) int32 {
	ret := libcCall(unsafe.Pointer(abi.FuncPCABI0(sem_wait_trampoline)), unsafe.Pointer(&sem))
	KeepAlive(sem)
	return ret
}
func sem_wait_trampoline()

//go:cgo_import_dynamic libc_sem_init sem_init "libfastlock.so"
//go:cgo_import_dynamic libc_sem_post sem_post "libfastlock.so"
//go:cgo_import_dynamic libc_sem_reltimedwait_np sem_reltimedwait_np "libfastlock.so"
//go:cgo_import_dynamic libc_sem_wait sem_wait "libfastlock.so"

//go:cgo_import_dynamic _ _ "libfastlock.so"
//go:cgo_import_dynamic _ _ "libvpmpdm.so"
