package getdata_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Daniel-C-R/t8-client-go/internal/getdata"
	"gonum.org/v1/gonum/floats"
)

// TestGetWaveformSuccess tests the successful retrieval of a waveform.
func TestGetWaveformSuccess(t *testing.T) {
	mockWaveformResponse := getdata.WaveformResponse{
		RawWaveform: "eJxjZPj//389QwMAEP4D/g==",
		Factor:      2.0,
		SampleRate:  2560,
	}

	expectedWaveform := []float64{1, -1, 3.2767e04, -3.2768e04}
	floats.Scale(mockWaveformResponse.Factor, expectedWaveform)

	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(mockWaveformResponse); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, sampleRate, err := getdata.GetWaveform(mockPmodeTimeParams)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(waveform, expectedWaveform) {
		t.Errorf("expected waveform %v, got %v", expectedWaveform, waveform)
	}

	if sampleRate != mockWaveformResponse.SampleRate {
		t.Errorf("expected sample rate %v, got %v", mockWaveformResponse.SampleRate, sampleRate)
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, sampleRate, err := getdata.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if waveform != nil {
		t.Errorf("expected nil waveform, got %v", waveform)
	}

	if sampleRate != 0 {
		t.Errorf("expected sample rate 0, got %v", sampleRate)
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, sampleRate, err := getdata.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if waveform != nil {
		t.Errorf("expected nil waveform, got %v", waveform)
	}

	if sampleRate != 0 {
		t.Errorf("expected sample rate 0, got %v", sampleRate)
	}
}

// TestGetWaveformInvalidTimestamp tests the case when an invalid timestamp is provided.
func TestGetWaveformInvalidTimestamp(t *testing.T) {
	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(getdata.WaveformResponse{}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"invalid_timestamp",
		"user",
		"password",
	)

	waveform, sampleRate, err := getdata.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if waveform != nil {
		t.Errorf("expected nil waveform, got %v", waveform)
	}

	if sampleRate != 0 {
		t.Errorf("expected sample rate 0, got %v", sampleRate)
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, sampleRate, err := getdata.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if waveform != nil {
		t.Errorf("expected nil waveform, got %v", waveform)
	}

	if sampleRate != 0 {
		t.Errorf("expected sample rate 0, got %v", sampleRate)
	}
}

// TestGetWaveformInvalidWaveformData tests the case when the server returns invalid waveform data.
func TestGetWaveformInvalidWaveformData(t *testing.T) {
	mockWaveformResponse := getdata.WaveformResponse{
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	waveform, sampleRate, err := getdata.GetWaveform(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if waveform != nil {
		t.Errorf("expected nil waveform, got %v", waveform)
	}

	if sampleRate != 0 {
		t.Errorf("expected sample rate 0, got %v", sampleRate)
	}
}

// TestGetSpectrumSuccess tests the successful retrieval of a spectrum.
func TestGetSpectrumSuccess(t *testing.T) {
	mockSpectrumResponse := getdata.SpectrumResponse{
		RawSpectrum: "eJxjZPj//389QwMAEP4D/g==",
		Factor:      2.0,
		Fmin:        0.0,
		Fmax:        1280.0,
	}

	expectedSpectrum := []float64{1, -1, 3.2767e04, -3.2768e04}
	floats.Scale(mockSpectrumResponse.Factor, expectedSpectrum)

	mock_server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(mockSpectrumResponse); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := getdata.GetSpectrum(mockPmodeTimeParams)
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := getdata.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if spectrum != nil {
		t.Errorf("expected nil spectrum, got %v", spectrum)
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := getdata.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if spectrum != nil {
		t.Errorf("expected nil spectrum, got %v", spectrum)
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
			if err := json.NewEncoder(w).Encode(getdata.SpectrumResponse{}); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}),
	)
	defer mock_server.Close()

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"invalid_timestamp",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := getdata.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if spectrum != nil {
		t.Errorf("expected nil spectrum, got %v", spectrum)
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := getdata.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if spectrum != nil {
		t.Errorf("expected nil spectrum, got %v", spectrum)
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
	mockSpectrumResponse := getdata.SpectrumResponse{
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

	mockPmodeTimeParams := getdata.NewPmodeUrlTimeParams(
		mock_server.URL,
		"test_machine",
		"test_point",
		"test_pmode",
		"2019-04-10T14:48:44",
		"user",
		"password",
	)

	spectrum, fmin, fmax, err := getdata.GetSpectrum(mockPmodeTimeParams)

	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if spectrum != nil {
		t.Errorf("expected nil spectrum, got %v", spectrum)
	}

	if fmin != 0 {
		t.Errorf("expected fmin 0, got %v", fmin)
	}

	if fmax != 0 {
		t.Errorf("expected fmax 0, got %v", fmax)
	}
}
