package plotutil

import (
	"github.com/Daniel-C-R/t8-client-go/internal/spectra"
	"github.com/Daniel-C-R/t8-client-go/internal/waveforms"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// PlotWaveform generates a plot of a waveform given its amplitude values and sample rate.
// It creates a line plot where the X-axis represents time in seconds and the Y-axis represents
// the amplitude of the waveform.
//
// Parameters:
//   - waveform: A waveforms.Waveform struct containing the waveform data.
//
// Returns:
//   - *plot.Plot: A pointer to the generated plot.Plot object.
//   - error: An error if the plot creation or line addition fails.
func PlotWaveform(waveform waveforms.Waveform) (*plot.Plot, error) {
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

// PlotSpectrum creates a plot of a spectrum given its magnitude values and corresponding
// frequency values. The function takes a Spectrum struct and a frequency range (fmin and fmax)
// as input, and returns a pointer to a plot.Plot object or an error if the plot cannot be created.
//
// Parameters:
//   - spectrum: A Spectrum struct containing the magnitude and frequency values.
//   - fmin: The minimum frequency value to be displayed on the plot (not currently used).
//   - fmax: The maximum frequency value to be displayed on the plot (not currently used).
//
// Returns:
//   - *plot.Plot: A pointer to the created plot.Plot object.
//   - error: An error if the plot creation fails.
func PlotSpectrum(spectrum spectra.Spectrum, fmin, fmax float64) (*plot.Plot, error) {
	p := plot.New()

	pts := make(plotter.XYs, len(spectrum.Magnitudes))
	for i := range spectrum.Magnitudes {
		pts[i].X = spectrum.Frequencies[i]
		pts[i].Y = spectrum.Magnitudes[i]
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
