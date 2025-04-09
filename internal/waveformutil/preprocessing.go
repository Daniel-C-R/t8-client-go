package waveformutil

import (
	"gonum.org/v1/gonum/dsp/window"
)

// ZeroPadWaveform takes a slice of float64 values representing a waveform
// and returns a new slice where the length is zero-padded to the next power
// of 2 greater than or equal to the original length. The original waveform
// values are preserved at the beginning of the new slice, and the remaining
// elements are initialized to zero.
//
// Parameters:
//   - waveform: A slice of float64 values representing the input waveform.
//
// Returns:
//   - A new slice of float64 values with a length that is a power of 2,
//     containing the original waveform values followed by zero padding.
func ZeroPadWaveform(waveform []float64) []float64 {
	// Find the next power of 2 greater than or equal to the length of the waveform
	n := len(waveform)
	paddedLength := 1
	for paddedLength < n {
		paddedLength *= 2
	}

	// Create a new slice with the padded length and copy the original waveform into it
	paddedWaveform := make([]float64, paddedLength)
	copy(paddedWaveform, waveform)

	return paddedWaveform
}

// PreprocessWaveform applies preprocessing steps to the given waveform.
// It performs the following operations:
// 1. Zero-padding: Extends the waveform with zeros to a desired length.
// 2. Hanning window: Applies a Hanning window function to smooth the waveform.
//
// Parameters:
//
//	waveform []float64 - The input waveform to be preprocessed.
//
// Returns:
//
//	[]float64 - The preprocessed waveform after applying zero-padding and the Hanning window.
func PreprocessWaveform(waveform []float64) (preprocessedWaveform []float64) {
	// Apply zero-padding to the waveform
	preprocessedWaveform = ZeroPadWaveform(waveform)
	// Apply a Hanning window to the waveform
	window.Hann(preprocessedWaveform)
	return
}
