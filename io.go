package pooh

import (
	"io"
	"net"
	"time"
)

const (
	BufferSize = 64 * 1024
)

func Swap(rwc1, rwc2 io.ReadWriteCloser) {
	done := make(chan struct{}, 2)
	go func() {
		Copy(rwc2, rwc1)
		done <- struct{}{}
	}()
	go func() {
		Copy(rwc1, rwc2)
		done <- struct{}{}
	}()
	<-done
	_ = Close(rwc1)
	_ = Close(rwc2)
	<-done
}

func Copy(dst io.Writer, src io.Reader) {
	if wt, ok := src.(io.WriterTo); ok {
		_, _ = wt.WriteTo(dst)
		return
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		_, _ = rt.ReadFrom(src)
		return
	}
	buf := make([]byte, BufferSize)
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

func SwapTimeout(conn1, conn2 net.Conn) {
	done := make(chan struct{}, 2)
	go func() {
		CopyTimeout(conn2, conn1)
		done <- struct{}{}
	}()
	go func() {
		CopyTimeout(conn1, conn2)
		done <- struct{}{}
	}()
	<-done
	_ = Close(conn1)
	_ = Close(conn2)
	<-done
}

func CopyTimeout(dst, src net.Conn) {
	buf := make([]byte, BufferSize)
	for {
		err := src.SetReadDeadline(time.Now().Add(PacketTimeout))
		if err != nil {
			return
		}
		nr, er := src.Read(buf)
		if nr > 0 {
			err = dst.SetWriteDeadline(time.Now().Add(PacketTimeout))
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
}

func SwapPacketTimeout(conn1, conn2 net.PacketConn) {
	done := make(chan struct{}, 2)
	go func() {
		CopyPacketTimeout(conn2, conn1)
		done <- struct{}{}
	}()
	go func() {
		CopyPacketTimeout(conn1, conn2)
		done <- struct{}{}
	}()
	<-done
	_ = Close(conn1)
	_ = Close(conn2)
	<-done
}

func CopyPacketTimeout(dst, src net.PacketConn) {
	buf := make([]byte, BufferSize)
	for {
		err := src.SetReadDeadline(time.Now().Add(PacketTimeout))
		if err != nil {
			return
		}
		nr, addr, er := src.ReadFrom(buf)
		if nr > 0 {
			err = dst.SetWriteDeadline(time.Now().Add(PacketTimeout))
			if err != nil {
				return
			}
			nw, ew := dst.WriteTo(buf[:nr], addr)
			if ew != nil || nw != nr {
				return
			}
		}
		if er != nil {
			return
		}
	}
}

func Close(c io.Closer) error {
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
