package main

import (
	"fmt"
	"os"

	"github.com/Daniel-C-R/t8-client-go/internal/getdata"
	"github.com/Daniel-C-R/t8-client-go/internal/plotutil"
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
	spectrum, fmin, fmax, err := getdata.GetSpectrum(urlParams)
	if err != nil {
		fmt.Println("Error getting T8 spectrum:", err)
		return
	}

	plot, err = plotutil.PlotSpectrum(spectrum, fmin, fmax)
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
}
