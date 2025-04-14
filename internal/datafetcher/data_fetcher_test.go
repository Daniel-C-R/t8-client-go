package datafetcher_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Daniel-C-R/t8-client-go/internal/datafetcher"
	"github.com/Daniel-C-R/t8-client-go/internal/spectra"
	"github.com/Daniel-C-R/t8-client-go/internal/waveforms"
	"gonum.org/v1/gonum/floats"
)

// TestGetWaveformSuccess tests the successful retrieval of a waveform.
func TestGetWaveformSuccess(t *testing.T) {
	mockWaveformResponse := datafetcher.WaveformResponse{
		RawWaveform: "eJxjZPj//389QwMAEP4D/g==",
		Factor:      2.0,
		SampleRate:  2560,
	}

	expectedWaveform := waveforms.Waveform{
		Samples:    []float64{1, -1, 3.2767e04, -3.2768e04},
		SampleRate: 2560,
	}
	floats.Scale(mockWaveformResponse.Factor, expectedWaveform.Samples)

	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(mockWaveformResponse); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, err := datafetcher.GetWaveform(mockPmodeTimeParams)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(waveform, expectedWaveform) {
		t.Errorf("expected waveform %v, got %v", expectedWaveform, waveform)
	}
}

// TestGetWaveformFailure tests the failure case when the server returns an error.
func TestGetWaveformFailure(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, err := datafetcher.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(waveform, waveforms.Waveform{}) {
		t.Errorf("expected empty waveform, got %v", waveform)
	}
}

// TestGetWaveformInvalidJSON tests the case when the server returns invalid JSON.
func TestGetWaveformInvalidJSON(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte("invalid json")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, err := datafetcher.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(waveform, waveforms.Waveform{}) {
		t.Errorf("expected empty waveform, got %v", waveform)
	}
}

// TestGetWaveformInvalidTimestamp tests the case when an invalid timestamp is provided.
func TestGetWaveformInvalidTimestamp(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(datafetcher.WaveformResponse{}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"invalid_timestamp",
		"user",
		"password",
	)

	waveform, err := datafetcher.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(waveform, waveforms.Waveform{}) {
		t.Errorf("expected empty waveform, got %v", waveform)
	}
}

// TestGetWaveformEmptyResponse tests the case when the server returns an empty response.
func TestGetWaveformEmptyResponse(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte("{}")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, err := datafetcher.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(waveform, waveforms.Waveform{}) {
		t.Errorf("expected empty waveform, got %v", waveform)
	}
}

// TestGetWaveformInvalidWaveformData tests the case when the server returns invalid waveform data.
func TestGetWaveformInvalidWaveformData(t *testing.T) {
	mockWaveformResponse := datafetcher.WaveformResponse{
		RawWaveform: "invalid_waveform_data",
	}

	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(mockWaveformResponse); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, err := datafetcher.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(waveform, waveforms.Waveform{}) {
		t.Errorf("expected empty waveform, got %v", waveform)
	}
}

// TestGetSpectrumSuccess tests the successful retrieval of a spectrum.
func TestGetSpectrumSuccess(t *testing.T) {
	mockSpectrumResponse := datafetcher.SpectrumResponse{
		RawSpectrum: "eJxjZPj//389QwMAEP4D/g==",
		Factor:      2.0,
		Fmin:        0.0,
		Fmax:        1280.0,
	}

	expected_magnitudes := []float64{1, -1, 3.2767e04, -3.2768e04}
	floats.Scale(mockSpectrumResponse.Factor, expected_magnitudes)
	expectedSpectrum := spectra.NewSpectrum(
		expected_magnitudes,
		mockSpectrumResponse.Fmin,
		mockSpectrumResponse.Fmax,
	)

	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(mockSpectrumResponse); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := datafetcher.GetSpectrum(mockPmodeTimeParams)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(spectrum, expectedSpectrum) {
		t.Errorf("expected spectrum %v, got %v", expectedSpectrum, spectrum)
	}

	if fmin != mockSpectrumResponse.Fmin {
		t.Errorf("expected fmin %v, got %v", mockSpectrumResponse.Fmin, fmin)
	}

	if fmax != mockSpectrumResponse.Fmax {
		t.Errorf("expected fmax %v, got %v", mockSpectrumResponse.Fmax, fmax)
	}
}

// TestGetSpectrumFailure tests the failure case when the server returns an error.
func TestGetSpectrumFailure(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := datafetcher.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(spectrum, spectra.Spectrum{}) {
		t.Errorf("expected empty spectrum, got %v", spectrum)
	}

	if fmin != 0 {
		t.Errorf("expected fmin 0, got %v", fmin)
	}

	if fmax != 0 {
		t.Errorf("expected fmax 0, got %v", fmax)
	}
}

// TestGetSpectrumInvalidJSON tests the case when the server returns invalid JSON.
func TestGetSpectrumInvalidJSON(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte("invalid json")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := datafetcher.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(spectrum, spectra.Spectrum{}) {
		t.Errorf("expected empty spectrum, got %v", spectrum)
	}

	if fmin != 0 {
		t.Errorf("expected fmin 0, got %v", fmin)
	}

	if fmax != 0 {
		t.Errorf("expected fmax 0, got %v", fmax)
	}
}

// TestGetSpectrumInvalidTimestamp tests the case when an invalid timestamp is provided.
func TestGetSpectrumInvalidTimestamp(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(datafetcher.SpectrumResponse{}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"invalid_timestamp",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := datafetcher.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(spectrum, spectra.Spectrum{}) {
		t.Errorf("expected empty spectrum, got %v", spectrum)
	}

	if fmin != 0 {
		t.Errorf("expected fmin 0, got %v", fmin)
	}

	if fmax != 0 {
		t.Errorf("expected fmax 0, got %v", fmax)
	}
}

// TestGetSpectrumEmptyResponse tests the case when the server returns an empty response.
func TestGetSpectrumEmptyResponse(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if _, err := w.Write([]byte("{}")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := datafetcher.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(spectrum, spectra.Spectrum{}) {
		t.Errorf("expected empty spectrum, got %v", spectrum)
	}

	if fmin != 0 {
		t.Errorf("expected fmin 0, got %v", fmin)
	}

	if fmax != 0 {
		t.Errorf("expected fmax 0, got %v", fmax)
	}
}

// TestGetSpectrumInvalidSpectrumData tests the case when the server returns invalid spectrum data.
func TestGetSpectrumInvalidSpectrumData(t *testing.T) {
	mockSpectrumResponse := datafetcher.SpectrumResponse{
		RawSpectrum: "invalid_spectrum_data",
	}

	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(mockSpectrumResponse); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := datafetcher.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := datafetcher.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if !reflect.DeepEqual(spectrum, spectra.Spectrum{}) {
		t.Errorf("expected empty spectrum, got %v", spectrum)
	}

	if fmin != 0 {
		t.Errorf("expected fmin 0, got %v", fmin)
	}

	if fmax != 0 {
		t.Errorf("expected fmax 0, got %v", fmax)
	}
}
