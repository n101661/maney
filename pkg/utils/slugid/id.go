package slugid

import (
	"math/bits"
	"math/rand/v2"
	"unsafe"
)

const alphaNumbers = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
	letterBits = 0
	mask       = uint64(0)
)

// New create an random ID with prefix and the size of random string.
func New(prefix string, size int) string {
	res := make([]byte, len(prefix)+1+size)

	// Copy the string of prefix.
	for i := range len(prefix) {
		res[i] = prefix[i]
	}
	// Use underscore as separator.
	res[len(prefix)] = '_'

	for i, v := len(prefix)+1, rand.Uint64(); i < len(res); {
		if v == 0 {
			v = rand.Uint64()
		}

		if masked := int(v & mask); masked < len(alphaNumbers) {
			res[i] = alphaNumbers[masked]
			i++
		}

		v >>= letterBits
	}

	return unsafe.String(unsafe.SliceData(res), len(res))
}

func init() {
	letterBits = bits.Len(uint(len(alphaNumbers)))
	mask = 1<<letterBits - 1
}
