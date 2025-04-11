package spectra

import (
	"math"
	"math/cmplx"

	"github.com/Daniel-C-R/t8-client-go/internal/waveforms"
	"gonum.org/v1/gonum/cmplxs"
	"gonum.org/v1/gonum/dsp/fourier"
)

type Spectrum struct {
	Magnitudes  []float64
	Frequencies []float64
}

// NewSpectrum creates a new Spectrum object with the given magnitudes and frequency range.
// It calculates the frequencies based on the provided minimum frequency (fmin).
//
// Parameters:
//   - magnitudes: A slice of float64 representing the magnitudes of the spectrum.
//   - fmin: The minimum frequency value (float64).
//   - fmax: The maximum frequency value (float64).
//
// Returns:
//
//	A Spectrum object containing the provided magnitudes and the calculated frequencies.
func NewSpectrum(magnitudes []float64, fmin, fmax float64) Spectrum {
	// Calculate frequencies
	frequencies := make([]float64, len(magnitudes))
	step := (fmax - fmin) / float64(len(magnitudes)-1)
	for i := range magnitudes {
		frequencies[i] = fmin + float64(i)*step
	}

	// Return spectrum object
	return Spectrum{
		Magnitudes:  magnitudes,
		Frequencies: frequencies,
	}
}

// SpectrumFromWaveform computes the spectrum of a given waveform using FFT (Fast Fourier Transform)
// and filters the resulting frequencies and magnitudes within a specified range.
//
// Parameters:
//   - waveform: A waveforms.Waveform struct containing the waveform data.
//   - fmin: The minimum frequency of interest in Hz.
//   - fmax: The maximum frequency of interest in Hz.
//
// Returns:
//   - A Spectrum struct containing the magnitudes and corresponding frequencies within the specified range.
func SpectrumFromWaveform(waveform waveforms.Waveform, fmin, fmax float64) Spectrum {
	// Perform FFT on the waveform
	fft := fourier.NewFFT(len(waveform.Samples))
	spectrum := fft.Coefficients(nil, waveform.Samples)
	cmplxs.Scale(complex(math.Sqrt(2), 0), spectrum)

	// Calculate magnitudes and frequencies
	magnitudes := make([]float64, len(spectrum))
	frequencies := make([]float64, len(spectrum))
	for i, c := range spectrum {
		magnitudes[i] = cmplx.Abs(c) / float64(len(spectrum))
		frequencies[i] = float64(i) * waveform.SampleRate / float64(len(waveform.Samples))
	}

	// Filter the spectrum and frequencies within the specified range
	var filteredSpectrum []float64
	var filteredFreqs []float64
	for i, freq := range frequencies {
		if freq >= fmin && freq <= fmax {
			filteredSpectrum = append(filteredSpectrum, magnitudes[i])
			filteredFreqs = append(filteredFreqs, freq)
		}
	}

	return Spectrum{
		Frequencies: filteredFreqs,
		Magnitudes:  filteredSpectrum,
	}
}
