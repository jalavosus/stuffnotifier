package datastore

const (
	base64DecodeErrMsg string = "error decoding base64 string"
	snappyDecodeErrMsg string = "error decompressing snappy-compressed bytes"
	gobDecodeErrMsg    string = "error decoding gob-encoded bytes"
)

const (
	// base64EncodeErrMsg string = "error encoding bytes to base64"
	// snappyEncodeErrMsg string = "error snappy-compressing bytes"
	gobEncodeErrMsg string = "error gob-encoding data"
)
