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

	var decompressedData bytes.Buffer
	if _, err := decompressedData.ReadFrom(r); err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}

	data := decompressedData.Bytes()
	array := make([]float64, len(data)/2)
	for i := 0; i < len(array); i++ {
		array[i] = float64(int16(binary.LittleEndian.Uint16(data[i*2 : (i+1)*2])))
	}

	return array, nil
}
