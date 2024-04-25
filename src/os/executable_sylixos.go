// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package os

import (
	"internal/itoa"
)

func executable() (string, error) {
	procfn := "/proc/" + itoa.Itoa(Getpid()) + "/exe"
	path, err := Readlink(procfn)

	// When the executable has been deleted then Readlink returns a
	// path appended with " (deleted)".
	return stringsTrimSuffix(path, " (deleted)"), err
}

// stringsTrimSuffix is the same as strings.TrimSuffix.
func stringsTrimSuffix(s, suffix string) string {
	if len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix {
		return s[:len(s)-len(suffix)]
	}
	return s
}
