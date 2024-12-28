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
	IsCreated() (bool, error)
	IsRunning() (bool, error)
	Build() (string, error)
	Create() (string, error)
	Start() (string, error)
	Exec(command []string, ex ExecType) (string, error)
}

type DevCube struct {
	ConfigPath	string
	Config  	*viper.Viper
	Engine  	Engine
}

var Version string = "1.0.0"
var Domain string = "devcube.io"
var configPath string
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

func (d *DevCube) PreRun(_ *cobra.Command, args []string) {
	d.ConfigPath = args[0]
	d.SetLogLevel()
	d.SetDefaults()
	d.ParseConfig()
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

	if created, _ := d.Engine.IsCreated(); !created {
		if _, err := d.Engine.Create(); err != nil {
			log.Fatal().Err(err).Msg("cannot create")
		}
	}
	
	if running, _ := d.Engine.IsRunning(); !running {
		if _, err := d.Engine.Start(); err != nil {
			log.Fatal().Err(err).Msg("cannot start")
		}
	}	

	if installed, _ := d.InstallDepends(); installed {
		if _, err := d.Engine.Exec([]string{"/usr/bin/nvim"}, Workspace); err != nil {
			log.Fatal().Err(err).Msg("cannot exec")
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

