package datafetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Daniel-C-R/t8-client-go/internal/decoder"
	"github.com/Daniel-C-R/t8-client-go/internal/timeconversion"
	"github.com/Daniel-C-R/t8-client-go/pkg/spectra"
	"github.com/Daniel-C-R/t8-client-go/pkg/waveforms"
	"gonum.org/v1/gonum/floats"
)

type HttpDataFetcher struct {
	Host     string
	User     string
	Password string
}

type WaveformResponse struct {
	RawWaveform string  `json:"data"`
	Factor      float64 `json:"factor"`
	SampleRate  float64 `json:"sample_rate"`
}

// Method to obtain the waveform from and HTTP server.
//
// Retivese the waveform data from a remote server using an HTTP GET request.
// The request is authenticated using basic authentication with a username and password.
// The URL is constructed using the provided parameters, including the machine identifier,
// point identifier, processing mode, and timestamp.
//
// Parameters:
//   - waveformIdentifier: A PmodeTimeIdentifier structh with the machine, point, processing mode,
//     and time in ISO format.
//
// Returns:
//   - waveforms.Waveform: A struct containing the decoded waveform data.
//   - error: An error if the request fails, the response cannot be decoded, or any other issue occurs.
func (h HttpDataFetcher) GetWaveform(
	waveformIdentifier PmodeTimeIdentifier,
) (waveforms.Waveform, error) {
	timestamp, err := timeconversion.IsoStringToTimestamp(waveformIdentifier.DateTime)
	if err != nil {
		return waveforms.Waveform{}, fmt.Errorf("error parsing timestamp: %w", err)
	}

	url := fmt.Sprintf(
		"%s/waves/%s/%s/%s/%d",
		h.Host,
		waveformIdentifier.Machine,
		waveformIdentifier.Point,
		waveformIdentifier.Pmode,
		timestamp,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return waveforms.Waveform{}, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(h.User, h.Password)

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

// Method to obtain the spectrum from and HTTP server.
//
// Retivese the spectrum data from a remote server using an HTTP GET request.
// The request is authenticated using basic authentication with a username and password.
// The URL is constructed using the provided parameters, including the machine identifier,
// point identifier, processing mode, and timestamp.
//
// Parameters:
//   - spectrumIdentifier: A PmodeTimeIdentifier structh with the machine, point, processing mode,
//     and time in ISO format.
//
// Returns:
//   - spectra.Spectrum: A struct containing the decoded spectrum data.
//   - error: An error if the request fails, the response cannot be decoded, or any other issue
// 	   occurs.
//   - fmin: The minimum frequency of the spectrum.
//   - fmax: The maximum frequency of the spectrum.

func (h HttpDataFetcher) GetSpectrum(
	spectrumIdentifier PmodeTimeIdentifier,
) (spectra.Spectrum, float64, float64, error) {
	timestamp, err := timeconversion.IsoStringToTimestamp(spectrumIdentifier.DateTime)
	if err != nil {
		return spectra.Spectrum{}, 0, 0, err
	}

	url := fmt.Sprintf(
		"%s/spectra/%s/%s/%s/%d",
		h.Host,
		spectrumIdentifier.Machine,
		spectrumIdentifier.Point,
		spectrumIdentifier.Pmode,
		timestamp,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return spectra.Spectrum{}, 0, 0, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(h.User, h.Password)

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
