// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// System calls and other sys.stuff for arm64, SylixOS
// System calls are implemented in libvpmpdm, this file
// contains trampolines that convert from Go to C calling convention.
// Some direct system call implementations currently remain.
//

//go:build sylixos && arm64

#include "go_asm.h"
#include "go_tls.h"
#include "textflag.h"
#include "cgo/abi_arm64.h"

#define CLOCK_REALTIME	$0
#define	CLOCK_MONOTONIC	$1

// mstart_stub is the first function executed on a new thread started by pthread_create.
// It just does some low-level setup and then calls mstart.
// Note: called with the C calling convention.
TEXT runtime·mstart_stub(SB),NOSPLIT,$144
	// R0 points to the m.
	// We are already on m's g0 stack.

	// Save callee-save registers.
	SAVE_R19_TO_R28(8)
	SAVE_F8_TO_F15(88)

	MOVD    m_g0(R0), g
	BL	runtime·save_g(SB)

	BL	runtime·mstart(SB)

	// Restore callee-save registers.
	RESTORE_R19_TO_R28(8)
	RESTORE_F8_TO_F15(88)

	// Go is all done with this OS thread.
	// Tell pthread everything is ok (we never join with this thread, so
	// the value here doesn't really matter).
	MOVD	$0, R0

	RET

TEXT runtime·sigfwd(SB),NOSPLIT,$0-32
	MOVW	sig+8(FP), R0
	MOVD	info+16(FP), R1
	MOVD	ctx+24(FP), R2
	MOVD	fn+0(FP), R11
	BL	(R11)			// Alignment for ELF ABI?
	RET

TEXT runtime·sigtramp(SB),NOSPLIT|TOPFRAME,$192
	// Save callee-save registers in the case of signal forwarding.
	// Please refer to https://golang.org/issue/31827 .
	SAVE_R19_TO_R28(8*4)
	SAVE_F8_TO_F15(8*14)

	// If called from an external code context, g will not be set.
	// Save R0, since runtime·load_g will clobber it.
	MOVW	R0, 8(RSP)		// signum
	BL	runtime·load_g(SB)

#ifdef GOEXPERIMENT_regabiargs
	// Restore signum to R0.
	MOVW	8(RSP), R0
	// R1 and R2 already contain info and ctx, respectively.
#else
	MOVD	R1, 16(RSP)
	MOVD	R2, 24(RSP)
#endif
	BL	runtime·sigtrampgo<ABIInternal>(SB)

	// Restore callee-save registers.
	RESTORE_R19_TO_R28(8*4)
	RESTORE_F8_TO_F15(8*14)

	RET

//
// These trampolines help convert from Go calling convention to C calling convention.
// They should be called with asmcgocall.
// A pointer to the arguments is passed in R0.
// A single int32 result is returned in R0.
// (For more results, make an args/results structure.)
TEXT runtime·pthread_attr_init_trampoline(SB),NOSPLIT,$0
	MOVD	0(R0), R0		// arg 1 - attr
	CALL	libc_pthread_attr_init(SB)
	RET

TEXT runtime·pthread_attr_destroy_trampoline(SB),NOSPLIT,$0
	MOVD	0(R0), R0		// arg 1 - attr
	CALL	libc_pthread_attr_destroy(SB)
	RET

TEXT runtime·pthread_attr_setdetachstate_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - state
	MOVD	0(R0), R0		// arg 1 - attr
	CALL	libc_pthread_attr_setdetachstate(SB)
	RET

TEXT runtime·pthread_create_trampoline(SB),NOSPLIT,$0
	MOVD	0(R0), R1		// arg 2 - attr
	MOVD	8(R0), R2		// arg 3 - start
	MOVD	16(R0), R3		// arg 4 - arg
	SUB	$16, RSP
	MOVD	RSP, R0			// arg 1 - &threadid (discard)
	CALL	libc_pthread_create(SB)
	ADD	$16, RSP
	RET

TEXT runtime·pthread_kill_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - signo
	MOVD	0(R0), R0		// arg 1 - thread
	CALL	libc_pthread_kill(SB)
	RET

TEXT runtime·pthread_self_trampoline(SB),NOSPLIT,$0
	CALL	libc_pthread_self(SB)
	RET

TEXT runtime·pthread_attr_setstack_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - addr
	MOVD	16(R0), R2		// arg 3 - size
	MOVD	0(R0), R0		// arg 1 - attr
	CALL	libc_pthread_attr_setstack(SB)
	RET

TEXT runtime·sysconf_trampoline(SB),NOSPLIT,$0
	MOVW	0(R0), R0		// arg 1 - conf
	CALL	libc_sysconf(SB)
	RET

TEXT runtime·raise_trampoline(SB),NOSPLIT,$0
	MOVW	0(R0), R0		// arg 1 - sig
	CALL	libc_raise(SB)
	RET

TEXT runtime·malloc_trampoline(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args
	MOVD	0(R19), R0		// arg 1 - size
	CALL	libc_malloc(SB)
	MOVD	R0, 8(R19)
	RET

TEXT runtime·sem_init_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - pshared
	MOVW	16(R0), R2		// arg 3 - value
	MOVD	0(R0), R0		// arg 1 - sem
	CALL	libc_sem_init(SB)
	CMPW	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·sem_post_trampoline(SB),NOSPLIT,$0
	MOVD	0(R0), R0		// arg 1 - sem
	CALL	libc_sem_post(SB)
	CMPW	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·sem_reltimedwait_np_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - timeout
	MOVD	0(R0), R0		// arg 1 - sem
	CALL	libc_sem_reltimedwait_np(SB)
	CMPW	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·sem_wait_trampoline(SB),NOSPLIT,$0
	MOVD	0(R0), R0		// arg 1 - sem
	CALL	libc_sem_wait(SB)
	CMPW	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·exit_trampoline(SB),NOSPLIT,$0
	MOVW	0(R0), R0		// arg 1 - status
	CALL	libc_exit(SB)
	MOVD	$0, R0			// crash on failure
	MOVD	R0, (R0)
	RET

TEXT runtime·raiseproc_trampoline(SB),NOSPLIT,$0
	MOVD	R0, R19			// pointer to args
	CALL	libc_getpid(SB)		// arg 1 - pid
	MOVW	0(R19), R1		// arg 2 - signal
	CALL	libc_kill(SB)
	RET

TEXT runtime·sched_yield_trampoline(SB),NOSPLIT,$0
	CALL	libc_sched_yield(SB)
	RET

TEXT runtime·mmap_trampoline(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args
	MOVD	0(R19), R0		// arg 1 - addr
	MOVD	8(R19), R1		// arg 2 - len
	MOVW	16(R19), R2		// arg 3 - prot
	MOVW	20(R19), R3		// arg 4 - flags
	MOVW	24(R19), R4		// arg 5 - fid
	MOVW	28(R19), R5		// arg 6 - offset
	CALL	libc_mmap(SB)
	MOVD	$0, R1
	CMP	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R1		// errno
	MOVD	$0, R0
noerr:
	MOVD	R0, 32(R19)
	MOVD	R1, 40(R19)
	RET

TEXT runtime·munmap_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - len
	MOVD	0(R0), R0		// arg 1 - addr
	CALL	libc_munmap(SB)
	CMPW	$-1, R0
	BNE	3(PC)
	MOVD	$0, R0			// crash on failure
	MOVD	R0, (R0)
	RET

TEXT runtime·open_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - flags
	MOVW	12(R0), R2		// arg 3 - mode
	MOVD	0(R0), R0		// arg 1 - path
	MOVD	$0, R3			// varargs
	CALL	libc_open(SB)
	RET

TEXT runtime·close_trampoline(SB),NOSPLIT,$0
	MOVD	0(R0), R0		// arg 1 - fd
	CALL	libc_close(SB)
	RET

TEXT runtime·read_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - buf
	MOVW	16(R0), R2		// arg 3 - count
	MOVW	0(R0), R0		// arg 1 - fd
	CALL	libc_read(SB)
	CMP	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·write_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - buf
	MOVW	16(R0), R2		// arg 3 - count
	MOVW	0(R0), R0		// arg 1 - fd
	CALL	libc_write(SB)
	CMP	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·pipe2_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - flags
	MOVD	0(R0), R0		// arg 1 - filedes
	CALL	libc_pipe2(SB)
	CMPW	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·setitimer_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - new
	MOVD	16(R0), R2		// arg 3 - old
	MOVW	0(R0), R0		// arg 1 - which
	CALL	libc_setitimer(SB)
	RET

TEXT runtime·usleep_trampoline(SB),NOSPLIT,$0
	MOVD	0(R0), R0		// arg 1 - usec
	CALL	libc_usleep(SB)
	RET

TEXT runtime·clock_gettime_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - tp
	MOVW	0(R0), R0		// arg 1 - clock_id
	CALL	libc_clock_gettime(SB)
	CMPW	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·fcntl_trampoline(SB),NOSPLIT,$0
	MOVD	R0, R19
	MOVW	0(R19), R0		// arg 1 - fd
	MOVW	4(R19), R1		// arg 2 - cmd
	MOVW	8(R19), R2		// arg 3 - arg
	MOVD	$0, R3			// vararg
	CALL	libc_fcntl(SB)
	MOVD	$0, R1
	CMP	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R1
	MOVW	$-1, R0
noerr:
	MOVW	R0, 12(R19)
	MOVW	R1, 16(R19)
	RET

TEXT runtime·sigaction_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - new
	MOVD	16(R0), R2		// arg 3 - old
	MOVW	0(R0), R0		// arg 1 - sig
	CALL	libc_sigaction(SB)
	CMPW	$-1, R0
	BNE	3(PC)
	MOVD	$0, R0			// crash on syscall failure
	MOVD	R0, (R0)
	RET

TEXT runtime·sigprocmask_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - new
	MOVD	16(R0), R2		// arg 3 - old
	MOVW	0(R0), R0		// arg 1 - how
	CALL	libc_pthread_sigmask(SB)
	CMPW	$-1, R0
	BNE	3(PC)
	MOVD	$0, R0			// crash on syscall failure
	MOVD	R0, (R0)
	RET

TEXT runtime·sigaltstack_trampoline(SB),NOSPLIT,$0
	MOVD	8(R0), R1		// arg 2 - old
	MOVD	0(R0), R0		// arg 1 - new
	CALL	libc_sigaltstack(SB)
	CMPW	$-1, R0
	BNE	3(PC)
	MOVD	$0, R0			// crash on syscall failure
	MOVD	R0, (R0)
	RET

TEXT runtime·poll_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - npfds
	MOVW	16(R0), R2		// arg 3 - timeout
	MOVD	0(R0), R0		// arg 1 - pfds
	CALL	libc_poll(SB)
	CMPW	$-1, R0
	BNE	noerr
	CALL	libc_errno(SB)
	MOVW	(R0), R0		// errno
	NEG	R0, R0			// caller expects negative errno value
noerr:
	RET

TEXT runtime·libc_epoll_create1_trampoline(SB),NOSPLIT,$0-0
	JMP	libc_epoll_create1(SB)

TEXT runtime·libc_epoll_ctl_trampoline(SB),NOSPLIT,$0-0
	JMP	libc_epoll_ctl(SB)

TEXT runtime·libc_epoll_wait_trampoline(SB),NOSPLIT,$0-0
	JMP	libc_epoll_wait(SB)

TEXT runtime·API_VutexPend_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - desired
	MOVD	16(R0), R2		// arg 3 - timeout
	MOVD	0(R0), R0		// arg 1 - addr
	CALL	libc_API_VutexPend(SB)
	RET

TEXT runtime·API_VutexPostEx_trampoline(SB),NOSPLIT,$0
	MOVW	8(R0), R1		// arg 2 - value
	MOVW	12(R0), R2		// arg 3 - flags
	MOVD	0(R0), R0		// arg 1 - addr
	CALL	libc_API_VutexPostEx(SB)
	RET

TEXT runtime·vprocExitModeSet_trampoline(SB),NOSPLIT,$0
	MOVD	R0, R19			// pointer to args
	CALL	libc_getpid(SB)		// arg 1 - pid
	MOVW	0(R19), R1		// arg 2 - mode
	CALL	libc_vprocExitModeSet(SB)
	RET

// syscall calls a function in libc on behalf of the syscall package.
// syscall takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscall must be called on the g0 stack with the
// C calling convention (use libcCall).
//
// syscall expects a 32-bit result and tests for 32-bit -1
// to decide there was an error.
TEXT runtime·syscall(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	$0, R3			// vararg

	CALL	R11

	MOVD	R0, (4*8)(R19)		// r1
	MOVD	R1, (5*8)(R19)		// r2

	// Standard libc functions return -1 on error
	// and set errno.
	CMPW	$-1, R0
	BNE	ok

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (6*8)(R19)		// err

ok:
	RET

// syscallX calls a function in libc on behalf of the syscall package.
// syscallX takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscallX must be called on the g0 stack with the
// C calling convention (use libcCall).
//
// syscallX is like syscall but expects a 64-bit result
// and tests for 64-bit -1 to decide there was an error.
TEXT runtime·syscallX(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	$0, R3			// vararg

	CALL	R11

	MOVD	R0, (4*8)(R19)		// r1
	MOVD	R1, (5*8)(R19)		// r2

	// Standard libc functions return -1 on error
	// and set errno.
	CMP	$-1, R0
	BNE	ok

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (6*8)(R19)		// err

ok:
	RET

// syscallXerrno calls a function in libc on behalf of the syscall package.
// syscallXerrno takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscallXerrno must be called on the g0 stack with the
// C calling convention (use libcCall).
TEXT runtime·syscallXerrno(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	CALL	libc_errno(SB)
	MOVW	$0, (R0)

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	$0, R3			// vararg

	CALL	R11

	MOVD	R0, (4*8)(R19)		// r1
	MOVD	R1, (5*8)(R19)		// r2

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (6*8)(R19)		// err

	RET

// syscallXnull calls a function in libc on behalf of the syscall package.
// syscallXnull takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscallXnull must be called on the g0 stack with the
// C calling convention (use libcCall).
//
// syscallXnull is like syscall but expects a 64-bit result
// and tests for 64-bit 0(null) to decide there was an error.
TEXT runtime·syscallXnull(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	$0, R3			// vararg

	CALL	R11

	MOVD	R0, (4*8)(R19)		// r1
	MOVD	R1, (5*8)(R19)		// r2

	// Standard libc functions return 0(null) on error
	// and set errno.
	CMP	$0, R0
	BNE	ok

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (6*8)(R19)		// err

ok:
	RET

// syscall6 calls a function in libc on behalf of the syscall package.
// syscall6 takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	a4    uintptr
//	a5    uintptr
//	a6    uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscall6 must be called on the g0 stack with the
// C calling convention (use libcCall).
//
// syscall6 expects a 32-bit result and tests for 32-bit -1
// to decide there was an error.
TEXT runtime·syscall6(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	(4*8)(R19), R3		// a4
	MOVD	(5*8)(R19), R4		// a5
	MOVD	(6*8)(R19), R5		// a6
	MOVD	$0, R6			// vararg

	CALL	R11

	MOVD	R0, (7*8)(R19)		// r1
	MOVD	R1, (8*8)(R19)		// r2

	// Standard libc functions return -1 on error
	// and set errno.
	CMPW	$-1, R0
	BNE	ok

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (9*8)(R19)		// err

ok:
	RET

// syscall6X calls a function in libc on behalf of the syscall package.
// syscall6X takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	a4    uintptr
//	a5    uintptr
//	a6    uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscall6X must be called on the g0 stack with the
// C calling convention (use libcCall).
//
// syscall6X is like syscall6 but expects a 64-bit result
// and tests for 64-bit -1 to decide there was an error.
TEXT runtime·syscall6X(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	(4*8)(R19), R3		// a4
	MOVD	(5*8)(R19), R4		// a5
	MOVD	(6*8)(R19), R5		// a6
	MOVD	$0, R6			// vararg

	CALL	R11

	MOVD	R0, (7*8)(R19)		// r1
	MOVD	R1, (8*8)(R19)		// r2

	// Standard libc functions return -1 on error
	// and set errno.
	CMP	$-1, R0
	BNE	ok

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (9*8)(R19)		// err

ok:
	RET

// syscall6Xerrno calls a function in libc on behalf of the syscall package.
// syscall6Xerrno takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	a4    uintptr
//	a5    uintptr
//	a6    uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscall6Xerrno must be called on the g0 stack with the
// C calling convention (use libcCall).
TEXT runtime·syscall6Xerrno(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	CALL	libc_errno(SB)
	MOVW	$0, (R0)

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	(4*8)(R19), R3		// a4
	MOVD	(5*8)(R19), R4		// a5
	MOVD	(6*8)(R19), R5		// a6
	MOVD	$0, R6			// vararg

	CALL	R11

	MOVD	R0, (7*8)(R19)		// r1
	MOVD	R1, (8*8)(R19)		// r2

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (9*8)(R19)		// err

	RET

// syscall10 calls a function in libc on behalf of the syscall package.
// syscall10 takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	a4    uintptr
//	a5    uintptr
//	a6    uintptr
//	a7    uintptr
//	a8    uintptr
//	a9    uintptr
//	a10   uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscall10 must be called on the g0 stack with the
// C calling convention (use libcCall).
TEXT runtime·syscall10(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	(4*8)(R19), R3		// a4
	MOVD	(5*8)(R19), R4		// a5
	MOVD	(6*8)(R19), R5		// a6
	MOVD	(7*8)(R19), R6		// a7
	MOVD	(8*8)(R19), R7		// a8
	MOVD	(9*8)(R19), R8		// a9
	MOVD	(10*8)(R19), R9		// a10
	MOVD	$0, R10			// vararg

	CALL	R11

	MOVD	R0, (11*8)(R19)		// r1
	MOVD	R1, (12*8)(R19)		// r2

	// Standard libc functions return -1 on error
	// and set errno.
	CMPW	$-1, R0
	BNE	ok

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (13*8)(R19)		// err

ok:
	RET

// syscall10X calls a function in libc on behalf of the syscall package.
// syscall10X takes a pointer to a struct like:
// struct {
//	fn    uintptr
//	a1    uintptr
//	a2    uintptr
//	a3    uintptr
//	a4    uintptr
//	a5    uintptr
//	a6    uintptr
//	a7    uintptr
//	a8    uintptr
//	a9    uintptr
//	a10   uintptr
//	r1    uintptr
//	r2    uintptr
//	err   uintptr
// }
// syscall10X must be called on the g0 stack with the
// C calling convention (use libcCall).
//
// syscall10X is like syscall10 but expects a 64-bit result
// and tests for 64-bit -1 to decide there was an error.
TEXT runtime·syscall10X(SB),NOSPLIT,$0
	MOVD    R0, R19			// pointer to args

	MOVD	(0*8)(R19), R11		// fn
	MOVD	(1*8)(R19), R0		// a1
	MOVD	(2*8)(R19), R1		// a2
	MOVD	(3*8)(R19), R2		// a3
	MOVD	(4*8)(R19), R3		// a4
	MOVD	(5*8)(R19), R4		// a5
	MOVD	(6*8)(R19), R5		// a6
	MOVD	(7*8)(R19), R6		// a7
	MOVD	(8*8)(R19), R7		// a8
	MOVD	(9*8)(R19), R8		// a9
	MOVD	(10*8)(R19), R9		// a10
	MOVD	$0, R10			// vararg

	CALL	R11

	MOVD	R0, (11*8)(R19)		// r1
	MOVD	R1, (12*8)(R19)		// r2

	// Standard libc functions return -1 on error
	// and set errno.
	CMP	$-1, R0
	BNE	ok

	// Get error code from libc.
	CALL	libc_errno(SB)
	MOVW	(R0), R0
	MOVD	R0, (13*8)(R19)		// err

ok:
	RET

TEXT runtime·issetugid_trampoline(SB),NOSPLIT,$0
	MOVD	R0, R19			// pointer to args
	CALL	libc_issetugid(SB)
	MOVW	R0, 0(R19)		// return value
	RET
