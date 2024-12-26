package cmd

import (
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize devcube configuration",
	// Run:   devc.Init,
}

func init() {
	rootCmd.AddCommand(initCmd)
}
