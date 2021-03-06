package po

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	PacketTimeout  = 30 * time.Second
	HttpBufferSize = 8 * 1024
)

func NewConn(c net.Conn) Conn {
	if conn, ok := c.(Conn); ok {
		return conn
	}
	conn := &connImpl{
		Conn: c,
	}
	return conn.WithBuffer()
}

type Conn interface {
	net.Conn
	Byte() (byte, error)
	Bytes(n int) ([]byte, error)
	BytesString(n int) (string, error)
	Int(size int) (i int, err error)
	Uint16() (u uint16, err error)
	Uint32() (u uint32, err error)
	Uint64() (u uint64, err error)
	Frame(size int) ([]byte, error)
	FrameString(size int) (string, error)
	WriteFrame(p []byte, size int) error
	WriteUint16(u uint16) error
	WriteUint32(u uint32) error
	WriteUint64(u uint64) error
	Preplace(replace []byte) Conn
	Reset() Conn
	WithBuffer(size ...int) Conn
	WithDontFragment(dontFragment ...bool) Conn
}

type connImpl struct {
	net.Conn
	buf          []byte
	index        int
	dontFragment bool
	init         bool
}

func (c *connImpl) Read(b []byte) (n int, err error) {
	switch {
	case c.buf == nil:
		return c.Conn.Read(b)
	case !c.init:
		c.init = true
		var nr int
		nr, err = c.Conn.Read(c.buf[:cap(c.buf)])
		if err != nil {
			return 0, err
		}
		c.buf = c.buf[:nr]
		fallthrough
	case c.index < len(c.buf):
		n = copy(b, c.buf[c.index:])
		c.index += n
	default:
		n, err = c.Conn.Read(b)
		if err != nil {
			return
		}
		c.index += n
		if c.index > cap(c.buf) {
			c.buf = nil
			c.index = 0
		} else {
			c.buf = append(c.buf, b[:n]...)
		}
	}
	return
}

func (c *connImpl) Close() error {
	return Close(c.Conn)
}

func (c *connImpl) Byte() (byte, error) {
	b, err := c.Bytes(1)
	return b[0], err
}

func (c *connImpl) Bytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, nil
	}
	b := make([]byte, n)
	if c.dontFragment {
		nr, err := c.Read(b)
		if err != nil {
			return nil, err
		}
		if n != nr {
			return b[:nr], errors.New("pooh.Conn: should not fragment")
		}
	} else {
		_, err := io.ReadFull(c, b)
		if err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (c *connImpl) BytesString(n int) (string, error) {
	bs, err := c.Bytes(n)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func (c *connImpl) Int(size int) (i int, err error) {
	bs, err := c.Bytes(size)
	if err != nil {
		return
	}
	i = Int(bs)
	return
}

func (c *connImpl) Uint16() (u uint16, err error) {
	bs, err := c.Bytes(2)
	if err != nil {
		return
	}
	u = binary.BigEndian.Uint16(bs)
	return
}

func (c *connImpl) Uint32() (u uint32, err error) {
	bs, err := c.Bytes(4)
	if err != nil {
		return
	}
	u = binary.BigEndian.Uint32(bs)
	return
}

func (c *connImpl) Uint64() (u uint64, err error) {
	bs, err := c.Bytes(8)
	if err != nil {
		return
	}
	u = binary.BigEndian.Uint64(bs)
	return
}

func (c *connImpl) Frame(size int) ([]byte, error) {
	bs, err := c.Bytes(size)
	if err != nil {
		return nil, err
	}
	return c.Bytes(Int(bs))
}

func (c *connImpl) FrameString(size int) (string, error) {
	bs, err := c.Bytes(size)
	if err != nil {
		return "", err
	}
	return c.BytesString(Int(bs))
}

func (c *connImpl) WriteFrame(p []byte, size int) error {
	_, err := c.Conn.Write(append(Int2Bytes(len(p), size), p...))
	return err
}

func (c *connImpl) WriteUint16(u uint16) error {
	_, err := c.Conn.Write(Uint162Bytes(u))
	return err
}

func (c *connImpl) WriteUint32(u uint32) error {
	_, err := c.Conn.Write(Uint322Bytes(u))
	return err
}

func (c *connImpl) WriteUint64(u uint64) error {
	_, err := c.Conn.Write(Uint642Bytes(u))
	return err
}

func (c *connImpl) Preplace(replace []byte) Conn {
	need := len(replace) + len(c.buf) - c.index
	if need > cap(c.buf) {
		buf := make([]byte, 0, need)
		buf = append(buf, replace...)
		buf = append(buf, c.buf[c.index:]...)
		c.buf = buf
	} else {
		copy(c.buf[len(replace):need], c.buf[c.index:])
		copy(c.buf, replace)
		c.buf = c.buf[:need]
	}
	return c.Reset()
}

func (c *connImpl) Reset() Conn {
	c.index = 0
	return c
}

func (c *connImpl) WithBuffer(size ...int) Conn {
	switch {
	case len(size) == 0:
		c.buf = make([]byte, HttpBufferSize)
	case size[0] > 0:
		c.buf = make([]byte, size[0])
	default:
		c.buf = nil
	}
	return c
}

func (c *connImpl) WithDontFragment(dontFragment ...bool) Conn {
	if len(dontFragment) == 0 || dontFragment[0] {
		c.dontFragment = true
	} else {
		c.dontFragment = false
	}
	return c
}

func JoinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.Itoa(port))
}

func NewListener(addr net.Addr) *Listener {
	return &Listener{
		addr:    addr,
		connCh:  make(chan net.Conn),
		closeCh: make(chan struct{}),
	}
}

type Listener struct {
	addr      net.Addr
	connCh    chan net.Conn
	closeCh   chan struct{}
	closeOnce Once
}

var _ net.Listener = &Listener{}

func (l *Listener) Send(conn net.Conn) {
	if conn == nil {
		return
	}
	select {
	case l.connCh <- conn:
	case <-time.After(3 * time.Second):
		conn.Close()
	}
}

func (l *Listener) Accept() (net.Conn, error) {
	select {
	case conn := <-l.connCh:
		return conn, nil
	case <-l.closeCh:
		return nil, net.ErrClosed
	}
}

func (l *Listener) Close() error {
	l.closeOnce.Do(func() {
		close(l.closeCh)
	})
	return nil
}

func (l *Listener) Addr() net.Addr {
	return l.addr
}
