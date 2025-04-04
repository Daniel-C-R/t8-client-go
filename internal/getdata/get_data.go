package getdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WaveformResponse struct {
	RawWaveform string  `json:"data"`
	Factor      float64 `json:"factor"`
	SampleRate  float64 `json:"sample_rate"`
}

// GetWaveform retrieves waveform data from a specified host and endpoint.
//
// Parameters:
//   - host: The base URL of the server hosting the waveform data.
//   - machine: The identifier of the machine to retrieve data for.
//   - point: The specific point of interest on the machine.
//   - pmode: The processing mode for the waveform data.
//   - timestamp: The timestamp for the desired waveform data.
//   - user: The username for authentication (not currently used in the function).
//   - password: The password for authentication (not currently used in the function).
//
// Returns:
//   - A string containing the raw waveform data.
//   - A float64 representing the sample rate of the waveform.
//   - An error if the request fails or the response cannot be decoded.
//
// TODO:
//   - Waveform data is encoded as a base64 zlib compressed string. It should be decoded as returned as a float64 array.
func GetWaveform(
	host string,
	machine string,
	point string,
	pmode string,
	timestamp int64,
	user string,
	password string,
) (string, float64, error) {
	url := fmt.Sprintf("%s/waves/%s/%s/%s/%d", host, machine, point, pmode, timestamp)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(user, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error reading response body: %w", err)
	}

	var waveformResponse WaveformResponse
	if err := json.Unmarshal(body, &waveformResponse); err != nil {
		return "", 0, fmt.Errorf("error decoding JSON response: %w", err)
	}

	return waveformResponse.RawWaveform, waveformResponse.SampleRate, nil
}
