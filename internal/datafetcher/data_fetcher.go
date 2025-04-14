package datafetcher

import (
	"github.com/Daniel-C-R/t8-client-go/internal/spectra"
	"github.com/Daniel-C-R/t8-client-go/internal/waveforms"
)

type DataFetcher interface {
	// GetWaveform retrieves waveform data from a remote server.
	//
	// Parameters:
	//   - urlParams: A PmodeUrlTimeParams struct containing the following fields:
	//       - Host: The base URL of the server.
	//       - Machine: The machine identifier.
	//       - Point: The measurement point identifier.
	//       - Pmode: The processing mode.
	//       - DateTime: The timestamp for the data request in ISO format.
	//       - User: The username for authentication.
	//       - Password: The password for authentication.
	//
	// Returns:
	//   - waveforms.Waveform: A struct containing the decoded waveform data, including samples and sample rate.
	//   - error: An error if the request fails, the response cannot be decoded, or any other issue occurs.
	GetWaveform(urlParams PmodeUrlTimeParams) (waveforms.Waveform, error)

	// GetSpectrum retrieves spectrum data from a remote server.
	//
	// Parameters:
	//   - urlParams: A PmodeUrlTimeParams struct containing the following fields:
	//       - Host: The base URL of the server.
	//       - Machine: The machine identifier.
	//       - Point: The measurement point identifier.
	//       - Pmode: The processing mode.
	//       - DateTime: The timestamp for the data request in ISO format.
	//       - User: The username for authentication.
	//       - Password: The password for authentication.
	//
	// Returns:
	//   - spectra.Spectrum: A struct containing the decoded spectrum data, including frequencies and magnitudes.
	//   - error: An error if the request fails, the response cannot be decoded, or any other issue occurs.
	GetSpectrum(urlParams PmodeUrlTimeParams) (spectra.Spectrum, error)
}
