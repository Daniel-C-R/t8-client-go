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
