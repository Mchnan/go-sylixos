// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

/*
Input to cgo -godefs.  See also mkerrors.sh and mkall.sh
*/

// +godefs map struct_in_addr [4]byte /* in_addr */
// +godefs map struct_in6_addr [16]byte /* in6_addr */

package syscall

/*
#define SYLIXOS
#include <SylixOS.h>
#include <sys/socket.h>
#include <sys/resource.h>
#include <sys/un.h>
#include <sys/utsname.h>
#include <net/if_dl.h>

#include <netinet6/in6.h>
#include <netinet6/icmp6.h>
#include <net/if.h>
#include <net/route.h>
#include <termios.h>
#include <poll.h>
#include <spawn.h>

enum {
	sizeofPtr = sizeof(void*),
};

union sockaddr_all {
	struct sockaddr s1;	// this one gets used for fields
	struct sockaddr_in s2;	// these pad it out
	struct sockaddr_in6 s3;
	struct sockaddr_un s4;
	struct sockaddr_dl s5;
};

struct sockaddr_any {
	struct sockaddr addr;
	char pad[sizeof(union sockaddr_all) - sizeof(struct sockaddr)];
};
*/
import "C"

// Machine characteristics; for internal use.

const (
	sizeofPtr      = C.sizeofPtr
	sizeofShort    = C.sizeof_short
	sizeofInt      = C.sizeof_int
	sizeofLong     = C.sizeof_long
	sizeofLongLong = C.sizeof_longlong
)

// Basic types

type (
	_C_short     C.short
	_C_int       C.int
	_C_long      C.long
	_C_long_long C.longlong
)

// Time

type Timespec C.struct_timespec

type Timeval C.struct_timeval

// Processes

type Rusage C.struct_rusage

type Rlimit C.struct_rlimit

type _Gid_t C.gid_t

// Files

type Stat_t C.struct_stat

type Statfs_t C.struct_statfs

type Flock_t C.struct_flock

type Dirent C.struct_dirent

type Fsid C.fsid_t

// File system limits

const (
	pathMax = C.PATH_MAX
)

// Sockets

type RawSockaddrInet4 C.struct_sockaddr_in

type RawSockaddrInet6 C.struct_sockaddr_in6

type RawSockaddrUnix C.struct_sockaddr_un

type RawSockaddrDatalink C.struct_sockaddr_dl

type RawSockaddr C.struct_sockaddr

type RawSockaddrAny C.struct_sockaddr_any

type _Socklen C.socklen_t

type Linger C.struct_linger

type Iovec C.struct_iovec

type IPMreq C.struct_ip_mreq

type IPv6Mreq C.struct_ipv6_mreq

type Msghdr C.struct_msghdr

type Cmsghdr C.struct_cmsghdr

type Inet6Pktinfo C.struct_in6_pktinfo

type IPv6MTUInfo C.struct_ip6_mtuinfo

type ICMPv6Filter C.struct_icmp6_filter

const (
	SizeofSockaddrInet4    = C.sizeof_struct_sockaddr_in
	SizeofSockaddrInet6    = C.sizeof_struct_sockaddr_in6
	SizeofSockaddrAny      = C.sizeof_struct_sockaddr_any
	SizeofSockaddrUnix     = C.sizeof_struct_sockaddr_un
	SizeofSockaddrDatalink = C.sizeof_struct_sockaddr_dl
	SizeofLinger           = C.sizeof_struct_linger
	SizeofIPMreq           = C.sizeof_struct_ip_mreq
	SizeofIPv6Mreq         = C.sizeof_struct_ipv6_mreq
	SizeofMsghdr           = C.sizeof_struct_msghdr
	SizeofCmsghdr          = C.sizeof_struct_cmsghdr
	SizeofInet6Pktinfo     = C.sizeof_struct_in6_pktinfo
	SizeofIPv6MTUInfo      = C.sizeof_struct_ip6_mtuinfo
	SizeofICMPv6Filter     = C.sizeof_struct_icmp6_filter
)

// Select

type FdSet C.fd_set

// Routing and interface messages

const (
	SizeofIfMsghdr         = C.sizeof_struct_if_msghdr
	SizeofIfData           = C.sizeof_struct_if_data
	SizeofIfaMsghdr        = C.sizeof_struct_ifa_msghdr
	SizeofIfAnnounceMsghdr = C.sizeof_struct_if_announcemsghdr
	SizeofRtMsghdr         = C.sizeof_struct_rt_msghdr
	SizeofRtMetrics        = C.sizeof_struct_rt_metrics
)

type IfMsghdr C.struct_if_msghdr

type IfData C.struct_if_data

type IfaMsghdr C.struct_ifa_msghdr

type IfAnnounceMsghdr C.struct_if_announcemsghdr

type RtMsghdr C.struct_rt_msghdr

type RtMetrics C.struct_rt_metrics

// Misc

// Terminal handling

type Termios C.struct_termios

type Utsname C.struct_utsname

type posix_spawnopt_t C.posix_spawnopt_t

type sched_param C.struct_sched_param

//type posix_spawnattr_t C.posix_spawnattr_t
type posix_spawnattr_t struct {
	SFlags        int16
	PidGroup      int32
	SigsetDefault uint64
	SigsetMask    uint64
	Schedparam    sched_param
	IPolicy       int32
	PcWd          *uint8
	Presraw       *uint8
	Opt           posix_spawnopt_t
	UlExts        uint64
	UlPad         [5]uint64
}

//type posix_spawn_file_actions_t C.posix_spawn_file_actions_t
type posix_spawn_file_actions_t struct {
	PplineActions **uint8
	Presraw       *uint8
	IInited       int32
	UlPad         [16]uint64
}
