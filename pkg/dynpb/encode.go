package dynpb

import "math"

func EncodeZigZag(v int64) uint64 {
	if v >= 0 {
		return uint64(v * 2)
	} else {
		return uint64(v*-2 - 1)
	}
}

func EncodeTwosComplement(v int64) uint64 {
	return uint64(v)
}

func EncodeFloat(f float32) uint32 {
	b := math.Float32bits(f)
	return b
}

func EncodeDouble(f float64) uint64 {
	b := math.Float64bits(f)
	return b
}
