package pooh

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConnTest struct {
	suite.Suite
	conn Conn
	zero net.Conn
	one  net.Conn
}

func (t *ConnTest) SetupTest() {
	t.zero, t.one = net.Pipe()
	t.conn = NewConn(t.zero)
}

func (t *ConnTest) TestB() {
	go func() {
		t.Require().NoError(Errors(t.one.Write([]byte("h"))))
	}()

	b, err := t.conn.Byte()
	t.NoError(err)
	t.Equal(byte('h'), b)
}

func (t *ConnTest) TestN() {
	go func() {
		t.Require().NoError(Errors(t.one.Write([]byte("hello"))))
	}()

	bs, err := t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("hello"), bs)
}

func (t *ConnTest) TestDontFragment() {
	t.conn.WithDontFragment()
	go func() {
		t.Require().NoError(Errors(t.one.Write([]byte("hell"))))
		t.Require().NoError(Errors(t.one.Write([]byte("o"))))
	}()

	bs, err := t.conn.Bytes(5)
	t.EqualError(err, "pooh.Conn: should not fragment")
	t.Equal([]byte("hell"), bs)

	t.conn.Reset().WithDontFragment(false)
	bs, err = t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("hello"), bs)
}

func (t *ConnTest) TestReset() {
	go func() {
		t.Require().NoError(Errors(t.one.Write([]byte("hello"))))
	}()

	bs, err := t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("hello"), bs)

	t.conn.Reset()
	bs, err = t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("hello"), bs)
}

func (t *ConnTest) TestPreplace() {
	go func() {
		t.Require().NoError(Errors(t.one.Write([]byte("hello"))))
	}()

	bs, err := t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("hello"), bs)

	t.conn.Reset()
	bs, err = t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("hello"), bs)

	t.conn.Preplace([]byte("world"))
	bs, err = t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("world"), bs)

	t.conn.Reset()
	bs, err = t.conn.Bytes(5)
	t.NoError(err)
	t.Equal([]byte("world"), bs)

	t.conn.Preplace([]byte{})
	go func() {
		t.NoError(Errors(t.one.Write([]byte("hello world"))))
	}()
	bs, err = t.conn.Bytes(6)
	t.NoError(err)
	t.Equal([]byte("hello "), bs)
	t.conn.Preplace([]byte("new "))
	bs, err = t.conn.Bytes(9)
	t.NoError(err)
	t.Equal([]byte("new world"), bs)

	t.conn.Preplace([]byte("end"))
	t.conn.Preplace(make([]byte, 8192))
	t.conn.Preplace([]byte("start"))
	bs, err = t.conn.Bytes(8200)
	t.NoError(err)
	t.Equal(append([]byte("start"), append(make([]byte, 8192), []byte("end")...)...), bs)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConnTest))
}

func TestIsIPv4(tt *testing.T) {
	t := assert.New(tt)
	t.False(IsIPv4(net.IP(nil)))
	t.False(IsIPv4(net.IPv6unspecified))
	t.True(IsIPv4(net.IPv4(1, 1, 1, 1)))
}

func TestIsIPv6(tt *testing.T) {
	t := assert.New(tt)
	t.False(IsIPv6(net.IP(nil)))
	t.False(IsIPv6(net.IPv4(1, 1, 1, 1)))
	t.True(IsIPv6(net.IPv6unspecified))
}
