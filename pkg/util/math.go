package util

import (
	"encoding/hex"
	"fmt"
	"math"
)

func Unsynchsafe(in uint32) uint32 {
	var out uint32 = 0
	mask := uint32(0x7F000000)

	for {
		out >>= 1
		out |= in & mask
		mask >>= 8
		if mask == 0 {
			break
		}
	}

	return out
}

func Synchsafe(in uint32) uint32 {
	var out uint32 = 0
	mask := uint32(0x7F)

	for {
		out = in & ^mask
		out = out << 1
		out = out | in&mask
		mask = ((mask + 1) << 8) - 1
		in = out
		if mask^0x7FFFFFFF == 0 {
			break
		}
	}
	return out
}

func BoolFromBit(data []byte, bitIndex int) bool {
	byteIndex := bitIndex / 8
	bytes := make([]byte, 1)
	_ = copy(bytes, data[byteIndex:byteIndex+1])

	// shift
	inByteIndex := bitIndex % 8
	if inByteIndex != 0 {
		ShiftLeft(bytes, inByteIndex)
	}
	ShiftLeft(bytes, -7)
	return int8(bytes[0]) == 1
}

func BitsFromBytes(data []byte, bitIndex int, bitCount int) []byte {
	byteIndex := bitIndex / 8
	inByteIndex := bitIndex % 8
	byteIndexOffset := int(math.Ceil(float64(bitCount+inByteIndex) / 8))
	bytes := make([]byte, byteIndexOffset)
	_ = copy(bytes, data[byteIndex:byteIndex+byteIndexOffset])

	// shift left
	if inByteIndex != 0 {
		ShiftLeft(bytes, inByteIndex)
	}

	// shift right
	lastShift := 8 - bitCount
	if lastShift < 0 {
		lastShift = 8 - bitCount%8
	}
	ShiftLeft(bytes, -lastShift)

	// return relevant bytes
	bytesMax := int(math.Ceil(float64(bitCount) / 8))
	return bytes[:bytesMax]
}

func BytesEqualHex(h string, compare []byte) bool {
	// decode
	value, err := hex.DecodeString(h)
	if err != nil {
		fmt.Println("Error: Can't decode hex string.", err)
		return false
	}

	// validate
	if len(value) > len(compare) {
		fmt.Println("Hex compare data too short.")
		return false
	}

	// compare
	for i, b := range value {
		if compare[i] != b {
			// fmt.Println(int8(compare[i]) != int8(b))
			return false
		}
	}
	return true
}

func BytesEqualHexWithMask(h string, mask string, buffer []byte) bool {
	// decode hex strings
	value, err := hex.DecodeString(h)
	if err != nil {
		fmt.Println("Error: Can't decode hex string.", err)
		return false
	}
	valueMask, err := hex.DecodeString(mask)
	if err != nil {
		fmt.Println("Error: Can't decode mask.", err)
		return false
	}

	// validate
	if len(value) != len(valueMask) || len(value) > len(buffer) {
		fmt.Println("Hex compare data too short.")
		return false
	}

	// compare
	for i, b := range value {
		if buffer[i]&valueMask[i] != b {
			// fmt.Println(int8(compare[i]&valueMask[i]) != int8(b))
			return false
		}
	}
	return true
}

// ShiftLeft performs a left bit shift operation on the provided bytes.
// If the bits count is negative, a right bit shift is performed.
func ShiftLeft(data []byte, bits int) {
	n := len(data)
	if bits < 0 {
		// shift right
		bits = -bits
		for i := n - 1; i > 0; i-- {
			data[i] = data[i]>>bits | data[i-1]<<(8-bits)
		}
		data[0] >>= bits
	} else {
		// shift left
		for i := 0; i < n-1; i++ {
			data[i] = data[i]<<bits | data[i+1]>>(8-bits)
		}
		data[n-1] <<= bits
	}
}
