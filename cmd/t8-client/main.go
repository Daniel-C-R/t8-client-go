package main

import (
	"fmt"
	"os"

	"github.com/Daniel-C-R/t8-client-go/internal/getdata"
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

	rawWaveform, sampleRate, err := getdata.GetWaveform(urlParams)
	if err != nil {
		fmt.Println("Error getting waveform:", err)
		return
	}
	fmt.Println("Waveform:", rawWaveform)
	fmt.Println("Sample Rate:", sampleRate)
}
