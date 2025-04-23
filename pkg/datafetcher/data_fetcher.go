package datafetcher

import (
	"github.com/Daniel-C-R/t8-client-go/pkg/spectra"
	"github.com/Daniel-C-R/t8-client-go/pkg/waveforms"
)

type DataFetcher interface {
	// Method to obtain the waveform given its identifier
	//
	// Retrieves a waveform record based on its identifier, which is determined by the machine, point, processing mode,
	// and registration date.
	//
	// Parameters:
	//   - waveformIdentifier: The identifier of the waveform, which includes the machine, point, and processing mode.
	//
	// Returns:
	//   - A waveform record containing the samples and sample rate.
	//   - An error if the retrieval fails.
	GetWaveform(waveformIdentifier PmodeTimeIdentifier) (waveforms.Waveform, error)

	// Method to obtain the spectrum given its time-based identifier
	//
	// Retrieves a spectrum record along with its frequency range based given its machine, point, processing mode, and
	// registration date.
	//
	// Parameters:
	//   - spectrumIdentifier: The identifier of the spectrum, which includes the machine, point, and processing mode.
	//
	// Returns:
	//   - A spectrum record containing the magnitudes and frequencies.
	//   - The minimum frequency of the spectrum.
	//   - The maximum frequency of the spectrum.
	//   - An error if the retrieval fails.
	GetSpectrum(spectrumIdentifier PmodeTimeIdentifier) (spectra.Spectrum, float64, float64, error)
}
