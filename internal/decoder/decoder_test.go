package decoder_test

import (
	"testing"

	"github.com/Daniel-C-R/t8-client-go/internal/decoder"
)

func TestZintToFloat(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []float64
		mustFail bool
	}{
		{
			name:     "Valid Base64 String",
			input:    "eJxjZPj//389QwMAEP4D/g==",
			expected: []float64{1, -1, 32767, -32768},
		},
		{
			name:     "Invalid Base64 String",
			input:    "invalid_base64",
			expected: nil,
			mustFail: true,
		},
		{
			name:     "Empty Base64 String",
			input:    "",
			expected: nil,
			mustFail: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := decoder.ZintToFloat(tc.input)
			if tc.mustFail {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if len(result) != len(tc.expected) {
					t.Errorf("Expected length %d but got %d", len(tc.expected), len(result))
				}
				for i, v := range result {
					if v != tc.expected[i] {
						t.Errorf("Expected %f but got %f at index %d", tc.expected[i], v, i)
					}
				}
			}
		})
	}
}
