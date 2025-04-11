package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Daniel-C-R/t8-client-go/internal/getdata"
	"github.com/Daniel-C-R/t8-client-go/internal/plotutil"
	"github.com/Daniel-C-R/t8-client-go/internal/spectra"
	"gonum.org/v1/plot/vg"
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

	urlParams := getdata.NewPmodeUrlTimeParams(
		*host,
		*machine,
		*point,
		*pmode,
		*dateTime,
		user,
		password,
	)

	// Waveform
	waveform, err := getdata.GetWaveform(urlParams)
	if err != nil {
		fmt.Println("Error getting waveform:", err)
		return
	}

	plot, err := plotutil.PlotWaveform(waveform)
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

	plot, err = plotutil.PlotSpectrum(t8_spectrum, fmin, fmax)
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
	waveform.Preprocess()

	spectrum := spectra.SpectrumFromWaveform(waveform, fmin, fmax)

	plot, err = plotutil.PlotSpectrum(spectrum, fmin, fmax)
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
