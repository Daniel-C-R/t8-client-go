package getdata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Daniel-C-R/t8-client-go/internal/decoder"
	"github.com/Daniel-C-R/t8-client-go/internal/timeutil"
	"github.com/Daniel-C-R/t8-client-go/internal/waveforms"
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
//   - waveforms.Waveform: The decoded waveform data.
//   - An error if the request fails or the response cannot be decoded.
func GetWaveform(urlParams PmodeUrlTimeParams) (waveforms.Waveform, error) {
	timestamp, err := timeutil.IsoStringToTimestamp(urlParams.DateTime)
	if err != nil {
		return waveforms.Waveform{}, fmt.Errorf("error parsing timestamp: %w", err)
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
		return waveforms.Waveform{}, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(urlParams.User, urlParams.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return waveforms.Waveform{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return waveforms.Waveform{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return waveforms.Waveform{}, fmt.Errorf("error reading response body: %w", err)
	}

	var waveformResponse WaveformResponse
	if err := json.Unmarshal(body, &waveformResponse); err != nil {
		return waveforms.Waveform{}, fmt.Errorf("error decoding JSON response: %w", err)
	}

	samples, err := decoder.ZintToFloat(waveformResponse.RawWaveform)
	if err != nil {
		return waveforms.Waveform{}, fmt.Errorf("error decoding waveform data: %w", err)
	}

	floats.Scale(waveformResponse.Factor, samples)

	waveform := waveforms.Waveform{Samples: samples, SampleRate: waveformResponse.SampleRate}

	return waveform, nil
}

type SpectrumResponse struct {
	RawSpectrum string  `json:"data"`
	Factor      float64 `json:"factor"`
	Fmin        float64 `json:"min_freq"`
	Fmax        float64 `json:"max_freq"`
}

// GetSpectrum retrieves a spectrum from a remote server based on the provided parameters.
// It sends an HTTP GET request to the specified endpoint, parses the response, and returns
// the spectrum data along with the minimum and maximum frequency values.
//
// Parameters:
//   - urlParams: A PmodeUrlTimeParams struct containing the necessary parameters for the request,
//     including host, machine, point, pmode, datetime, user, and password.
//
// Returns:
//   - []float64: The decoded spectrum data as a slice of float64 values.
//   - float64: The minimum frequency value (Fmin) of the spectrum.
//   - float64: The maximum frequency value (Fmax) of the spectrum.
//   - error: An error object if any issues occur during the process, otherwise nil.
//
// Errors:
//   - Returns an error if the datetime conversion fails.
//   - Returns an error if the HTTP request cannot be created or executed.
//   - Returns an error if the response status code is not 200 OK.
//   - Returns an error if the response body cannot be read or parsed as JSON.
//   - Returns an error if the spectrum data cannot be decoded.
//
// Example usage:
//
//	spectrum, fmin, fmax, err := GetSpectrum(urlParams)
//	if err != nil {
//	    log.Fatalf("Failed to get spectrum: %v", err)
//	}
//	fmt.Printf("Spectrum: %v, Fmin: %f, Fmax: %f\n", spectrum, fmin, fmax)
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
