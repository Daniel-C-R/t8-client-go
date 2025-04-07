package getdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Daniel-C-R/t8-client-go/internal/decoder"
	"github.com/Daniel-C-R/t8-client-go/internal/timeutil"
	"gonum.org/v1/gonum/floats"
)

type WaveformResponse struct {
	RawWaveform string  `json:"data"`
	Factor      float64 `json:"factor"`
	SampleRate  float64 `json:"sample_rate"`
}

// GetWaveform retrieves waveform data from a specified host and endpoint.
//
// Parameters:
//   - urlParams: A PmodeUrlTimeParams struct containing the host, machine, point, pmode,
//     dateTime, user, and password for the request.
//
// Returns:
//   - A slice of float64 representing the decoded waveform data.
//   - A float64 representing the sample rate of the waveform.
//   - An error if the request fails or the response cannot be decoded.
func GetWaveform(urlParams PmodeUrlTimeParams) ([]float64, float64, error) {
	timestamp, err := timeutil.IsoStringToTimestamp(urlParams.DateTime)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing timestamp: %w", err)
	}

	url := fmt.Sprintf(
		"%s/waves/%s/%s/%s/%d",
		urlParams.Host,
		urlParams.Machine,
		urlParams.Point,
		urlParams.Pmode,
		timestamp,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(urlParams.User, urlParams.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading response body: %w", err)
	}

	var waveformResponse WaveformResponse
	if err := json.Unmarshal(body, &waveformResponse); err != nil {
		return nil, 0, fmt.Errorf("error decoding JSON response: %w", err)
	}

	waveform, err := decoder.ZintToFloat(waveformResponse.RawWaveform)
	if err != nil {
		return nil, 0, fmt.Errorf("error decoding waveform data: %w", err)
	}

	floats.Scale(waveformResponse.Factor, waveform)

	return waveform, waveformResponse.SampleRate, nil
}

type SpectrumResponse struct {
	RawSpectrum string  `json:"data"`
	Factor      float64 `json:"factor"`
	Fmin        float64 `json:"min_freq"`
	Fmax        float64 `json:"max_freq"`
}

func GetSpectrum(urlParams PmodeUrlTimeParams) ([]float64, float64, float64, error) {
	timestamp, err := timeutil.IsoStringToTimestamp(urlParams.DateTime)
	if err != nil {
		return nil, 0, 0, err
	}

	url := fmt.Sprintf(
		"%s/spectra/%s/%s/%s/%d",
		urlParams.Host,
		urlParams.Machine,
		urlParams.Point,
		urlParams.Pmode,
		timestamp,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(urlParams.User, urlParams.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error reading response body: %w", err)
	}

	var spectrumResponse SpectrumResponse
	if err := json.Unmarshal(body, &spectrumResponse); err != nil {
		return nil, 0, 0, fmt.Errorf("error decoding JSON response: %w", err)
	}

	spectrum, err := decoder.ZintToFloat(spectrumResponse.RawSpectrum)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("error decoding spectrum data: %w", err)
	}

	floats.Scale(spectrumResponse.Factor, spectrum)

	return spectrum, spectrumResponse.Fmin, spectrumResponse.Fmax, nil
}
