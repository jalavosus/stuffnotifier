package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/sha3"
)

// BytesLike is any string or byteslice.
type BytesLike interface{ ~string | []byte }

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

// SHA3 returns the SHA3-256 sum of the passed bytes/string,
// as bytes.
func SHA3[T BytesLike](data T) (sum []byte) {
	sum = make([]byte, 32)
	shaSum := sha3.Sum256(bytesLikeToBytes(data))
	copy(sum, shaSum[:])

	return
}

// SHA3Hex returns the result of SHA3 as a
// human-readable string.
func SHA3Hex[T BytesLike](data T) (sumHex string) {
	shaSum := SHA3(data)
	sumHex = fmt.Sprintf("%x", shaSum)

	return
}

// SHA256 returns the SHA256 sum of the passed bytes/string
// as bytes.
func SHA256[T BytesLike](data T) (sum []byte) {
	sum = make([]byte, 32)
	shaSum := sha256.Sum256(bytesLikeToBytes(data))
	copy(sum, shaSum[:])

	return
}

// SHA256Hex returns the result of SHA256 as a
// human-readable string.
func SHA256Hex[T BytesLike](data T) (sumHex string) {
	shaSum := SHA256(data)
	sumHex = fmt.Sprintf("%x", shaSum)

	return
}
