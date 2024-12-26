package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"muzzammil.xyz/jsonc"
)

func execCmd(command []string, capture bool) (string, error) {
	var stdout []byte
	var err error

	cwd, _ := os.Getwd()
	cmd := exec.Command(command[0], command[1:]...)
	log.Info().Str("workdir", cwd).Str("command", cmd.String()).Send()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if capture {
		stdout, err = cmd.Output()
	} else {
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}

	return strings.TrimSpace(string(stdout)), err
}

func (d *DevCube) SetLogLevel() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func (d *DevCube) ParseConfig() {
	_, j, err := jsonc.ReadFromFile(filepath.Join(rootConfigDir, "devcube.json"))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read devcube settings")
	}

	d.Config = viper.New()
	d.Config.SetConfigType("json")
	if err := d.Config.ReadConfig(bytes.NewBuffer(j)); err != nil {
		log.Fatal().Err(err).Msg("cannot read json")
	}

	log.Debug().Str("devcube", fmt.Sprintf("%+v", d.Config)).Send()
}

func (d *DevCube) SetDefaults() {
	d.ConfigDir = rootConfigDir
	d.WorkingDirectoryPath, _ = os.Getwd()
	d.WorkingDirectoryName = filepath.Base(d.WorkingDirectoryPath)
	d.Config.SetDefault("name", d.WorkingDirectoryName)
	d.Config.SetDefault("user", "cube")
	d.Config.SetDefault("workspaceFolder", "/workspace")
	d.Config.SetDefault("workspaceMount", "type=bind,source="+d.WorkingDirectoryPath+",target="+d.Config.GetString("workspaceFolder")+",consistency=cached")
}

func (d *DevCube) SetEngine() {
	d.Engine = &Podman{}

	if err := d.Engine.Init(d); err != nil {
		log.Fatal().Err(err).Msg("cannot initialize")
	}
	log.Debug().Str("engine", fmt.Sprintf("%+v", d.Engine)).Send()
}
