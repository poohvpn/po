package pooh

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

func Uint162Bytes(u uint16) []byte {
	var bs [2]byte
	binary.BigEndian.PutUint16(bs[:], u)
	return bs[:]
}

func Uint322Bytes(u uint32) []byte {
	var bs [4]byte
	binary.BigEndian.PutUint32(bs[:], u)
	return bs[:]
}

func Uint642Bytes(u uint64) []byte {
	var bs [8]byte
	binary.BigEndian.PutUint64(bs[:], u)
	return bs[:]
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
