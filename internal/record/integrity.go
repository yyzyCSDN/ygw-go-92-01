package record

import (
	"encoding/hex"

	"github.com/cespare/xxhash/v2"
)

func Digest(data []byte) string {
	sum := xxhash.Sum64(data)
	return hex.EncodeToString(uint64ToBytes(sum))
}

func uint64ToBytes(value uint64) []byte {
	buf := make([]byte, 8)
	for i := 0; i < 8; i++ {
		buf[i] = byte(value >> (8 * i))
	}
	return buf
}
