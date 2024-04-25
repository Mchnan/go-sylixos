// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package pprof

import (
	"bytes"
	"errors"
	"internal/itoa"
	"os"
	"strconv"
	"syscall"
)

func parseProcSelfModules(data []byte, addMapping func(lo, hi, offset uint64, file, buildID string)) {
	// [root@sylixos:/proc/2]# cat modules
	// NAME HANDLE TYPE GLB BASE SIZE SYMCNT
	// pprof.test 045c5e90 USER YES 100040000 35acf0 0
	// libfastlock.so 045baad0 USER YES 1003a0000 13f08 50
	// libvpmpdm.so 045c6540 USER YES 1003c0000 292c0 315
	// <VP Ver:2.2.0 pt-malloc>

	var line []byte
	// next removes and returns the next field in the line.
	// It also removes from line any spaces following the field.
	next := func() []byte {
		var f []byte
		f, line, _ = bytes.Cut(line, space)
		line = bytes.TrimLeft(line, " ")
		return f
	}

	exe, _ := os.Executable()
	exeLen := len(exe)

	for len(data) > 0 {
		line, data, _ = bytes.Cut(data, newline)
		name := next()

		next() // handle
		next() // type
		next() // global

		baseStr := next()
		lo, err := strconv.ParseUint(string(baseStr), 16, 64)
		if err != nil {
			continue
		}

		sizeStr := next()
		size, err := strconv.ParseUint(string(sizeStr), 16, 64)
		if err != nil {
			continue
		}

		next() // symcnt

		file := string(name)
		if len(file) <= exeLen && exe[exeLen-len(file):] == file {
			file = exe
		}

		buildID := peBuildID(file)
		addMapping(lo, lo+size, 0, file, buildID)
	}
}

// readMapping reads /proc/self/modules and writes mappings to b.pb.
// It saves the address ranges of the mappings in b.mem for use
// when emitting locations.
func (b *profileBuilder) readMapping() {
	data, _ := os.ReadFile("/proc/" + itoa.Itoa(syscall.Getpid()) + "/modules")
	parseProcSelfModules(data, b.addMapping)
	if len(b.mem) == 0 { // pprof expects a map entry, so fake one.
		b.addMappingEntry(0, 0, 0, "", "", true)
		// TODO(hyangah): make addMapping return *memMap or
		// take a memMap struct, and get rid of addMappingEntry
		// that takes a bunch of positional arguments.
	}
}

func readMainModuleMapping() (start, end uint64, err error) {
	data, _ := os.ReadFile("/proc/" + itoa.Itoa(syscall.Getpid()) + "/modules")

	var line []byte
	// next removes and returns the next field in the line.
	// It also removes from line any spaces following the field.
	next := func() []byte {
		var f []byte
		f, line, _ = bytes.Cut(line, space)
		line = bytes.TrimLeft(line, " ")
		return f
	}

	for len(data) > 0 {
		line, data, _ = bytes.Cut(data, newline)
		next() // name
		next() // handle
		next() // type
		next() // global

		baseStr := next()
		lo, err := strconv.ParseUint(string(baseStr), 16, 64)
		if err != nil {
			continue
		}

		sizeStr := next()
		size, err := strconv.ParseUint(string(sizeStr), 16, 64)
		if err != nil {
			continue
		}

		next() // symcnt

		return lo, lo + size, nil
	}

	return 0, 0, errors.New("read /proc/self/modules failed")
}
