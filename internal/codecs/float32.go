package codecs

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/goccy/go-reflect"
)

const (
	typeByteFloat32        byte = 0x30
	typeByteFloat32Inverse      = typeByteFloat32 ^ 0xff
	sizeFloat32                 = int(unsafe.Sizeof(float32(0))) + 1
)

func EncodeFloat32(b []byte, v float32, inverse bool) int {
	if cap(b) < sizeFloat32 {
		panic("slice is too small to hold a float32")
	}
	b = b[:sizeFloat32]
	b[0] = typeByteFloat32
	int32Val := int32(math.Float32bits(v))
	int32Val ^= (int32Val >> 31) | (-1 << 31)
	binary.BigEndian.PutUint32(b[1:], uint32(int32Val))
	if inverse {
		invertArray(b)
	}
	return sizeFloat32
}

func DecodeFloat32(b []byte, v reflect.Value) (int, error) {
	encoded := b[1:5]
	if b[0] == typeByteFloat32Inverse {
		encoded = make([]byte, 4)
		copy(encoded, b[1:5])
		invertArray(encoded)
	}
	val := int32(binary.BigEndian.Uint32(encoded))
	val ^= (^val >> 31) | (-1 << 31)
	fv := math.Float32frombits(uint32(val))
	ptr := v.Pointer()
	**(**float32)(unsafe.Pointer(&ptr)) = *(*float32)(unsafe.Pointer(&fv))
	return sizeFloat32, nil
}
