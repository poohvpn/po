package pooh

import (
	"net"

	"github.com/rs/zerolog/log"
)

type Analog struct {
	SrcConn   Conn
	SrcProto  string
	SrcAddr   *Addr
	Context   Object
	DstAddr   *Addr
	Modulator Modulator
	DstConn   net.Conn
}

func NewAnalog(src net.Conn) *Analog {
	a := &Analog{
		SrcConn: NewConn(src),
		SrcAddr: NewAddr(src.RemoteAddr()),
	}
	return a
}

func (a *Analog) Fork(src net.Conn) *Analog {
	res := *a
	res.SrcConn = NewConn(src)
	return &res
}

func (a *Analog) Accepted() bool {
	return a.SrcProto != ""
}

func (a *Analog) Simulate() {
	defer func() {
		_ = a.Close()
	}()
	switch a.DstAddr.Header {
	case HeaderTCP:
		Swap(a.SrcConn, a.DstConn)
	case HeaderUDP, HeaderICMDP4, HeaderICMDP6:
		Swap(a.SrcConn, a.DstConn, CopyOption{
			Timeout: PacketTimeout,
		})
	default:
		log.Panic().Interface("dst_header", a.DstAddr.Header).Msg("unknown destination protocol")
	}
}

func (a *Analog) Close() error {
	return Errors(Close(a.SrcConn), Close(a.DstConn))
}
