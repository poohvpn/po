package po

import (
	"io"
	"net"
	"time"

	"github.com/hashicorp/go-multierror"
)

const (
	BufferSize = 64 * 1024
)

type CopyOption struct {
	BufferSize int
	Timeout    time.Duration
}

func Swap(rwc1, rwc2 net.Conn, option ...CopyOption) {
	done := make(chan struct{})
	go func() {
		Copy(rwc2, rwc1, option...)
		done <- struct{}{}
	}()
	go func() {
		Copy(rwc1, rwc2, option...)
		done <- struct{}{}
	}()
	<-done
	_ = Close(rwc1)
	_ = Close(rwc2)
	<-done
}

func Copy(dst, src net.Conn, option ...CopyOption) {
	var opt CopyOption
	if len(option) > 0 {
		opt = option[0]
	}
	if opt.BufferSize <= 0 {
		opt.BufferSize = BufferSize
	}
	if opt.Timeout <= 0 {
		if wt, ok := src.(io.WriterTo); ok {
			_, _ = wt.WriteTo(dst)
			return
		}
		// Similarly, if the writer has a ReadFrom method, use it to do the copy.
		if rt, ok := dst.(io.ReaderFrom); ok {
			_, _ = rt.ReadFrom(src)
			return
		}
	}
	buf := make([]byte, opt.BufferSize)
	if opt.Timeout > 0 {
		for {
			err := src.SetReadDeadline(time.Now().Add(opt.Timeout))
			if err != nil {
				return
			}
			nr, er := src.Read(buf)
			if nr > 0 {
				err = dst.SetWriteDeadline(time.Now().Add(opt.Timeout))
				if err != nil {
					return
				}
				nw, ew := dst.Write(buf[:nr])
				if ew != nil || nw != nr {
					return
				}
			}
			if er != nil {
				return
			}
		}
	} else {
		for {
			nr, er := src.Read(buf)
			if nr > 0 {
				nw, ew := dst.Write(buf[0:nr])
				if ew != nil || nr != nw {
					return
				}
			}
			if er != nil {
				return
			}
		}
	}
}

func tryClose(c io.Closer) error {
	if c == nil || IsNil(c) {
		return nil
	}
	if closer, ok := c.(interface{ CloseRead() error }); ok {
		_ = closer.CloseRead()
	}
	if closer, ok := c.(interface{ CloseWrite() error }); ok {
		_ = closer.CloseWrite()
	}
	return c.Close()
}

func Close(cs ...io.Closer) error {
	errs := new(multierror.Error)
	for _, c := range cs {
		errs = multierror.Append(errs, tryClose(c))
	}
	return errs.ErrorOrNil()
}
