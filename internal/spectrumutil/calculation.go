package spectrumutil

import (
	"math"
	"math/cmplx"

	"github.com/Daniel-C-R/t8-client-go/internal/waveforms"
	"gonum.org/v1/gonum/cmplxs"
	"gonum.org/v1/gonum/dsp/fourier"
)

// CalculateSpectrum computes the spectrum of a given waveform using FFT (Fast Fourier Transform)
// and filters the resulting frequencies and magnitudes within a specified range.
//
// Parameters:
//   - waveform: A waveforms.Waveform struct containing the waveform data.
//   - fmin: The minimum frequency of interest in Hz.
//   - fmax: The maximum frequency of interest in Hz.
//
// Returns:
//   - A slice of float64 containing the magnitudes of the spectrum within the specified frequency range.
//   - A slice of float64 containing the corresponding frequencies within the specified range.
func CalculateSpectrum(waveform waveforms.Waveform, fmin, fmax float64) ([]float64, []float64) {
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

	return filteredSpectrum, filteredFreqs
}
