package pooh

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChecksum(tt *testing.T) {
	t := require.New(tt)
	cs := func(hexStr string) uint16 {
		data, err := hex.DecodeString(hexStr)
		t.NoError(err)
		return Checksum(data)
	}
	t.Equal(uint16(0x8b78), cs("00000000a91dc7365cc861240a090002ffffff000a0900010808080801010101051408bad5e789aafe821aca0aedc5538d2f3d"))
	t.Equal(uint16(0x8b78), cs("00000000a91dc7365cc861240a090002ffffff000a0900010808080801010101051408bad5e789aafe821aca0aedc5538d2f3d"))
	t.Equal(uint16(0xffff), cs("00000000"))
	t.Equal(uint16(0xb72f), cs("000000001234123412341234"))
}
