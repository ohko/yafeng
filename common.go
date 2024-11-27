package yafeng

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"strings"
)

func Hash(data string) string {
	hash := sha512.Sum512([]byte(data))
	return hex.EncodeToString(hash[:])
}

func GenerateNonce(n int) string {
	if n == 0 {
		n = 32
	}
	b := make([]byte, n/2)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return strings.ToUpper(hex.EncodeToString(b))
}

func InArray[T bool | string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](
	arr []T, x T) bool {
	for _, v := range arr {
		if x == v {
			return true
		}
	}
	return false
}
