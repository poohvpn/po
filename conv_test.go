package pooh

import (
	"github.com/stretchr/testify/require"
	"testing"
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
