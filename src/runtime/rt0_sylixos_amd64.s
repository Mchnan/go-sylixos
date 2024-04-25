// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos && amd64

#include "textflag.h"

TEXT _rt0_amd64_sylixos(SB),NOSPLIT,$-8
	JMP	runtime·rt0_go(SB)

TEXT _rt0_amd64_sylixos_lib(SB),NOSPLIT,$0
	JMP	_rt0_amd64_lib(SB)
