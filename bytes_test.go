package pooh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInt(tt *testing.T) {
	t := require.New(tt)
	for i := 0; i < 100; i++ {
		t.Equal(0, Int(make([]byte, i)))
	}
	t.Equal(1, Int([]byte{1}))
	t.Equal(0x1234, Int([]byte{0x12, 0x34}))
	t.Equal(0x1234, Int([]byte{0x00, 0x00, 0x12, 0x34}))
	t.Equal(0x12345678, Int([]byte{0x12, 0x34, 0x56, 0x78}))
	t.Equal(0x1234567890abcdef, Int([]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}))
}

func TestInt2Bytes(tt *testing.T) {
	t := require.New(tt)
	for i := 0; i < 100; i++ {
		t.Equal(make([]byte, i), Int2Bytes(0, i))
	}
	t.Equal([]byte{1}, Int2Bytes(1, 1))
	t.Equal([]byte{0x12, 0x34}, Int2Bytes(0x1234, 2))
	t.Equal([]byte{0x00, 0x00, 0x12, 0x34}, Int2Bytes(0x1234, 4))
	t.Equal([]byte{0x56, 0x78}, Int2Bytes(0x12345678, 2))
	t.Equal([]byte{0x12, 0x34, 0x56, 0x78}, Int2Bytes(0x12345678, 4))
	t.Equal([]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}, Int2Bytes(0x1234567890abcdef, 8))
}

func TestUint2Bytes(tt *testing.T) {
	t := require.New(tt)
	t.Equal([]byte{0, 0}, Uint162Bytes(0))
	t.Equal([]byte{0, 0, 0, 0}, Uint322Bytes(0))
	t.Equal([]byte{0, 0, 0, 0, 0, 0, 0, 0}, Uint642Bytes(0))
	t.Equal([]byte{0xff, 0xff}, Uint162Bytes(0xffff))
	t.Equal([]byte{0xff, 0xff, 0xff, 0xff}, Uint322Bytes(0xffffffff))
	t.Equal([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, Uint642Bytes(0xffffffffffffffff))
}

func TestDuplicate(tt *testing.T) {
	t := require.New(tt)
	t.Equal([]byte(nil), Duplicate(nil))
	t.Equal([]byte{}, Duplicate([]byte{}))
	t.NotEqual(nil, Duplicate([]byte{}))
	t.Equal([]byte{1, 2, 3}, Duplicate([]byte{1, 2, 3}))
}
