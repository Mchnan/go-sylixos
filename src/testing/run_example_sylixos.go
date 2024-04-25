// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package testing

import (
	"fmt"
	"io"
	"os"
	"time"
)

func runExample(eg InternalExample) (ok bool) {
	if chatty.on {
		fmt.Printf("%s=== RUN   %s\n", chatty.prefix(), eg.Name)
	}

	// Capture stdout.
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Stdout = w
	outC := make(chan string)
	finished := false

	go func() {
		var allOutput string
		for {
			buf, err := io.ReadAll(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "testing: copying pipe: %v\n", err)
				os.Exit(1)
			}
			if len(buf) > 0 {
				allOutput += string(buf)
			} else if finished {
				break
			}
		}
		outC <- allOutput
		r.Close()
	}()

	start := time.Now()

	// Clean up in a deferred call so we can recover if the example panics.
	defer func() {
		timeSpent := time.Since(start)

		// Close pipe, restore stdout, get output.
		w.Close()
		os.Stdout = stdout
		out := <-outC

		err := recover()
		ok = eg.processRunResult(out, timeSpent, finished, err)
	}()

	// Run example.
	eg.F()
	finished = true
	return
}
