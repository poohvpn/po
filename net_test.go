package pooh

import (
	"github.com/stretchr/testify/suite"
	"net"
	"testing"
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

func (t *ConnTest) TestPreload() {
	go func() {
		t.Require().NoError(Errors(t.one.Write([]byte("hello"))))
	}()

	bs, err := t.conn.Preload()
	t.NoError(err)
	t.Equal([]byte("hello"), bs)
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
