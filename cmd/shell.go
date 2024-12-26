package cmd

import (
	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Execute a shell inside devcube",
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
