package pooh

import (
	"github.com/rs/zerolog/log"
	"net"
)

type Analog struct {
	SrcConn     Conn
	SrcProxy    string
	SrcAddr     *Addr
	Context     Map
	DstAddr     *Addr
	Demodulator Demodulator
	DstConn     net.Conn
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

func (a *Analog) Simulate() {
	defer func() {
		_ = a.Close()
	}()
	switch a.DstAddr.Header {
	case HeaderTCP:
		Swap(a.SrcConn, a.DstConn)
	case HeaderUDP, HeaderICMDP4, HeaderICMDP6:
		SwapTimeout(a.SrcConn, a.DstConn)
	default:
		log.Panic().Interface("dst_header", a.DstAddr.Header).Msg("unknown destination protocol")
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
