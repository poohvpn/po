package po

func Checksum(data []byte) uint16 {
	var csum uint32
	lastIndex := len(data) - 1
	for i := 0; i < lastIndex; i += 2 {
		csum += uint32(data[i])<<8 + uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[lastIndex]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum)
}
