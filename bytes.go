package po

import "encoding/binary"

func Int(bs []byte) (i int) {
	for _, b := range bs {
		i = (i << 8) + int(b)
	}
	return
}

func Int2Bytes(v, size int) (bs []byte) {
	bs = make([]byte, size)
	for i := range bs {
		bs[size-i-1] = byte(v)
		v >>= 8
	}
	return
}

func Uint162Bytes(u ...uint16) []byte {
	bs := make([]byte, len(u)*2)
	for i, v := range u {
		binary.BigEndian.PutUint16(bs[i*2:], v)
	}
	return bs
}

func Uint322Bytes(u ...uint32) []byte {
	bs := make([]byte, len(u)*4)
	for i, v := range u {
		binary.BigEndian.PutUint32(bs[i*4:], v)
	}
	return bs
}

func Uint642Bytes(u ...uint64) []byte {
	bs := make([]byte, len(u)*8)
	for i, v := range u {
		binary.BigEndian.PutUint64(bs[i*8:], v)
	}
	return bs
}

func Duplicate(bs []byte) []byte {
	if bs == nil {
		return nil
	}
	bsDup := make([]byte, len(bs))
	if len(bs) > 0 {
		copy(bsDup, bs)
	}
	return bsDup
}
