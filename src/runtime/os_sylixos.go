// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package runtime

import (
	"internal/abi"
	"internal/goarch"
	"unsafe"
)

type mOS struct {
	waitsema uintptr // semaphore for parking on locks
}

var urandom_dev = []byte("/dev/urandom\x00")

var startupRandomData []byte

//go:nosplit
func readRandom(r []byte) int {
	if startupRandomData != nil {
		n := copy(r, startupRandomData)
		return int(n)
	}
	fd := open(&urandom_dev[0], 0 /* O_RDONLY */, 0)
	n := read(fd, unsafe.Pointer(&r[0]), int32(len(r)))
	closefd(fd)
	return int(n)
}
func setProcessCPUProfiler(hz int32) {
	setProcessCPUProfilerTimer(hz)
}

func setThreadCPUProfiler(hz int32) {
	setThreadCPUProfilerHz(hz)
}

//go:nosplit
func validSIGPROF(mp *m, c *sigctxt) bool {
	return true
}

//go:nosplit
func semacreate(mp *m) {
	if mp.waitsema != 0 {
		return
	}

	var sem *semt

	// Call libc's malloc rather than malloc. This will
	// allocate space on the C heap. We can't call mallocgc
	// here because it could cause a deadlock.
	sem = (*semt)(malloc(unsafe.Sizeof(*sem)))
	if sem_init(sem, 0, 0) != 0 {
		throw("sem_init")
	}
	mp.waitsema = uintptr(unsafe.Pointer(sem))
}

//go:nosplit
func semasleep(ns int64) int32 {
	mp := getg().m
	if ns >= 0 {
		var ts timespec

		ts.tv_sec = ns / 1e9
		ts.tv_nsec = ns % 1e9

		if err := sem_reltimedwait_np((*semt)(unsafe.Pointer(mp.waitsema)), &ts); err != 0 {
			if err == -_ETIMEDOUT || err == -_EAGAIN || err == -_EINTR {
				return -1
			}
			println("sem_reltimedwait_np err ", err, " ts.tv_sec ", ts.tv_sec, " ts.tv_nsec ", ts.tv_nsec, " ns ", ns, " id ", mp.id)
			throw("sem_reltimedwait_np")
		}
		return 0
	}
	for {
		err := sem_wait((*semt)(unsafe.Pointer(mp.waitsema)))
		if err == 0 {
			break
		}
		if err == -_EINTR {
			continue
		}
		throw("sem_wait")
	}
	return 0
}

//go:nosplit
func semawakeup(mp *m) {
	if sem_post((*semt)(unsafe.Pointer(mp.waitsema))) != 0 {
		throw("sem_post")
	}
}

//go:nosplit
func futexsleep(addr *uint32, val uint32, ns int64) {
	var timeout int32
	if ns >= 0 {
		// The timeout is specified in microseconds - ensure that we
		// do not end up dividing to zero, which would put us to sleep
		// indefinitely...
		timeout = timediv(ns, 1000000, nil)
		if timeout == 0 {
			timeout = 1
		}
	}
	API_VutexPend(addr, val, uintptr(timeout))
}

//go:nosplit
func futexwakeup(addr *uint32, cnt uint32) {
	API_VutexPostEx(addr, *addr, _LW_OPTION_VUTEX_FLAG_DONTSET|_LW_OPTION_VUTEX_FLAG_DEEPWAKE)
}

func osinit() {
	vprocExitModeSet(_LW_VPROC_EXIT_FORCE)

	ncpu = sysconf(__SC_NPROCESSORS_ONLN)
	physPageSize = uintptr(sysconf(__SC_PAGE_SIZE))
}

// mstart_stub provides glue code to call mstart from pthread_create.
func mstart_stub()

// May run with m.p==nil, so write barriers are not allowed.
//
//go:nowritebarrier
func newosproc(mp *m) {
	var (
		attr pthreadattr
		oset sigset
	)

	if pthread_attr_init(&attr) != 0 {
		throw("pthread_attr_init")
	}

	if pthread_attr_setstack(&attr, unsafe.Pointer(mp.g0.stack.lo), mp.g0.stack.hi-mp.g0.stack.lo) != 0 {
		throw("pthread_attr_setstack")
	}

	if pthread_attr_setdetachstate(&attr, _PTHREAD_CREATE_DETACHED) != 0 {
		throw("pthread_attr_setdetachstate")
	}

	// Disable signals during create, so that the new thread starts
	// with signals disabled. It will enable them in minit.
	sigprocmask(_SIG_SETMASK, &sigset_all, &oset)
	ret := retryOnEAGAIN(func() int32 {
		return pthread_create(&attr, abi.FuncPCABI0(mstart_stub), unsafe.Pointer(mp))
	})
	sigprocmask(_SIG_SETMASK, &oset, nil)
	if ret != 0 {
		print("runtime: failed to create new OS thread (have ", mcount(), " already; errno=", ret, ")\n")
		if ret == _EAGAIN {
			println("runtime: may need to increase max user processes (ulimit -u)")
		}
		throw("newosproc")
	}

	pthread_attr_destroy(&attr)
}

//go:nosplit
//go:nowritebarrierrec
func setsigstack(i uint32) {
	var sa sigactiont
	sigaction(i, nil, &sa)
	if sa.sa_flags&_SA_ONSTACK != 0 {
		return
	}
	sa.sa_flags |= _SA_ONSTACK
	sigaction(i, &sa, nil)
}

// setSignalstackSP sets the ss_sp field of a stackt.
//
//go:nosplit
func setSignalstackSP(s *stackt, sp uintptr) {
	*(*uintptr)(unsafe.Pointer(&s.ss_sp)) = sp
}

// sigPerThreadSyscall is only used on SylixOS, so we assign a bogus signal
// number.
const sigPerThreadSyscall = 1 << 31

//go:nosplit
func runPerThreadSyscall() {
	throw("runPerThreadSyscall only valid on SylixOS")
}

//go:nosplit
//go:nowritebarrierrec
func getsig(i uint32) uintptr {
	var sa sigactiont
	sigaction(i, nil, &sa)
	return sa.sa_handler
}

// It's hard to tease out exactly how big a Sigset is, but
// rt_sigprocmask crashes if we get it wrong, so if binaries
// are running, this is right.
type sigset [2]uint32

var sigset_all = sigset{^uint32(0) & ^uint32(1<<(_SIGTERM-1)), ^uint32(0)}

//go:nosplit
//go:nowritebarrierrec
func sigaddset(mask *sigset, i int) {
	(*mask)[(i-1)/32] |= 1 << ((uint32(i) - 1) & 31)
}

func sigdelset(mask *sigset, i int) {
	(*mask)[(i-1)/32] &^= 1 << ((uint32(i) - 1) & 31)
}

//go:nosplit
func sigfillset(mask *uint64) {
	*mask = ^uint64(0) & ^uint64(1<<(_SIGTERM-1))
}

//go:nosplit
func (c *sigctxt) fixsigcode(sig uint32) {
	if sig == _SIGTERM {
		exit(int32(c.sigaddr()))
	}
}

func sigtramp()

//go:nosplit
//go:nowritebarrierrec
func setsig(i uint32, fn uintptr) {
	var sa sigactiont

	sa.sa_flags = _SA_SIGINFO | _SA_ONSTACK | _SA_RESTART
	sa.sa_mask = ^uint64(0)
	if fn == abi.FuncPCABIInternal(sighandler) { // abi.FuncPCABIInternal(sighandler) matches the callers in signal_unix.go
		fn = abi.FuncPCABI0(sigtramp)
	}
	sa.sa_handler = fn
	sigaction(i, &sa, nil)
}

func signalM(mp *m, sig int) {
	pthread_kill(pthread(mp.procid), sig)
}

var (
	env **byte
)

//go:nosplit
func env_index(argv **byte, i int32) *byte {
	return *(**byte)(add(unsafe.Pointer(argv), uintptr(i)*goarch.PtrSize))
}

func goenvs() {
	n := int32(0)
	for env_index(env, n) != nil {
		n++
	}

	envs = make([]string, n)
	for i := int32(0); i < n; i++ {
		envs[i] = gostring(env_index(env, i))
	}
}

func sylixosenvs(e **byte) {
	env = e
}

// Called to initialize a new m (including the bootstrap m).
// Called on the parent thread (main thread in case of bootstrap), can allocate memory.
func mpreinit(mp *m) {
	mp.gsignal = malg(32 * 1024) // SylixOS wants >= 2K
	mp.gsignal.m = mp
}

// Called to initialize a new m (including the bootstrap m).
// Called on the new thread, cannot allocate memory.
func minit() {
	minitSignals()

	getg().m.procid = uint64(pthread_self())
}

// Called from dropm to undo the effect of an minit.
//
//go:nosplit
func unminit() {
	unminitSignals()
}

// Called from exitm, but not from drop, to undo the effect of thread-owned
// resources in minit, semacreate, or elsewhere. Do not take locks after calling this.
func mdestroy(mp *m) {
}
