package pooh

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
