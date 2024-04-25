// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package syscall

type SysProcAttr struct {
	Chroot     string      // Chroot.
	Credential *Credential // Credential.
	Ptrace     bool        // Enable tracing.
	Setsid     bool        // Create session.
	Setpgid    bool        // Set process group ID to Pgid, or, if Pgid == 0, to new pid.
	Setctty    bool        // Set controlling terminal to fd Ctty
	Noctty     bool        // Detach fd 0 from controlling terminal
	Ctty       int         // Controlling TTY fd
	Foreground bool        // Place child's process group in foreground. (Implies Setpgid. Uses Ctty as fd of controlling TTY)
	Pgid       int         // Child's process group ID if Setpgid.
}

// Flags to be set in the `posix_spawnattr_t'.
const (
	_POSIX_SPAWN_RESETIDS      = 0x01
	_POSIX_SPAWN_SETPGROUP     = 0x02
	_POSIX_SPAWN_SETSIGDEF     = 0x04
	_POSIX_SPAWN_SETSIGMASK    = 0x08
	_POSIX_SPAWN_SETSCHEDPARAM = 0x10
	_POSIX_SPAWN_SETSCHEDULER  = 0x20
)

// Fork, dup fd onto 0..len(fd), and exec(argv0, argvv, envv) in child.
// If a dup or exec fails, write the errno error to pipe.
// (Pipe is close-on-exec so if exec succeeds, it will be closed.)
// In the child, this function must not acquire any locks, because
// they might have been locked at the time of the fork. This means
// no rescheduling, no malloc calls, and no new stack segments.
// For the same reason compiler does not race instrument it.
// The calls to RawSyscall are okay because they are assembly
// functions that do not grow the stack.
//
//go:norace
func forkAndExecInChild(argv0 *byte, argv, envv []*byte, chroot, dir *byte, attr *ProcAttr, sys *SysProcAttr, pipe int) (pid int, err Errno) {
	var spawnattr posix_spawnattr_t
	var actions posix_spawn_file_actions_t
	var maxFd = -1
	var closeFds []int

	fdsLen := len(attr.Files)
	fds := make([]int, fdsLen)
	for i, ufd := range attr.Files {
		fds[i] = int(ufd)
		if fds[i] > maxFd {
			maxFd = fds[i]
		}
	}
	closeexec := make([]int, fdsLen)

	if maxFd > -1 {
		closeFds = make([]int, maxFd+1)
	}

	posix_spawnattr_init(&spawnattr)

	if dir != nil {
		posix_spawnattr_setwd(&spawnattr, dir)
	}

	posix_spawn_file_actions_init(&actions)

	for i := 0; i < fdsLen; i++ {
		if fds[i] < fdsLen {
			if fds[fds[i]] != fds[i] {
				return -1, EINVAL
			}
		}

		closeexec[i], _ = fcntl(fds[i], F_GETFD, closeexec[i])
		fcntl(fds[i], F_SETFD, 0)
		if fds[i] != i {
			posix_spawn_file_actions_adddup2(&actions, fds[i], i)
		}
	}

	for i := 0; i < fdsLen; i++ {
		if fds[i] != i {
			if closeexec[i] != 0 {
				if closeFds[fds[i]] == 0 {
					posix_spawn_file_actions_addclose(&actions, fds[i])
					closeFds[fds[i]] = 1
				}
			}
		}
	}

	// Set process group
	if sys.Setpgid {
		posix_spawnattr_setpgroup(&spawnattr, sys.Pgid)
		posix_spawnattr_setflags(&spawnattr, _POSIX_SPAWN_SETPGROUP)
	}

	err = posix_spawn(&pid, argv0, &actions, &spawnattr, argv, envv)

	posix_spawn_file_actions_destroy(&actions)
	posix_spawnattr_destroy(&spawnattr)

	for i := 0; i < fdsLen; i++ {
		fcntl(fds[i], F_SETFD, closeexec[i])
	}

	return
}
