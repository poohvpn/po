package pooh

import "net"

type (
	Header   byte
	AddrType byte
)

const (
	HeaderICMDP4    Header = 0x01
	HeaderTCP       Header = 0x06
	HeaderUDP       Header = 0x11
	HeaderICMDP6    Header = 0x3A
	HeaderKeepalive Header = 0x7A
	HeaderAuth      Header = 0xAF
	HeaderCountry   Header = 0xC0
	HeaderTravel    Header = 0xC2
	HeaderError     Header = 0xEE

	AddrTypeIPv4             AddrType = 0x04
	AddrTypeIPv6             AddrType = 0x06
	AddrTypeDomain           AddrType = 0x10
	AddrTypeDomainPreferIPv4 AddrType = 0x14
)

type Addr struct {
	Header Header
	Type   AddrType
	IP     net.IP
	Domain string
	Port   int
	Zone   string
}

var _ net.Addr = &Addr{}

func (a *Addr) Network() string {
	if a == nil {
		return ""
	}
	switch a.Header {
	case HeaderTCP:
		return "tcp"
	case HeaderUDP:
		return "udp"
	case HeaderICMDP4, HeaderICMDP6:
		return "icmdp"
	default:
		return ""
	}
}

func (a *Addr) String() string {
	switch {
	case a == nil:
		return ""
	case len(a.IP) > 0:
		return (&net.TCPAddr{
			IP:   a.IP,
			Port: a.Port,
			Zone: a.Zone,
		}).String()
	case a.Domain != "":
		return JoinHostPort(a.Domain, a.Port)
	default:
		return ""
	}
}

func NewAddr(addr net.Addr) *Addr {
	switch v := addr.(type) {
	case *Addr:
		return v
	case *net.TCPAddr:
		return &Addr{
			Header: HeaderTCP,
			IP:     v.IP,
			Port:   v.Port,
			Zone:   v.Zone,
		}
	case *net.UDPAddr:
		return &Addr{
			Header: HeaderUDP,
			IP:     v.IP,
			Port:   v.Port,
			Zone:   v.Zone,
		}
	default:
		udpAddr, err := net.ResolveUDPAddr("udp", v.String())
		if err != nil {
			panic(err)
		}
		return &Addr{
			IP:   udpAddr.IP,
			Port: udpAddr.Port,
			Zone: udpAddr.Zone,
		}
	}
}
