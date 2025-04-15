package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Daniel-C-R/t8-client-go/pkg/datafetcher"
	"github.com/Daniel-C-R/t8-client-go/pkg/spectra"
	"gonum.org/v1/plot/vg"
)

const (
	outputDir        = "output"
	waveformPlotPath = outputDir + "/waveform.png"
	spectrumPlotPath = outputDir + "/spectrum.png"
	fftSpectrumPath  = outputDir + "/fft_spectrum.png"
)

func main() {
	host := flag.String("host", "", "Host URL")
	machine := flag.String("machine", "", "Machine name")
	point := flag.String("point", "", "Point name")
	pmode := flag.String("pmode", "", "Pmode value")
	dateTime := flag.String("datetime", "", "Date and time")
	flag.Parse()

	if *host == "" || *machine == "" || *point == "" || *pmode == "" || *dateTime == "" {
		fmt.Println("All parameters (host, machine, point, pmode, datetime) are required.")
		flag.Usage()
		return
	}

	user := os.Getenv("T8_CLIENT_USER")
	password := os.Getenv("T8_CLIENT_PASSWORD")

	urlParams := datafetcher.NewPmodeUrlTimeParams(
		*host,
		*machine,
		*point,
		*pmode,
		*dateTime,
		user,
		password,
	)

	// Updated to use the HttpDataFetcher implementation
	fetcher := datafetcher.HttpDataFetcher{}

	// Waveform
	waveform, err := fetcher.GetWaveform(urlParams)
	if err != nil {
		fmt.Println("Error getting waveform:", err)
		return
	}

	plot, err := waveform.Plot()
	if err != nil {
		fmt.Println("Error plotting waveform:", err)
		return
	}
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	err = plot.Save(8*vg.Inch, 4*vg.Inch, waveformPlotPath)
	if err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

	fmt.Println("Waveform plot saved to", waveformPlotPath)

	// T8 Spectrum
	t8_spectrum, fmin, fmax, err := fetcher.GetSpectrum(urlParams)
	if err != nil {
		fmt.Println("Error getting T8 spectrum:", err)
		return
	}

	plot, err = t8_spectrum.Plot(fmin, fmax)
	if err != nil {
		fmt.Println("Error plotting T8 spectrum:", err)
		return
	}
	err = plot.Save(8*vg.Inch, 4*vg.Inch, spectrumPlotPath)
	if err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}

	fmt.Println("T8 spectrum plot saved to", spectrumPlotPath)

	// FFT Spectrum
	waveform.Preprocess()

	spectrum := spectra.SpectrumFromWaveform(waveform, fmin, fmax)

	plot, err = spectrum.Plot(fmin, fmax)
	if err != nil {
		fmt.Println("Error plotting FFT spectrum:", err)
		return
	}
	err = plot.Save(8*vg.Inch, 4*vg.Inch, fftSpectrumPath)
	if err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}
	fmt.Println("FFT spectrum plot saved to", fftSpectrumPath)
}
