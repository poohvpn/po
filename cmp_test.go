package po

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNil(tt *testing.T) {
	t := require.New(tt)
	t.True(IsNil(nil))
	t.True(IsNil(net.Conn(nil)))
	t.True(IsNil([]byte(nil)))
}
