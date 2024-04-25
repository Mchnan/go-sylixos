// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build sylixos

package net

import (
	"internal/poll"
	"syscall"
	"unsafe"
)

const (
	_IFF_UP          = 0x0001 /* Interface is enable                  */
	_IFF_BROADCAST   = 0x0002 /* Interface support broadcast          */
	_IFF_POINTOPOINT = 0x0004 /* Interface is point to point          */
	_IFF_RUNNING     = 0x0010 /* Interface is linked                  */
	_IFF_MULTICAST   = 0x0080 /* Interface support multicast          */
	_IFF_LOOPBACK    = 0x0100 /* Loop back interface                  */
	_IFF_NOARP       = 0x0200 /* Do not use ARP protocol              */
	_IFF_PROMISC     = 0x0400 /* Receive all packets                  */
	_IFF_ALLMULTI    = 0x0800 /* Receive all multicast packets        */

	_MAX_IF   = 32
	_MAX_IPV6 = 10
	_IFNAMSIZ = 16
)

type Ifreq struct {
	Ifrn [_IFNAMSIZ]byte
	Ifru [32]byte
}

type IfreqMtu struct {
	Ifrn    [_IFNAMSIZ]byte
	IfruMtu int32
}

type IfreqIfIndex struct {
	Ifrn        [_IFNAMSIZ]byte
	IfruIfIndex int32
}

type IfreqFlags struct {
	Ifrn      [_IFNAMSIZ]byte
	IfruFlags int16
}

type Ifconf struct {
	Len int32
	Req *Ifreq
}

type In6Ifraddr struct {
	Addr      [16]byte /* in6_addr */
	PrefixLen int32
}

type In6Ifreq struct {
	IfIndex int32
	Len     int32
	Addr6   *In6Ifraddr
}

func ioctlIfreq(fd int, req int, value *Ifreq) error {
	return syscall.IoctlPtr(fd, req, unsafe.Pointer(value))
}

func interfaceName(name [_IFNAMSIZ]byte) string {
	var i = 0

	for ; i < len(name); i++ {
		if name[i] == 0 {
			break
		}
	}
	return string(name[0:i])
}

// If the ifindex is zero, interfaceTable returns mappings of all
// network interfaces. Otherwise it returns a mapping of a specific
// interface.
func interfaceTable(ifindex int) ([]Interface, error) {
	var (
		req  [_MAX_IF]Ifreq
		conf Ifconf
	)

	sock, err := sysSocket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	defer poll.CloseFunc(sock)

	conf.Len = int32(unsafe.Sizeof(req))
	conf.Req = &req[0]
	err = syscall.IoctlPtr(sock, syscall.SIOCGIFCONF, unsafe.Pointer(&conf))
	if err != nil {
		return nil, err
	}

	if_num := int(conf.Len) / int(unsafe.Sizeof(req[0]))

	var ift []Interface

	for i := 0; i < if_num; i++ {
		var en Interface

		err := ioctlIfreq(sock, syscall.SIOCGIFINDEX, &req[i])
		if err != nil {
			return nil, err
		}
		en.Index = int(((*IfreqIfIndex)(unsafe.Pointer(&req[i]))).IfruIfIndex)

		if ifindex != 0 {
			if en.Index != ifindex {
				continue
			}
		}

		en.Name = interfaceName(req[i].Ifrn)

		err = ioctlIfreq(sock, syscall.SIOCGIFMTU, &req[i])
		if err != nil {
			return nil, err
		}
		en.MTU = int(((*IfreqMtu)(unsafe.Pointer(&req[i]))).IfruMtu)

		err = ioctlIfreq(sock, syscall.SIOCGIFHWADDR, &req[i])
		if err == nil {
			en.HardwareAddr = req[i].Ifru[2:8]
		}

		err = ioctlIfreq(sock, syscall.SIOCGIFFLAGS, &req[i])
		if err != nil {
			return nil, err
		}
		flags := ((*IfreqFlags)(unsafe.Pointer(&req[i]))).IfruFlags

		if flags&_IFF_UP != 0 {
			en.Flags |= FlagUp
		}
		if flags&_IFF_BROADCAST != 0 {
			en.Flags |= FlagBroadcast
		}
		if flags&_IFF_LOOPBACK != 0 {
			en.Flags |= FlagLoopback
			if en.MTU == 0 {
				en.MTU = 65536
			}
		}
		if flags&_IFF_POINTOPOINT != 0 {
			en.Flags |= FlagPointToPoint
		}
		if flags&_IFF_MULTICAST != 0 {
			en.Flags |= FlagMulticast
		}
		if flags&_IFF_RUNNING != 0 {
			en.Flags |= FlagRunning
		}

		ift = append(ift, en)
	}

	return ift, nil
}

func ifiAddrTable(name string, index int) ([]Addr, error) {
	sock, err := sysSocket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		return nil, err
	}
	defer poll.CloseFunc(sock)

	var (
		ifat []Addr
		req  Ifreq
		sa   *syscall.RawSockaddrInet4
	)

	copy(req.Ifrn[:], name)
	err = ioctlIfreq(sock, syscall.SIOCGIFADDR, &req)
	if err != nil {
		return nil, err
	}

	sa = (*syscall.RawSockaddrInet4)(unsafe.Pointer(&req.Ifru))
	ifat = append(ifat, &IPAddr{IP: IPv4(sa.Addr[0], sa.Addr[1], sa.Addr[2], sa.Addr[3])})

	var (
		req6  In6Ifreq
		addr6 [_MAX_IPV6]In6Ifraddr
	)

	req6.IfIndex = int32(index)
	req6.Len = int32(unsafe.Sizeof(addr6))
	req6.Addr6 = &addr6[0]
	err = syscall.IoctlPtr(sock, syscall.SIOCGIFADDR6, unsafe.Pointer(&req6))
	if err != nil {
		return nil, err
	}

	ipv6_num := int(req6.Len) / int(unsafe.Sizeof(addr6[0]))
	for i := 0; i < ipv6_num; i++ {
		ifa := &IPNet{IP: make(IP, IPv6len), Mask: CIDRMask(int(addr6[i].PrefixLen), 8*IPv6len)}
		copy(ifa.IP, addr6[i].Addr[:])
		ifat = append(ifat, ifa)
	}

	return ifat, nil
}

// If the ifi is nil, interfaceAddrTable returns addresses for all
// network interfaces. Otherwise it returns addresses for a specific
// interface.
func interfaceAddrTable(ifi *Interface) ([]Addr, error) {
	if ifi != nil {
		return ifiAddrTable(ifi.Name, ifi.Index)

	} else {
		var (
			req  [_MAX_IF]Ifreq
			conf Ifconf
		)

		sock, err := sysSocket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
		if err != nil {
			return nil, err
		}
		defer poll.CloseFunc(sock)

		conf.Len = int32(unsafe.Sizeof(req))
		conf.Req = &req[0]
		err = syscall.IoctlPtr(sock, syscall.SIOCGIFCONF, unsafe.Pointer(&conf))
		if err != nil {
			return nil, err
		}

		if_num := int(conf.Len) / int(unsafe.Sizeof(req[0]))

		var ifat []Addr

		for i := 0; i < if_num; i++ {
			name := interfaceName(req[i].Ifrn)

			err = ioctlIfreq(sock, syscall.SIOCGIFINDEX, &req[i])
			if err != nil {
				return nil, err
			}

			index := int(((*IfreqIfIndex)(unsafe.Pointer(&req[i]))).IfruIfIndex)

			ifa, err := ifiAddrTable(name, index)
			if err != nil {
				return nil, err
			}
			ifat = append(ifat, ifa...)
		}

		return ifat, nil
	}
}

// interfaceMulticastAddrTable returns addresses for a specific
// interface.
func interfaceMulticastAddrTable(ifi *Interface) ([]Addr, error) {
	ifmat4 := parseProcNetIGMP("/proc/net/igmp", ifi)
	ifmat6 := parseProcNetIGMP6("/proc/net/igmp6", ifi)
	return append(ifmat4, ifmat6...), nil
}

func parseProcNetIGMP(path string, ifi *Interface) []Addr {
	fd, err := open(path)
	if err != nil {
		return nil
	}
	defer fd.close()
	var ifmat []Addr
	fd.readLine() // skip first line
	b := make([]byte, IPv4len)
	for l, ok := fd.readLine(); ok; l, ok = fd.readLine() {
		f := splitAtBytes(l, " :\r\t\n")
		if len(f) < 3 {
			continue
		}
		if ifi == nil || f[0] == ifi.Name {
			// The SylixOS kernel puts the IP
			// address in /proc/net/igmp in native
			// endianness.
			for i := 0; i+1 < len(f[1]); i += 2 {
				b[i/2], _ = xtoi2(f[1][i:i+2], 0)
			}
			i := *(*uint32)(unsafe.Pointer(&b[:4][0]))
			ifma := &IPAddr{IP: IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i))}
			ifmat = append(ifmat, ifma)
		}
	}
	return ifmat
}

func parseProcNetIGMP6(path string, ifi *Interface) []Addr {
	fd, err := open(path)
	if err != nil {
		return nil
	}
	defer fd.close()
	var ifmat []Addr
	b := make([]byte, IPv6len)
	for l, ok := fd.readLine(); ok; l, ok = fd.readLine() {
		f := splitAtBytes(l, " \r\t\n")
		if len(f) < 3 {
			continue
		}
		if ifi == nil || f[0] == ifi.Name {
			for i := 0; i+1 < len(f[1]); i += 2 {
				b[i/2], _ = xtoi2(f[1][i:i+2], 0)
			}
			ifma := &IPAddr{IP: IP{b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8], b[9], b[10], b[11], b[12], b[13], b[14], b[15]}}
			ifmat = append(ifmat, ifma)
		}
	}
	return ifmat
}
