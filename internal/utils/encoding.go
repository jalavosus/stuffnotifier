package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type BytesLike interface{ string | []byte }

func bytesLikeToBytes[T BytesLike](data T) (b []byte) {
	switch val := (any)(data).(type) {
	case []byte:
		b = val
	case string:
		b = []byte(val)
	}

	return
}

// EncodeB64 encodes arbitrary string or bytes into a base64 string.
func EncodeB64[T BytesLike](data T) string {
	return base64.StdEncoding.EncodeToString(bytesLikeToBytes(data))
}

// DecodeB64 decodes data stored as a base64 string into
// bytes.
func DecodeB64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func Sha256[T BytesLike](data T) (sum []byte) {
	sum = make([]byte, 32)

	shaSum := sha256.Sum256(bytesLikeToBytes(data))
	copy(sum, shaSum[:])

	return
}

func Sha256Hex[T BytesLike](data T) string {
	shaSum := Sha256(data)
	hash := fmt.Sprintf("%x", shaSum)

	return hash
}
