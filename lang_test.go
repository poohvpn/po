package pooh

import (
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestIsNil(tt *testing.T) {
	t := require.New(tt)
	t.True(IsNil(nil))
	t.True(IsNil(net.Conn(nil)))
	t.True(IsNil([]byte(nil)))
}

func TestDuplicate(tt *testing.T) {
	t := require.New(tt)
	t.Equal([]byte(nil), Duplicate(nil))
	t.Equal([]byte{}, Duplicate([]byte{}))
	t.NotEqual(nil, Duplicate([]byte{}))
	t.Equal([]byte{1, 2, 3}, Duplicate([]byte{1, 2, 3}))
}
