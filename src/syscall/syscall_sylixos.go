// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

// SylixOS system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and wrap
// it in our own nicer implementation, either here or in
// syscall_sylixos.go , syscall_sylixos1.go , syscall_sylixos2.go
// or syscall_unix.go.

package syscall

import (
	"internal/abi"
	"sync"
	"unsafe"
)

// SylixOS is different from Linux and BSD.
// It does not have native Indirect system call.
// We try to use the libc function to imitate it.
// TODO: Add more syscall for SylixOS
func Syscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	if trap == SYS_EXECVE {
		r1, r2, err = syscall(abi.FuncPCABI0(libc_execve_trampoline), a1, a2, a3)
		return
	} else if trap == SYS_IOCTL {
		r1, r2, err = syscall(abi.FuncPCABI0(libc_ioctl_trampoline), a1, a2, a3)
		return
	} else {
		// Do not panic on unknown syscall
		return 0, 0, ENOSYS
	}
}

// TODO: Add more syscall for SylixOS
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	if trap == SYS_EXECVE {
		r1, r2, err = syscall(abi.FuncPCABI0(libc_execve_trampoline), a1, a2, a3)
		return
	} else if trap == SYS_IOCTL {
		r1, r2, err = syscall(abi.FuncPCABI0(libc_ioctl_trampoline), a1, a2, a3)
		return
	} else {
		// Do not panic on unknown syscall
		return 0, 0, ENOSYS
	}
}

// TODO: Add more syscall for SylixOS
func RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno) {
	if trap == SYS_EXECVE {
		r1, r2, err = rawSyscall(abi.FuncPCABI0(libc_execve_trampoline), a1, a2, a3)
		return
	} else if trap == SYS_IOCTL {
		r1, r2, err = rawSyscall(abi.FuncPCABI0(libc_ioctl_trampoline), a1, a2, a3)
		return
	} else {
		// Do not panic on unknown syscall
		return 0, 0, ENOSYS
	}
}

// TODO: Add more syscall for SylixOS
func RawSyscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno) {
	if trap == SYS_EXECVE {
		r1, r2, err = rawSyscall(abi.FuncPCABI0(libc_execve_trampoline), a1, a2, a3)
		return
	} else if trap == SYS_IOCTL {
		r1, r2, err = rawSyscall(abi.FuncPCABI0(libc_ioctl_trampoline), a1, a2, a3)
		return
	} else {
		// Do not panic on unknown syscall
		return 0, 0, ENOSYS
	}
}

const ImplementsGetwd = true

func Getwd() (string, error) {
	var buf [pathMax]byte
	_, err := getcwd(buf[:])
	if err != nil {
		return "", err
	}
	n := clen(buf[:])
	if n < 1 {
		return "", EINVAL
	}
	return string(buf[:n]), nil
}

/*
 * Wrapped
 */

//sysnb	getgroups(ngid int, gid *_Gid_t) (n int, err error)
//sysnb	setgroups(ngid int, gid *_Gid_t) (err error)

func Getgroups() (gids []int, err error) {
	n, err := getgroups(0, nil)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	// Sanity check group count. Max is 16 on BSD.
	if n < 0 || n > 1000 {
		return nil, EINVAL
	}

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != nil {
		return nil, err
	}
	gids = make([]int, n)
	for i, v := range a[0:n] {
		gids[i] = int(v)
	}
	return
}

func Setgroups(gids []int) (err error) {
	if len(gids) == 0 {
		return setgroups(0, nil)
	}

	a := make([]_Gid_t, len(gids))
	for i, v := range gids {
		a[i] = _Gid_t(v)
	}
	return setgroups(len(a), &a[0])
}

var fdDirMap = map[int]uintptr{}
var fdDirMapMu sync.Mutex

func Close(fd int) (err error) {
	fdDirMapMu.Lock()
	dir, ok := fdDirMap[fd]
	if ok {
		delete(fdDirMap, fd)
	}
	fdDirMapMu.Unlock()
	if ok {
		return closedir(dir)
	} else {
		return rawClose(fd)
	}
}

func fdToDir(fd int) (uintptr, error) {
	fdDirMapMu.Lock()
	dir, ok := fdDirMap[fd]
	if ok {
		fdDirMapMu.Unlock()
		return dir, nil
	} else {
		dir, err := fdopendir(fd)
		if err != nil {
			fdDirMapMu.Unlock()
			return 0, err
		}
		fdDirMap[fd] = dir
		fdDirMapMu.Unlock()
		return dir, nil
	}
}

func ReadDirent(fd int, buf []byte) (n int, err error) {
	dir, e := fdToDir(fd)
	if e != nil {
		return -1, e
	}

	var entry Dirent
	var entryp *Dirent
	readdir_r(dir, &entry, &entryp)
	if entryp == nil {
		return 0, nil
	}

	n = 0
	max := len(buf)
	for i, c := range entryp.Name {
		if n >= max {
			break
		}
		n++
		buf[i] = c
		if c == 0 {
			break
		}
	}

	return n, nil
}

// ParseDirent parses up to max directory entries in buf,
// appending the names to names. It returns the number of
// bytes consumed from buf, the number of entries added
// to names, and the new names slice.
func ParseDirent(buf []byte, max int, names []string) (consumed int, count int, newnames []string) {
	consumed = 0
	for _, c := range buf {
		if max >= 0 {
			if consumed >= max {
				break
			}
		}
		consumed++
		if c == 0 {
			break
		}
	}

	if consumed > 0 {
		name := buf[0 : consumed-1]
		names = append(names, string(name))
		return consumed, 1, names
	} else {
		return 0, 0, names
	}
}

// FcntlFlock performs a fcntl syscall for the F_GETLK, F_SETLK or F_SETLKW command.
func FcntlFlock(fd uintptr, cmd int, lk *Flock_t) error {
	_, e1 := fcntlPtr(int(fd), cmd, unsafe.Pointer(lk))
	return e1
}

// Wait status
// 31     24       16        8        0
// +-------+--------+--------+--------+
// |       | stop   |  term  |  exit  |
// |       | signal | signal | status |
// +-------+--------+--------+--------+

type WaitStatus uint32

const (
	termSigMask  = 0xFF00
	termSigShift = 8

	stopSigMask  = 0xFF0000
	stopSigShift = 16

	sigMask   = termSigMask | stopSigMask
	noStopSig = 0

	_WUNTRACED = 2 /* Tell about stopped	*/

	CLD_EXITED    = 1 << stopSigShift
	CLD_KILLED    = 2 << stopSigShift
	CLD_DUMPED    = 3 << stopSigShift
	CLD_TRAPPED   = 4 << stopSigShift
	CLD_STOPPED   = 5 << stopSigShift
	CLD_CONTINUED = 6 << stopSigShift
)

func (w WaitStatus) Exited() bool {
	var stopSig = w & stopSigMask
	switch stopSig {
	case noStopSig, CLD_EXITED, CLD_KILLED, CLD_DUMPED:
		return true
	default:
		return false
	}
}

func (w WaitStatus) ExitStatus() int {
	if !w.Exited() {
		return -1
	}
	return int(w & 0xFF)
}

func (w WaitStatus) Signaled() bool {
	return w&termSigMask != noStopSig
}

func (w WaitStatus) Signal() Signal {
	if !w.Signaled() {
		return -1
	}
	return Signal(w>>termSigShift) & 0xFF
}

func (w WaitStatus) CoreDump() bool { return w&stopSigMask == CLD_DUMPED }

func (w WaitStatus) Stopped() bool { return w&stopSigMask == CLD_STOPPED }

func (w WaitStatus) Continued() bool { return w&stopSigMask == CLD_CONTINUED }

func (w WaitStatus) StopSignal() Signal {
	if w.Stopped() {
		return SIGSTOP
	} else if w.Continued() {
		return SIGCONT
	} else {
		return -1
	}
}

func (w WaitStatus) TrapCause() int { return -1 }

//sys	wait4(pid int, wstatus *_C_int, options int, rusage *Rusage) (wpid int, err error)

func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) {
	var status _C_int
	wpid, err = wait4(pid, &status, options|_WUNTRACED, rusage)
	if wstatus != nil {
		*wstatus = WaitStatus(status)
	}
	return
}

//sys	accept(s int, rsa *RawSockaddrAny, addrlen *_Socklen) (fd int, err error)
//sys	bind(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sys	connect(s int, addr unsafe.Pointer, addrlen _Socklen) (err error)
//sysnb	socket(domain int, typ int, proto int) (fd int, err error)
//sys	getsockopt(s int, level int, name int, val unsafe.Pointer, vallen *_Socklen) (err error)
//sys	setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error)
//sysnb	getpeername(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sysnb	getsockname(fd int, rsa *RawSockaddrAny, addrlen *_Socklen) (err error)
//sys	Shutdown(s int, how int) (err error)

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Len = SizeofSockaddrInet4
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, EINVAL
	}
	sa.raw.Len = SizeofSockaddrInet6
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) {
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) || n == 0 {
		return nil, 0, EINVAL
	}
	sa.raw.Len = byte(3 + n) // 2 for Family, Len; 1 for NUL
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = uint8(name[i])
	}
	return unsafe.Pointer(&sa.raw), _Socklen(sa.raw.Len), nil
}

func (sa *SockaddrDatalink) sockaddr() (unsafe.Pointer, _Socklen, error) {
	if sa.Index == 0 {
		return nil, 0, EINVAL
	}
	sa.raw.Len = sa.Len
	sa.raw.Family = AF_LINK
	sa.raw.Index = sa.Index
	sa.raw.Type = sa.Type
	sa.raw.Nlen = sa.Nlen
	sa.raw.Alen = sa.Alen
	sa.raw.Slen = sa.Slen
	sa.raw.Data = sa.Data
	return unsafe.Pointer(&sa.raw), SizeofSockaddrDatalink, nil
}

func anyToSockaddr(rsa *RawSockaddrAny) (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_LINK:
		pp := (*RawSockaddrDatalink)(unsafe.Pointer(rsa))
		sa := new(SockaddrDatalink)
		sa.Len = pp.Len
		sa.Family = pp.Family
		sa.Index = pp.Index
		sa.Type = pp.Type
		sa.Nlen = pp.Nlen
		sa.Alen = pp.Alen
		sa.Slen = pp.Slen
		sa.Data = pp.Data
		return sa, nil

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		if pp.Len < 2 || pp.Len > SizeofSockaddrUnix {
			return nil, EINVAL
		}
		sa := new(SockaddrUnix)

		// Some BSDs include the trailing NUL in the length, whereas
		// others do not. Work around this by subtracting the leading
		// family and len. The path is then scanned to see if a NUL
		// terminator still exists within the length.
		n := int(pp.Len) - 2 // subtract leading Family, Len
		for i := 0; i < n; i++ {
			if pp.Path[i] == 0 {
				// found early NUL; assume Len included the NUL
				// or was overestimating.
				n = i
				break
			}
		}
		bytes := (*[len(pp.Path)]byte)(unsafe.Pointer(&pp.Path[0]))[0:n]
		sa.Name = string(bytes)
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.Addr = pp.Addr
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		sa.Addr = pp.Addr
		return sa, nil
	}
	return nil, EAFNOSUPPORT
}

func Accept(fd int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if err != nil {
		return
	}
	sa, err = anyToSockaddr(&rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

func Getsockname(fd int) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil {
		return
	}
	return anyToSockaddr(&rsa)
}

//sysnb socketpair(domain int, typ int, proto int, fd *[2]int32) (err error)

func GetsockoptByte(fd, level, opt int) (value byte, err error) {
	var n byte
	vallen := _Socklen(1)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&n), &vallen)
	return n, err
}

func GetsockoptInet4Addr(fd, level, opt int) (value [4]byte, err error) {
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	return value, err
}

func GetsockoptIPMreq(fd, level, opt int) (*IPMreq, error) {
	var value IPMreq
	vallen := _Socklen(SizeofIPMreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptIPv6Mreq(fd, level, opt int) (*IPv6Mreq, error) {
	var value IPv6Mreq
	vallen := _Socklen(SizeofIPv6Mreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptIPv6MTUInfo(fd, level, opt int) (*IPv6MTUInfo, error) {
	var value IPv6MTUInfo
	vallen := _Socklen(SizeofIPv6MTUInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

func GetsockoptICMPv6Filter(fd, level, opt int) (*ICMPv6Filter, error) {
	var value ICMPv6Filter
	vallen := _Socklen(SizeofICMPv6Filter)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
}

//sys   recvfrom(fd int, p []byte, flags int, from *RawSockaddrAny, fromlen *_Socklen) (n int, err error)
//sys   sendto(s int, buf []byte, flags int, to unsafe.Pointer, addrlen _Socklen) (err error)
//sys	recvmsg(s int, msg *Msghdr, flags int) (n int, err error)

func recvmsgRaw(fd int, p, oob []byte, flags int, rsa *RawSockaddrAny) (n, oobn int, recvflags int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(unsafe.Pointer(rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = &p[0]
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		// receive at least one normal byte
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); err != nil {
		return
	}
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	return
}

//sys	sendmsg(s int, msg *Msghdr, flags int) (n int, err error)

func sendmsgN(fd int, p, oob []byte, ptr unsafe.Pointer, salen _Socklen, flags int) (n int, err error) {
	var msg Msghdr
	msg.Name = (*byte)(ptr)
	msg.Namelen = uint32(salen)
	var iov Iovec
	if len(p) > 0 {
		iov.Base = &p[0]
		iov.SetLen(len(p))
	}
	var dummy byte
	if len(oob) > 0 {
		// send at least one normal byte
		if len(p) == 0 {
			iov.Base = &dummy
			iov.SetLen(1)
		}
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	}
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = sendmsg(fd, &msg, flags); err != nil {
		return 0, err
	}
	if len(oob) > 0 && len(p) == 0 {
		n = 0
	}
	return n, nil
}

//sys	utimes(path string, timeval *[2]Timeval) (err error)

func Utimes(path string, tv []Timeval) (err error) {
	if len(tv) != 2 {
		return EINVAL
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

func UtimesNano(path string, ts []Timespec) error {
	if len(ts) != 2 {
		return EINVAL
	}
	// Not as efficient as it could be because Timespec and
	// Timeval have different types in the different OSes
	tv := [2]Timeval{
		NsecToTimeval(TimespecToNsec(ts[0])),
		NsecToTimeval(TimespecToNsec(ts[1])),
	}
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

//sys	futimes(fd int, timeval *[2]Timeval) (err error)

func Futimes(fd int, tv []Timeval) (err error) {
	if len(tv) != 2 {
		return EINVAL
	}
	return futimes(fd, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
}

//sys	fcntl(fd int, cmd int, arg int) (val int, err error)

var mapper = &mmapper{
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
}

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) {
	return mapper.Mmap(fd, offset, length, prot, flags)
}

func Munmap(b []byte) (err error) {
	return mapper.Munmap(b)
}

// Implemented in the runtime package (runtime/sys_sylixos3.go)
func syscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscallX(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscallXnull(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscallXerrno(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func syscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func syscall6X(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func syscall6Xerrno(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
func rawSyscall(fn, a1, a2, a3 uintptr) (r1, r2 uintptr, err Errno)
func rawSyscall6(fn, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2 uintptr, err Errno)
