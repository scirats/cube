package main

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log zerolog.Logger

type Engine interface {
	Init(config *DevCube) error
	IsBuilt() (bool, error)
	Build() (string, error)
	Start() (string, error)
	Exec(command []string) (string, error)
}

type DevCube struct {
	ConfigDir            string
	Config               *viper.Viper
	Engine               Engine
	WorkingDirectoryPath string
	WorkingDirectoryName string
}

var Version string = "1.0.0"
var rootConfigDir string = ".devcube"
var cube DevCube

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
}

var rootCmd = &cobra.Command{
	Use:              "cube",
	Version:          Version,
	Short:            "cube is a devcontainer managment tool",
	Long:             ``,
	PersistentPreRun: cube.PreRun,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize devcube configuration",
	Run:   cube.Init,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Open a devcube workspace",
	Run:   cube.Start,
}

func (d *DevCube) PreRun(_ *cobra.Command, _ []string) {
	d.SetLogLevel()
	d.ParseConfig()
	d.SetDefaults()
	d.SetEngine()
}

func (d *DevCube) Init(cmd *cobra.Command, args []string) {
	fmt.Println("init")
}

func (d *DevCube) Start(cmd *cobra.Command, args []string) {
	if built, _ := d.Engine.IsBuilt(); !built {
		if _, err := d.Engine.Build(); err != nil {
			log.Fatal().Err(err).Msg("cannot build")
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
