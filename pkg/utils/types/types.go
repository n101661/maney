package types

import (
	"fmt"
)

func Int32ToBytes(v int32) []byte {
	res := make([]byte, 4)
	for i := range 4 {
		res[3-i] = byte(v)
		v >>= 8
	}
	return res
}

func BytesToInt32(v []byte) (int32, error) {
	if len(v) > 4 {
		return 0, fmt.Errorf("invalid int32: %v", v)
	}
	res := int32(0)
	for i, b := range v {
		res |= int32(b) << ((3 - i) * 8)
	}
	return res, nil
}
