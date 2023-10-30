package dynpb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAndDecode_TwosComplement(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		var value int64 = 2
		for i := 1; i < 64; i++ {
			t.Run(fmt.Sprintf("2 to the power of %d", i), func(t *testing.T) {
				encoded := EncodeTwosComplement(value)
				decoded := DecodeTwosComplement(encoded)

				assert.Equal(t, value, decoded)
			})
			value *= 2
		}
	})

	t.Run("Negative", func(t *testing.T) {
		var value int64 = -2
		for i := 1; i < 64; i++ {
			t.Run(fmt.Sprintf("2 to the power of %d", i), func(t *testing.T) {
				encoded := EncodeTwosComplement(value)
				decoded := DecodeTwosComplement(encoded)

				assert.Equal(t, value, decoded)
			})
			value *= 2
		}
	})
}

func TestEncodeAndDecode_ZigZag(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		var value int64 = 2
		for i := 1; i < 3; i++ {
			t.Run(fmt.Sprintf("2 to the power of %d", i), func(t *testing.T) {
				encoded := EncodeZigZag(value)
				decoded := DecodeZigZag(encoded)

				assert.Equal(t, value, decoded)
			})
			value *= 2
		}
	})

	t.Run("Negative", func(t *testing.T) {
		var value int64 = -2
		for i := 1; i < 63; i++ {
			t.Run(fmt.Sprintf("2 to the power of %d", i), func(t *testing.T) {
				encoded := EncodeZigZag(value)
				decoded := DecodeZigZag(encoded)

				assert.Equal(t, value, decoded)
			})
			value *= 2
		}
	})
}
