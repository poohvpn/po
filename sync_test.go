package po

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestOnce(tt *testing.T) {
	t := require.New(tt)
	var once Once
	t.False(once.Done())
	ok := false
	once.Do(func() {
		ok = true
	})
	once.Do(func() {
		ok = false
	})
	t.True(ok)
	t.True(once.Done())
	select {
	case <-once.Wait():
	case <-time.After(time.Microsecond):
		t.Fail("once.Wait() should not block")
	}

	once = Once{}
	go func() {
		select {
		case <-once.Wait():
		case <-time.After(3 * time.Microsecond):
			t.Fail("once.Wait() should not block after once done")
		}
	}()
	time.Sleep(time.Microsecond)
	once.Do(func() {
		ok = false
	})
	once.Do(func() {
		ok = true
	})
	t.False(ok)
	t.True(once.Done())
	select {
	case <-once.Wait():
	case <-time.After(3 * time.Microsecond):
		t.Fail("once.Wait() should not block")
	}
}
