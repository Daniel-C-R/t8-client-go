package plotutil

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// PlotWaveform generates a plot of a waveform given its amplitude values and sample rate.
// It creates a line plot where the X-axis represents time in seconds and the Y-axis represents
// the amplitude of the waveform.
//
// Parameters:
//   - waveform: A slice of float64 values representing the amplitude of the waveform.
//   - sampleRate: A float64 value representing the number of samples per second.
//
// Returns:
//   - *plot.Plot: A pointer to the generated plot.Plot object.
//   - error: An error if the plot creation or line addition fails.
func PlotWaveform(waveform []float64, sampleRate float64) (*plot.Plot, error) {
	p := plot.New()

	pts := make(plotter.XYs, len(waveform))
	for i := range waveform {
		pts[i].X = float64(i) / sampleRate
		pts[i].Y = waveform[i]
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

// PlotSpectrum creates a plot of a spectrum given its magnitude values and corresponding
// frequency values. The function takes the spectrum data, frequency data, and a frequency
// range (fmin and fmax) as input, and returns a pointer to a plot.Plot object or an error
// if the plot cannot be created.
//
// Parameters:
//   - spectrum: A slice of float64 representing the magnitude values of the spectrum.
//   - freqs: A slice of float64 representing the corresponding frequency values.
//   - fmin: The minimum frequency value to be displayed on the plot (not currently used).
//   - fmax: The maximum frequency value to be displayed on the plot (not currently used).
//
// Returns:
//   - *plot.Plot: A pointer to the created plot.Plot object.
//   - error: An error if the plot creation fails.
func PlotSpectrum(spectrum, freqs []float64, fmin, fmax float64) (*plot.Plot, error) {
	p := plot.New()

	pts := make(plotter.XYs, len(spectrum))
	for i := range spectrum {
		pts[i].X = freqs[i]
		pts[i].Y = spectrum[i]
	}

	line, err := plotter.NewLine(pts)
	if err != nil {
		return nil, err
	}

	p.Add(line)
	p.Title.Text = "Spectrum"
	p.X.Label.Text = "Frequency (Hz)"
	p.Y.Label.Text = "Magnitude"

	return p, nil
}
