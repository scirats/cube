package cmd

import (
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Open a devcube workspace",
}

func init() {
	rootCmd.AddCommand(startCmd)
}
