package cmd

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var log zerolog.Logger

var Version string = "1.0.0"

var rootCmd = &cobra.Command{
	Use:     "cube",
	Version: Version,
	Short:   "cube is a devcontainer managment tool",
	Long:    ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
