package datastore

import (
	"bytes"
	"encoding/gob"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/snappy"
	"github.com/pkg/errors"

	"github.com/stoicturtle/stuffnotifier/internal/utils"
)

// CompressSnappy compresses a byte slice using standard Snappy compression.
func CompressSnappy(data []byte) []byte {
	return snappy.Encode(nil, data)
}

// DecompressSnappy decompresses bytes which were compressed
// using standard Snappy compression.
func DecompressSnappy(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}

// CompressS2 compresses a byte slice using S2,
// an "improved" Snappy algorithm.
func CompressS2(data []byte) []byte {
	return s2.Encode(nil, data)
}

// DecompressS2 decompresses bytes which were compressed
// using the S2 Snappy algorithm.
func DecompressS2(data []byte) ([]byte, error) {
	return s2.Decode(nil, data)
}

// EncodeData encodes arbitrary data using encoding/gob,
// returning the byte-encoded result of gob.Encode.
func EncodeData[T any](data T) ([]byte, error) {
	var b bytes.Buffer

	enc := gob.NewEncoder(&b)
	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// DecodeData decodes gob-encoded bytes,
// returning the original data.
func DecodeData[T any](data []byte) (*T, error) {
	var (
		b       bytes.Buffer
		decoded T
	)

	if _, err := b.Write(data); err != nil {
		return nil, err
	}

	dec := gob.NewDecoder(&b)
	if err := dec.Decode(&decoded); err != nil {
		return nil, err
	}

	return &decoded, nil
}

func fullDecode[T any](dataBase64 string) (*T, bool, error) {
	bytesData, err := utils.DecodeB64(dataBase64)
	if err != nil {
		return nil, false, errors.Wrap(err, base64DecodeErrMsg)
	}

	decompressed, err := DecompressSnappy(bytesData)
	if err != nil {
		return nil, false, errors.Wrap(err, snappyDecodeErrMsg)
	}

	decoded, err := DecodeData[T](decompressed)
	if err != nil {
		return nil, false, errors.Wrap(err, gobDecodeErrMsg)
	}

	return decoded, decoded != nil, nil
}

func fullEncode[T any](data T) (string, error) {
	encoded, err := EncodeData[T](data)
	if err != nil {
		return "", errors.Wrap(err, gobEncodeErrMsg)
	}

	compressed := CompressSnappy(encoded)
	encodedStr := utils.EncodeB64(compressed)

	return encodedStr, nil
}
