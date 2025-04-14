package datafetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Daniel-C-R/t8-client-go/internal/decoder"
	"github.com/Daniel-C-R/t8-client-go/internal/spectra"
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
// the spectrum data as a Spectrum struct.
//
// Parameters:
//   - urlParams: A PmodeUrlTimeParams struct containing the necessary parameters for the request.
//
// Returns:
//   - Spectrum: The decoded spectrum data as a Spectrum struct.
//   - fmin: The minimum frequency of the spectrum.
//   - fmax: The maximum frequency of the spectrum.
//   - error: An error object if any issues occur during the process, otherwise nil.
func GetSpectrum(urlParams PmodeUrlTimeParams) (spectra.Spectrum, float64, float64, error) {
	timestamp, err := timeutil.IsoStringToTimestamp(urlParams.DateTime)
	if err != nil {
		return spectra.Spectrum{}, 0, 0, err
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
		return spectra.Spectrum{}, 0, 0, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(urlParams.User, urlParams.Password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return spectra.Spectrum{}, 0, 0, fmt.Errorf("error making GET request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return spectra.Spectrum{}, 0, 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return spectra.Spectrum{}, 0, 0, fmt.Errorf("error reading response body: %w", err)
	}

	var spectrumResponse SpectrumResponse
	if err := json.Unmarshal(body, &spectrumResponse); err != nil {
		return spectra.Spectrum{}, 0, 0, fmt.Errorf("error decoding JSON response: %w", err)
	}

	spectrum, err := decoder.ZintToFloat(spectrumResponse.RawSpectrum)
	if err != nil {
		return spectra.Spectrum{}, 0, 0, fmt.Errorf("error decoding spectrum data: %w", err)
	}

	floats.Scale(spectrumResponse.Factor, spectrum)

	frequencies := make([]float64, len(spectrum))
	step := (spectrumResponse.Fmax - spectrumResponse.Fmin) / float64(len(spectrum)-1)
	for i := range frequencies {
		frequencies[i] = spectrumResponse.Fmin + step*float64(i)
	}

	spectrumObject := spectra.Spectrum{
		Frequencies: frequencies,
		Magnitudes:  spectrum,
	}

	return spectrumObject, spectrumResponse.Fmin, spectrumResponse.Fmax, nil
}
