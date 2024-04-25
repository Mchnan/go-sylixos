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
func osyield() {
	libcCall(unsafe.Pointer(abi.FuncPCABI0(sched_yield_trampoline)), unsafe.Pointer(nil))
}
func sched_yield_trampoline()

//go:nosplit
func osyield_no_g() {
	asmcgocall_no_g(unsafe.Pointer(abi.FuncPCABI0(sched_yield_trampoline)), unsafe.Pointer(nil))
}

//go:cgo_import_dynamic libc_sched_yield sched_yield "libvpmpdm.so"
