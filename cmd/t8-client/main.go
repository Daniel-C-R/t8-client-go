package main

import (
	"fmt"
	"os"

	"github.com/Daniel-C-R/t8-client-go/internal/getdata"
	"github.com/Daniel-C-R/t8-client-go/internal/plotutil"
	"github.com/Daniel-C-R/t8-client-go/internal/spectrumutil"
	"github.com/Daniel-C-R/t8-client-go/internal/waveformutil"
	"gonum.org/v1/plot/vg"
)

const (
	host     = "https://lzfs45.mirror.twave.io/lzfs45/rest"
	machine  = "LP_Turbine"
	point    = "MAD31CY005"
	pmode    = "AM1"
	dateTime = "2019-04-11T18:25:54"
)

func main() {
	user := os.Getenv("T8_CLIENT_USER")
	password := os.Getenv("T8_CLIENT_PASSWORD")

	urlParams := getdata.NewPmodeUrlTimeParams(
		host,
		machine,
		point,
		pmode,
		dateTime,
		user,
		password,
	)

	// Waveform
	waveform, sampleRate, err := getdata.GetWaveform(urlParams)
	if err != nil {
		fmt.Println("Error getting waveform:", err)
		return
	}

	plot, err := plotutil.PlotWaveform(waveform, sampleRate)
	if err != nil {
		fmt.Println("Error plotting waveform:", err)
		return
	}
	err = os.MkdirAll("output", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	err = plot.Save(8*vg.Inch, 4*vg.Inch, "output/waveform.png")
	if err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

	fmt.Println("Waveform plot saved to output/waveform.png")

	// T8 Spectrum
	t8_spectrum, fmin, fmax, err := getdata.GetSpectrum(urlParams)
	if err != nil {
		fmt.Println("Error getting T8 spectrum:", err)
		return
	}

	t8_freqs := make([]float64, len(t8_spectrum))
	step := (fmax - fmin) / float64(len(t8_spectrum)-1)
	for i := range t8_freqs {
		t8_freqs[i] = fmin + step*float64(i)
	}

	plot, err = plotutil.PlotSpectrum(t8_spectrum, t8_freqs, fmin, fmax)
	if err != nil {
		fmt.Println("Error plotting T8 spectrum:", err)
		return
	}
	err = plot.Save(8*vg.Inch, 4*vg.Inch, "output/spectrum.png")
	if err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

	fmt.Println("T8 spectrum plot saved to output/spectrum.png")

	// FFT Spectrum
	preprocessedWaveform := waveformutil.PreprocessWaveform(waveform)

	spectrum, freqs := spectrumutil.CalculateSpectrum(preprocessedWaveform, sampleRate, fmin, fmax)

	plot, err = plotutil.PlotSpectrum(spectrum, freqs, fmin, fmax)
	if err != nil {
		fmt.Println("Error plotting FFT spectrum:", err)
		return
	}
	err = plot.Save(8*vg.Inch, 4*vg.Inch, "output/fft_spectrum.png")
	if err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}
	fmt.Println("FFT spectrum plot saved to output/fft_spectrum.png")
}
