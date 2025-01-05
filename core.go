package main

import (
	"os"
	"github.com/spf13/viper"
	"github.com/spf13/cobra"
	"github.com/rs/zerolog"
)

const (
	Domain = "cube.io"
)

type Cube struct {
	Config		*viper.Viper
	Podman		*Podman
}

func (c *Cube) SetLogLevel() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func (c *Cube) SetDefaults() {
	c.Config = viper.New()
	c.Config.SetDefault("workspace", "/workspace")
	c.Config.SetDefault("engine", "podman")
}

func (c *Cube) SetEngine() {
	c.Podman = &Podman{
		Name: c.Config.GetString("name"),
		Engine: c.Config.GetString("engine"),
		WorkDir: c.Config.GetString("workspace"),
	}
}

func (c *Cube) PreRun(_ *cobra.Command, args []string) {
	c.SetLogLevel()
	c.SetDefaults()
	c.ParseConfig(args[0])
	c.SetEngine()
}

func (c *Cube) Init(cmd *cobra.Command, args []string) {
	
}

func (c *Cube) Start(cmd *cobra.Command, args []string) {
	if built, _ := c.Podman.IsBuilt(); !built {
		cfg := c.GenerateImage()
		if _, err := c.Podman.Build(cfg); err != nil {
			log.Fatal().Err(err).Msg("cannot build")
		}
	}

	if created, _ := c.Podman.IsCreated(); !created {
		if _, err := c.Podman.Create(); err != nil {
			log.Fatal().Err(err).Msg("cannot create")
		}
	}
	
	if running, _ := c.Podman.IsRunning(); !running {
		if _, err := c.Podman.Start(); err != nil {
			log.Fatal().Err(err).Msg("cannot start")
		}
	}	

	if installed := c.InstallDeps(); installed {
		if _, err := c.Podman.Exec("/usr/bin/nvim"); err != nil {
			log.Fatal().Err(err).Msg("cannot exec")
		}
	}
}
