package pooh

import (
	"github.com/rs/zerolog/log"
	"net"
	"syscall"
)

const (
	ProtoIcmp   byte = syscall.IPPROTO_ICMP
	ProtoIcmpV6 byte = syscall.IPPROTO_ICMPV6
	ProtoTcp    byte = syscall.IPPROTO_TCP
	ProtoUdp    byte = syscall.IPPROTO_UDP
)

const (
	AddrTypeDomain byte = 0x01
	AddrTypeIPv4   byte = 0x04
	AddrTypeIPv6   byte = 0x06
)

func NewAnalog(src net.Conn) *Analog {
	raddr := src.RemoteAddr()
	a := &Analog{
		SrcConn:  NewConn(src),
		SrcProto: raddr.Network(),
	}
	switch addr := raddr.(type) {
	case *net.TCPAddr:
		a.SrcIP = addr.IP
		a.SrcPort = addr.Port
		a.SrcZone = addr.Zone
	case *net.UDPAddr:
		a.SrcIP = addr.IP
		a.SrcPort = addr.Port
		a.SrcZone = addr.Zone
	default:
		udpAddr, err := net.ResolveUDPAddr("udp", addr.String())
		if err != nil {
			panic(err)
		}
		a.SrcIP = udpAddr.IP
		a.SrcPort = udpAddr.Port
		a.SrcZone = udpAddr.Zone
	}
	return a
}

type Analog struct {
	SrcConn          Conn        `json:"-"`
	SrcProto         string      `json:"src_proto"`
	SrcIP            net.IP      `json:"src_ip"`
	SrcPort          int         `json:"src_port"`
	SrcZone          string      `json:"src_zone"`
	SrcProxyProtocol string      `json:"src_proxy_protocol"`
	DstProto         byte        `json:"dst_proto"`
	DstAddrType      byte        `json:"dst_addr_type"`
	DstIP            net.IP      `json:"dst_ip"`
	DstDomain        string      `json:"dst_domain"`
	DstPort          int         `json:"dst_port"`
	Demodulator      Demodulator `json:"-"`
	Context          Map         `json:"context"`
	DstConn          net.Conn    `json:"-"`
}

func (a *Analog) Fork(src net.Conn) *Analog {
	res := *a
	res.SrcConn = NewConn(src)
	return &res
}

func (a *Analog) Simulate() {
	defer func() {
		_ = a.Close()
	}()
	switch a.DstProto {
	case ProtoTcp:
		Swap(a.SrcConn, a.DstConn)
	case ProtoUdp, ProtoIcmp, ProtoIcmpV6:
		SwapTimeout(a.SrcConn, a.DstConn)
	default:
		log.Panic().Uint8("dst_proto", a.DstProto).Msg("unknown destination protocol")
	}
}

func (a *Analog) Close() error {
	return Errors(Close(a.SrcConn), Close(a.DstConn))
}

type Modulator interface {
	Modulate(a *Analog) (err error, ok bool)
	Proto() string
}

type Demodulator interface {
	Demodulate(a *Analog) (err error, ok bool)
}
