package pooh

import (
	"testing"
)

func TestRand(t *testing.T) {
	t.Log(Rand.Uint32())
	t.Log(Rand.Uint64())
}
