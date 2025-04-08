package spectrumutil

import (
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/cmplxs"
	"gonum.org/v1/gonum/dsp/fourier"
)

func CalculateSpectrum(waveform []float64, sampleRate, fmin, fmax float64) ([]float64, []float64) {
	// Perform FFT on the waveform
	fft := fourier.NewFFT(len(waveform))
	spectrum := fft.Coefficients(nil, waveform)
	cmplxs.Scale(complex(math.Sqrt(2), 0), spectrum)

	// Calculate magnitudes and frequencies
	magnitudes := make([]float64, len(spectrum))
	frequencies := make([]float64, len(spectrum))
	for i, c := range spectrum {
		magnitudes[i] = cmplx.Abs(c) / float64(len(spectrum))
		frequencies[i] = float64(i) * sampleRate / float64(len(waveform))
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
