package edf

func uint64ToBytes(in uint64, out []byte) {
	var i uint64
	for i = 0; i < 8; i++ {
		out[7-i] = byte((in & (0xFF << i * 8) >> i * 8))
	}
}

func uint64FromBytes(in []byte) uint64 {
	var i uint64
	out := uint64(0)
	for i = 0; i < 8; i++ {
		out |= uint64(in[7-i] << uint64(i*0x8))
	}
	return out
}
func uint32FromBytes(in []byte) uint32 {
	ret := uint32(0)
	ret |= uint32(in[0]) << 24
	ret |= uint32(in[1]) << 16
	ret |= uint32(in[2]) << 8
	ret |= uint32(in[3])
	return ret
}

func uint32ToBytes(in uint32, out []byte) {
	out[0] = byte(in & (0xFF << 24) >> 24)
	out[1] = byte(in & (0xFF << 16) >> 16)
	out[2] = byte(in & (0xFF << 8) >> 8)
	out[3] = byte(in & 0xFF)
}
