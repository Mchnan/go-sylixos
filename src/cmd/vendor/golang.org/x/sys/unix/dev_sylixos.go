// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Functions to access/create device major and minor numbers matching the
// encoding used in sylixos's system/ioLib/ioInterface.c.

package unix

// sylixos not use major and minor
// Major returns the major component of a sylixos device number.
func Major(dev uint64) uint32 {
	return uint32(dev >> 20)
}

// sylixos not use major and minor
// Minor returns the minor component of a sylixos device number.
func Minor(dev uint64) uint32 {
	return uint32(dev & 0xfffff)
}

// sylixos not use Mkdev
// Mkdev returns a sylixos device number generated from the given major and minor
// components.
func Mkdev(major, minor uint32) uint64 {
	return (uint64(major) << 20) | uint64(minor)
}
