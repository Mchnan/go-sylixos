// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

/*
Input to cgo -cdefs

GOARCH=amd64 go tool cgo -cdefs defs_linux.go defs1_linux.go >defs_linux_amd64.h
*/

package runtime

/*
// Linux glibc and Linux kernel define different and conflicting
// definitions for struct sigaction, struct timespec, etc.
// We want the kernel ones, which are in the asm/* headers.
// But then we'd get conflicts when we include the system
// headers for things like ucontext_t, so that happens in
// a separate file, defs1.go.

#define SYLIXOS
#include <SylixOS.h>
#include <sys/socket.h>
#include <sys/resource.h>
#include <sys/un.h>
#include <net/if_dl.h>

#include <netinet6/in6.h>
#include <netinet6/icmp6.h>
#include <net/if.h>
#include <net/route.h>
#include <termios.h>
#include <poll.h>
#include <sys/epoll.h>
#include <sys/time.h>
#include <sys/mman.h>
#include <pthread.h>
#include <semaphore.h>
#include <sched.h>

*/
import "C"

const (
	_EPERM     = C.EPERM
	_ENOENT    = C.ENOENT
	_EINTR     = C.EINTR
	_EAGAIN    = C.EAGAIN
	_ENOMEM    = C.ENOMEM
	_EACCES    = C.EACCES
	_EFAULT    = C.EFAULT
	_EINVAL    = C.EINVAL
	_ETIMEDOUT = C.ETIMEDOUT

	_PROT_NONE  = C.PROT_NONE
	_PROT_READ  = C.PROT_READ
	_PROT_WRITE = C.PROT_WRITE
	_PROT_EXEC  = C.PROT_EXEC

	_MAP_ANON      = C.MAP_ANONYMOUS
	_MAP_PRIVATE   = C.MAP_PRIVATE
	_MAP_FIXED     = C.MAP_FIXED
	_MADV_DONTNEED = C.MADV_DONTNEED

	_SIGHUP    = C.SIGHUP
	_SIGINT    = C.SIGINT
	_SIGQUIT   = C.SIGQUIT
	_SIGILL    = C.SIGILL
	_SIGTRAP   = C.SIGTRAP
	_SIGABRT   = C.SIGABRT
	_SIGBUS    = C.SIGBUS
	_SIGFPE    = C.SIGFPE
	_SIGKILL   = C.SIGKILL
	_SIGUSR1   = C.SIGUSR1
	_SIGSEGV   = C.SIGSEGV
	_SIGUSR2   = C.SIGUSR2
	_SIGPIPE   = C.SIGPIPE
	_SIGALRM   = C.SIGALRM
	_SIGCHLD   = C.SIGCHLD
	_SIGCONT   = C.SIGCONT
	_SIGSTOP   = C.SIGSTOP
	_SIGTSTP   = C.SIGTSTP
	_SIGTTIN   = C.SIGTTIN
	_SIGTTOU   = C.SIGTTOU
	_SIGURG    = C.SIGURG
	_SIGXCPU   = C.SIGXCPU
	_SIGXFSZ   = C.SIGXFSZ
	_SIGVTALRM = C.SIGVTALRM
	_SIGPROF   = C.SIGPROF
	_SIGWINCH  = C.SIGWINCH
	_SIGIO     = C.SIGIO
	_SIGPWR    = C.SIGPWR
	_SIGSYS    = C.SIGSYS
	_SIGTERM   = C.SIGTERM

	_NSIG = 64

	_FPE_INTDIV = C.FPE_INTDIV
	_FPE_INTOVF = C.FPE_INTOVF
	_FPE_FLTDIV = C.FPE_FLTDIV
	_FPE_FLTOVF = C.FPE_FLTOVF
	_FPE_FLTUND = C.FPE_FLTUND
	_FPE_FLTRES = C.FPE_FLTRES
	_FPE_FLTINV = C.FPE_FLTINV
	_FPE_FLTSUB = C.FPE_FLTSUB

	_BUS_ADRALN = C.BUS_ADRALN
	_BUS_ADRERR = C.BUS_ADRERR
	_BUS_OBJERR = C.BUS_OBJERR

	_SEGV_MAPERR = C.SEGV_MAPERR
	_SEGV_ACCERR = C.SEGV_ACCERR

	_ITIMER_REAL    = C.ITIMER_REAL
	_ITIMER_VIRTUAL = C.ITIMER_VIRTUAL
	_ITIMER_PROF    = C.ITIMER_PROF

	_O_RDONLY   = C.O_RDONLY
	_O_WRONLY   = C.O_WRONLY
	_O_NONBLOCK = C.O_NONBLOCK
	_O_CREAT    = C.O_CREAT
	_O_TRUNC    = C.O_TRUNC
	_O_CLOEXEC  = C.O_CLOEXEC

	_SS_DISABLE  = C.SS_DISABLE
	_SI_USER     = C.SI_USER
	_SIG_BLOCK   = C.SIG_BLOCK
	_SIG_UNBLOCK = C.SIG_UNBLOCK
	_SIG_SETMASK = C.SIG_SETMASK

	_SA_SIGINFO = C.SA_SIGINFO
	_SA_RESTART = C.SA_RESTART
	_SA_ONSTACK = C.SA_ONSTACK

	_PTHREAD_CREATE_DETACHED = C.PTHREAD_CREATE_DETACHED

	__SC_PAGE_SIZE        = C._SC_PAGE_SIZE
	__SC_NPROCESSORS_ONLN = C._SC_NPROCESSORS_ONLN

	_F_SETFD    = C.F_SETFD
	_F_SETFL    = C.F_SETFL
	_F_GETFD    = C.F_GETFD
	_F_GETFL    = C.F_GETFL
	_FD_CLOEXEC = C.FD_CLOEXEC

	_CLOCK_REALTIME           = C.CLOCK_REALTIME
	_CLOCK_MONOTONIC          = C.CLOCK_MONOTONIC
	_CLOCK_PROCESS_CPUTIME_ID = C.CLOCK_PROCESS_CPUTIME_ID
	_CLOCK_THREAD_CPUTIME_ID  = C.CLOCK_THREAD_CPUTIME_ID

	_AF_UNIX    = C.AF_UNIX
	_SOCK_DGRAM = C.SOCK_DGRAM

	EPOLLIN      = C.EPOLLIN
	EPOLLPRI     = C.EPOLLPRI
	EPOLLOUT     = C.EPOLLOUT
	EPOLLERR     = C.EPOLLERR
	EPOLLHUP     = C.EPOLLHUP
	EPOLLRDHUP   = C.EPOLLHUP
	EPOLLONESHOT = C.EPOLLONESHOT
	EPOLLET      = 0x80000000

	EPOLL_CLOEXEC  = C.EPOLL_CLOEXEC
	EPOLL_NONBLOCK = C.EPOLL_NONBLOCK

	EPOLL_CTL_ADD = C.EPOLL_CTL_ADD
	EPOLL_CTL_DEL = C.EPOLL_CTL_DEL
	EPOLL_CTL_MOD = C.EPOLL_CTL_MOD
)

// type timespec C.struct_timespec
type timespec struct {
	tv_sec  int64
	tv_nsec int64 // ACOINFO TODO long type
}

//go:nosplit
func (ts *timespec) setNsec(ns int64) {
	ts.tv_sec = ns / 1e9
	ts.tv_nsec = ns % 1e9
}

// type timeval C.struct_timeval
type timeval struct {
	tv_sec  int64
	tv_usec int64 // ACOINFO TODO long type
}

func (tv *timeval) set_usec(x int32) {
	tv.tv_usec = int64(x)
}

// type itimerspec C.struct_itimerspec
type itimerspec struct {
	it_interval timespec
	it_value    timespec
}

// type itimerval C.struct_itimerval
type itimerval struct {
	it_interval timeval
	it_value    timeval
}

// type sigactiont C.struct_sigaction
type sigactiont struct {
	sa_handler  uintptr
	sa_mask     uint64
	sa_flags    uint32
	sa_restorer uintptr
}

// type siginfo C.siginfo_t
type siginfoFields struct {
	si_signo  int32
	si_errno  int32
	si_code   int32
	si_pid    int32
	si_uid    int32
	si_status int32
	si_utime  uintptr
	si_stime  uintptr

	// below here is a union; si_addr is the only field we use
	si_addr uint64
}

type siginfo struct {
	siginfoFields

	// Pad struct to the max size in the kernel.
	_ [4]uintptr
}

// type sigevent C.struct_sigevent
type sigeventFields struct {
	signo                   int32
	value                   uintptr
	notify                  int32
	sigev_notify_function   uintptr
	sigev_notify_attributes uintptr
	// below here is a union; sigev_notify_thread_id is the only field we use
	sigev_notify_thread_id uint64
}

type sigevent struct {
	sigeventFields

	// Pad struct to the max size in the kernel.
	_ [8]uint64
}

// type stackt C.stack_t
type stackt struct {
	ss_sp    *byte
	ss_size  uintptr
	ss_flags int32
}

// type semt C.sem_t
type semt struct {
	pxsem  uintptr
	resraw uintptr
	pad    [5]uintptr
}

// type sockaddr_un C.struct_sockaddr_un
type sockaddr_un struct {
	len    uint8
	family uint8
	path   [104]byte
}

// type sched_param C.struct_sched_param
type sched_param struct {
	sched_priority        int32
	sched_ss_low_priority int32
	sched_ss_repl_period  timespec
	sched_ss_init_budget  timespec
	sched_ss_max_repl     int32
	sched_pad             [12]uintptr
}

// type pthread C.pthread_t
type pthread uintptr

// type pthreadattr C.pthread_attr_t
type pthreadattr struct {
	name             *byte
	stack            uintptr
	stack_guard_size uint64
	stack_size       uint64
	sched_policy     int32
	option           uint64
	schedparam       sched_param
	pad              [8]uintptr
}

type EpollEvent C.struct_epoll_event
