package decoder

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/binary"
)

// ZintToFloat decodes a base64-encoded, zlib-compressed string into a slice of float64 values.
//
// The function performs the following steps:
// 1. Decodes the input string from base64 encoding.
// 2. Decompresses the decoded data using zlib.
// 3. Reads the decompressed data as a sequence of float64 values in little-endian byte order.
//
// Parameters:
// - raw: A base64-encoded string containing zlib-compressed binary data.
//
// Returns:
// - A slice of float64 values decoded from the input string.
// - An error if any step of the decoding or decompression process fails.
func ZintToFloat(raw string) ([]float64, error) {
	compressedData, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(compressedData)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}

	var decompsressedData bytes.Buffer
	if _, err := decompsressedData.ReadFrom(r); err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}

	array := make([]float64, decompsressedData.Len()/8)
	err = binary.Read(&decompsressedData, binary.LittleEndian, array)
	if err != nil {
		return nil, err
	}

	return array, nil
}
