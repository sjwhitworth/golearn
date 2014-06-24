package base

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

func PackU64ToBytesInline(val uint64, ret []byte) {
	ret[7] = byte(val & (0xFF << 56) >> 56)
	ret[6] = byte(val & (0xFF << 48) >> 48)
	ret[5] = byte(val & (0xFF << 40) >> 40)
	ret[4] = byte(val & (0xFF << 32) >> 32)
	ret[3] = byte(val & (0xFF << 24) >> 24)
	ret[2] = byte(val & (0xFF << 16) >> 16)
	ret[1] = byte(val & (0xFF << 8) >> 8)
	ret[0] = byte(val & (0xFF << 0) >> 0)
}

func PackFloatToBytesInline(val float64, ret []byte) {
	PackU64ToBytesInline(math.Float64bits(val), ret)
}

func PackU64ToBytes(val uint64) []byte {
	ret := make([]byte, 8)
	ret[7] = byte(val & (0xFF << 56) >> 56)
	ret[6] = byte(val & (0xFF << 48) >> 48)
	ret[5] = byte(val & (0xFF << 40) >> 40)
	ret[4] = byte(val & (0xFF << 32) >> 32)
	ret[3] = byte(val & (0xFF << 24) >> 24)
	ret[2] = byte(val & (0xFF << 16) >> 16)
	ret[1] = byte(val & (0xFF << 8) >> 8)
	ret[0] = byte(val & (0xFF << 0) >> 0)
	return ret
}

func UnpackBytesToU64(val []byte) uint64 {
	pb := unsafe.Pointer(&val[0])
	return *(*uint64)(pb)
}

func PackFloatToBytes(val float64) []byte {
	return PackU64ToBytes(math.Float64bits(val))
}

func UnpackBytesToFloat(val []byte) float64 {
	pb := unsafe.Pointer(&val[0])
	return *(*float64)(pb)
}

func xorFloatOp(item float64) float64 {
	var ret float64
	var tmp int64
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, item)
	binary.Read(buf, binary.LittleEndian, &tmp)
	tmp ^= -1 << 63
	binary.Write(buf, binary.LittleEndian, tmp)
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret
}

func printFloatByteArr(arr [][]byte) {
	buf := bytes.NewBuffer(nil)
	var f float64
	for _, b := range arr {
		buf.Write(b)
		binary.Read(buf, binary.LittleEndian, &f)
		f = xorFloatOp(f)
		fmt.Println(f)
	}
}
