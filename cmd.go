package main

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var log zerolog.Logger
var cube Cube

const (
	Version = "1.0.0"
)

var rootCmd = &cobra.Command{
	Use:		"cube",
	Version: 	Version,
	Short:		"cube is a devcontainer managment tool",
	Long:		``,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize devcube configuration",
	Args:  cobra.ExactArgs(1), 
	PersistentPreRun: cube.PreRun,
	Run:   cube.Init,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Open a devcube workspace",
	Args:  cobra.ExactArgs(1),
	PersistentPreRun: cube.PreRun,
	Run:   cube.Start,
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
}
