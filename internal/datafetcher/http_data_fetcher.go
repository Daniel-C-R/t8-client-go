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

type HttpDataFetcher struct{}

type WaveformResponse struct {
	RawWaveform string  `json:"data"`
	Factor      float64 `json:"factor"`
	SampleRate  float64 `json:"sample_rate"`
}

// GetWaveform retrieves waveform data from a remote server.
//
// Parameters:
//   - urlParams: A PmodeUrlTimeParams struct containing the following fields:
//   - Host: The base URL of the server.
//   - Machine: The machine identifier.
//   - Point: The measurement point identifier.
//   - Pmode: The processing mode.
//   - DateTime: The timestamp for the data request in ISO format.
//   - User: The username for authentication.
//   - Password: The password for authentication.
//
// Returns:
//   - waveforms.Waveform: A struct containing the decoded waveform data, including samples and sample rate.
//   - error: An error if the request fails, the response cannot be decoded, or any other issue occurs.
func (h HttpDataFetcher) GetWaveform(urlParams PmodeUrlTimeParams) (waveforms.Waveform, error) {
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

// GetSpectrum retrieves spectrum data from a remote server.
//
// Parameters:
//   - urlParams: A PmodeUrlTimeParams struct containing the following fields:
//   - Host: The base URL of the server.
//   - Machine: The machine identifier.
//   - Point: The measurement point identifier.
//   - Pmode: The processing mode.
//   - DateTime: The timestamp for the data request in ISO format.
//   - User: The username for authentication.
//   - Password: The password for authentication.
//
// Returns:
//   - spectra.Spectrum: A struct containing the decoded spectrum data, including frequencies and magnitudes.
//   - float64: The minimum frequency of the spectrum.
//   - float64: The maximum frequency of the spectrum.
//   - error: An error if the request fails, the response cannot be decoded, or any other issue occurs.
func (h HttpDataFetcher) GetSpectrum(
	urlParams PmodeUrlTimeParams,
) (spectra.Spectrum, float64, float64, error) {
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
