package main

import (
	"os"

	"github.com/Daniel-C-R/t8-client-go/cmd/t8-client/commands"
	"github.com/Daniel-C-R/t8-client-go/pkg/datafetcher"
	"github.com/spf13/cobra"
)

var baseUrlParams datafetcher.BaseUrlParams

var rootCmd = &cobra.Command{
	Use:   "t8-client",
	Short: "T8 Client CLI",
	Long:  "Client CLI for interacting with T8 devices to retrieve and analyze waveform and spectrum data",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	baseUrlParams.Host = os.Getenv("T8_CLIENT_HOST")
	baseUrlParams.User = os.Getenv("T8_CLIENT_USER")
	baseUrlParams.Password = os.Getenv("T8_CLIENT_PASSWORD")

	commands.AddSpectraComparisonCommand(rootCmd, &baseUrlParams)
}
