package timeconversion_test

import (
	"testing"

	"github.com/Daniel-C-R/t8-client-go/internal/timeconversion"
)

func TestIsoStringToTimestamp(t *testing.T) {
	test := []struct {
		name     string
		input    string
		expected int64
		mustFail bool
	}{
		{
			name:     "Valid ISO String",
			input:    "2019-04-11T18:25:54",
			expected: 1555007154,
			mustFail: false,
		},
		{
			name:     "Invalid ISO String",
			input:    "2019-04-11T18:25:54Z",
			expected: 0,
			mustFail: true,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			result, err := timeconversion.IsoStringToTimestamp(tc.input)
			if tc.mustFail {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if result != tc.expected {
					t.Errorf("Expected %d but got %d", tc.expected, result)
				}
			}
		})
	}
}
