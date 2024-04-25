// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos && arm64

#include "textflag.h"
#include "cgo/abi_arm64.h"

TEXT _rt0_arm64_sylixos(SB),NOSPLIT|NOFRAME,$0
	MOVD	$runtime·rt0_go(SB), R3
	BL	(R3)
	MOVW	$-1, R0
	CALL	libc_exit(SB)

TEXT main(SB),NOSPLIT|NOFRAME,$0
	MOVD	$runtime·rt0_go(SB), R3
	BL	(R3)
	MOVW	$-1, R0
	CALL	libc_exit(SB)
