package waveforms

import (
	"gonum.org/v1/gonum/dsp/window"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

type Waveform struct {
	Samples    []float64
	SampleRate float64
}

// ZeroPadding adjusts the waveform's sample slice to have a length that is
// the next power of 2 greater than or equal to its current length. This is
// achieved by creating a new slice with the padded length and copying the
// original samples into it, leaving the additional elements initialized to zero.
func (waveform Waveform) ZeroPadding() {
	// Find the next power of 2 greater than or equal to the length of the waveform
	n := len(waveform.Samples)
	paddedLength := 1
	for paddedLength < n {
		paddedLength *= 2
	}

	// Create a new slice with the padded length and copy the original waveform into it
	paddedWaveform := make([]float64, paddedLength)
	copy(paddedWaveform, waveform.Samples)
}

// Preprocess applies preprocessing steps to the waveform.
// It performs zero-padding to extend the waveform and applies
// a Hanning window function to the waveform samples to reduce
// spectral leakage during analysis.
func (waveform Waveform) Preprocess() {
	// Apply zero-padding to the waveform
	waveform.ZeroPadding()
	// Apply a Hanning window to the waveform
	window.Hann(waveform.Samples)
}

// Plot genera una gráfica de la forma de onda actual.
func (waveform Waveform) Plot() (*plot.Plot, error) {
	p := plot.New()

	pts := make(plotter.XYs, len(waveform.Samples))
	for i := range waveform.Samples {
		pts[i].X = float64(i) / waveform.SampleRate
		pts[i].Y = waveform.Samples[i]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return nil, err
	}

	p.Add(line)
	p.Title.Text = "Waveform"
	p.X.Label.Text = "Time (s)"
	p.Y.Label.Text = "Amplitude"

	return p, nil
}
