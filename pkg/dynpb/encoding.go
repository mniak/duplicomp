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

func DecodeZigZag(v uint64) int64 {
	if v%2 == 0 {
		return int64(v / 2)
	} else {
		return int64((v+1)/2) * -1
	}
}

func DecodeTwosComplement(v uint64) int64 {
	return int64(v)
}

func DecodeFloat(v uint64) float32 {
	f := math.Float32frombits(uint32(v))
	return f
}

func DecodeDouble(v uint64) float64 {
	f := math.Float64frombits(uint64(v))
	return f
}
