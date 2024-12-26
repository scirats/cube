package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

// import (
// 	"os"
// 	"os/exec"
// 	"strings"

// 	"github.com/rs/zerolog"
// )

// func execCmd(command []string, capture bool) (string, error) {
// 	var stdout []byte
// 	var err error

// 	cwd, _ := os.Getwd()
// 	cmd := exec.Command(command[0], command[1:]...)
// 	log.Info().Str("workdir", cwd).Str("command", cmd.String()).Send()
// 	cmd.Stdin = os.Stdin
// 	cmd.Stderr = os.Stderr
// 	if capture {
// 		stdout, err = cmd.Output()
// 	} else {
// 		cmd.Stdout = os.Stdout
// 		err = cmd.Run()
// 	}

// 	return strings.TrimSpace(string(stdout)), err
// }

func (d *DevCube) SetLogLevel() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func (d *DevCube) SetEngine() {
	d.Engine = &Podman{}

	if err := d.Engine.Init(d); err != nil {
		log.Fatal().Err(err).Msg("cannot initialize")
	}
	log.Debug().Str("engine", fmt.Sprintf("%+v", d.Engine)).Send()
}
