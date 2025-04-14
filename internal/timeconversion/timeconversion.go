package timeconversion

import "time"

// IsoStringToTimestamp converts an ISO 8601 formatted string into a Unix timestamp.
// It takes a string in the RFC3339 format as input and returns the corresponding
// Unix timestamp as an int64. If the input string is not in a valid RFC3339 format,
// an error is returned.
//
// Parameters:
//   - isoString: A string representing the date and time in ISO 8601 (RFC3339) format.
//
// Returns:
//   - int64: The Unix timestamp corresponding to the input date and time.
//   - error: An error if the input string is not in a valid RFC3339 format.
func IsoStringToTimestamp(isoString string) (int64, error) {
	t, err := time.Parse("2006-01-02T15:04:05", isoString)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
