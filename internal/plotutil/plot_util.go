package plotutil

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

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

func PlotSpectrum(spectrum []float64, fmin, fmax float64) (*plot.Plot, error) {
	p := plot.New()

	pts := make(plotter.XYs, len(spectrum))
	for i := range spectrum {
		pts[i].X = fmin + float64(i)*(fmax-fmin)/float64(len(spectrum)-1)
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
