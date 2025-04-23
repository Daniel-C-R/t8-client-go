package commands

import (
	"fmt"
	"os"

	"github.com/Daniel-C-R/t8-client-go/pkg/datafetcher"
	"github.com/Daniel-C-R/t8-client-go/pkg/spectra"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot/vg"
)

const (
	outputDir        = "output"
	waveformPlotPath = outputDir + "/waveform.png"
	spectrumPlotPath = outputDir + "/spectrum.png"
	fftSpectrumPath  = outputDir + "/fft_spectrum.png"
)

func AddSpectraComparisonCommand(rootCmd *cobra.Command, baseUrlParams *datafetcher.BaseUrlParams) {
	spectraComparisonCmd := &cobra.Command{
		Use:   "spectra-comparison",
		Short: "Compare gonum and T8 spectrum calculation",
		Long: `Compares the spectrum calculated by the T8 device with the one 
that would be calculated with the gonum package. First, it fetches the 
waveform and the spectrum data of a signal from the T8 device given a 
machine, a point, a processing mode, and a datetime. Then, it calculates 
the spectrum of the waveform using the gonum package and compares it 
with the spectrum fetched from the T8 device, saving some plots in a 
directory called output.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Create data fetcher
			waveformIdentifier := datafetcher.NewPmodeTimeIdentifier(
				cmd.Flag("machine").Value.String(),
				cmd.Flag("point").Value.String(),
				cmd.Flag("pmode").Value.String(),
				cmd.Flag("datetime").Value.String(),
			)

			dataFetcher := datafetcher.HttpDataFetcher{
				Host:     baseUrlParams.Host,
				User:     baseUrlParams.User,
				Password: baseUrlParams.Password,
			}

			// Waveform
			waveform, err := dataFetcher.GetWaveform(waveformIdentifier)
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
			t8_spectrum, fmin, fmax, err := dataFetcher.GetSpectrum(waveformIdentifier)
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
		},
	}

	var machine, point, processingMode, datetime string

	spectraComparisonCmd.Flags().StringVar(&machine, "machine", "", "Machine identifier (required)")
	spectraComparisonCmd.Flags().StringVar(&point, "point", "", "Point identifier (required)")
	spectraComparisonCmd.Flags().
		StringVar(&processingMode, "pmode", "", "Processing mode (required)")
	spectraComparisonCmd.Flags().
		StringVar(&datetime, "datetime", "", "Datetime for the data (required)")

	if err := spectraComparisonCmd.MarkFlagRequired("machine"); err != nil {
		fmt.Println("Error marking 'machine' flag as required:", err)
		return
	}
	if err := spectraComparisonCmd.MarkFlagRequired("point"); err != nil {
		fmt.Println("Error marking 'point' flag as required:", err)
		return
	}
	if err := spectraComparisonCmd.MarkFlagRequired("pmode"); err != nil {
		fmt.Println("Error marking 'pmode' flag as required:", err)
		return
	}
	if err := spectraComparisonCmd.MarkFlagRequired("datetime"); err != nil {
		fmt.Println("Error marking 'datetime' flag as required:", err)
		return
	}

	rootCmd.AddCommand(spectraComparisonCmd)
}
