package main

import (
	"fmt"
	"bufio"
	"os"
	"os/exec"
	"bytes"
	"strings"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"github.com/samber/lo"
	"github.com/scirats/cube/xo"
)

var spec = &xo.Spec{
	Properties: []*xo.PropertySpec{
		{
			Type:    xo.TypeStringList,
			Name:    "ports",
			Repeat:  false,
			Require: false,
		},
		{
			Type:    xo.TypeString,
			Name:    "dots",
			Repeat:  false,
			Require: true,
		},
		{
			Type:    xo.TypeStringMultiline,
			Name:    "container",
			Repeat:  false,
			Require: true,
		},
	},
	Blocks: []*xo.BlockSpec{
		{
			Name:    "source",
			Repeat:  false,
			Require: false,
			Properties: []*xo.PropertySpec{
				{
					Type:    xo.TypeString,
					Name:    "repo",
					Repeat:  false,
					Require: false,
				},
				{
					Type:    xo.TypeString,
					Name:    "name",
					Repeat:  false,
					Require: false,
				},
				{
					Type:    xo.TypeString,
					Name:    "email",
					Repeat:  false,
					Require: false,
				},
				{
					Type:    xo.TypeString,
					Name:    "token",
					Repeat:  false,
					Require: false,
				},
			},
		},
	},
}

func execCmd(command []string, capture bool, buffer *string) (string, error) {
	var stdout []byte
	var err error

	cmd := exec.Command("podman", command...)
	
	if buffer != nil {
		cmd.Stdin = bytes.NewReader([]byte(*buffer))
	} else {
		cmd.Stdin = os.Stdin
	}
	
	cmd.Stderr = os.Stderr
	
	if capture {
		stdout, err = cmd.Output()
	} else {
		cmd.Stdout = os.Stdout
		err = cmd.Run()
	}

	return strings.TrimSpace(string(stdout)), err
}

func GetResult(out string) (string, bool, error) {
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	var records []string
	for scanner.Scan() {
		line := scanner.Text()
		records = append(records, line)
	}

	value, ok := lo.First(records)

	return value, ok, scanner.Err()
}

func (d *DevCube) StatContainer(path string) (bool, error) {
	cmdArgs := []string{"sh", "-c"}
	cmdArgs = append(cmdArgs, "test -d "+path+" && echo 1 || echo 0")
	out, err := d.Engine.Exec(cmdArgs, Capture)
	test, _, _ := GetResult(out)
	
	return test == "1", err
}

func (d *DevCube) SetLogLevel() {
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

func (d *DevCube) ParseConfig() {
	config, err := os.ReadFile(d.ConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read <file> settings")
	}

	cfg, err := xo.Parse(spec, string(config))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot parse config settings")
	}

	d.Config.Set("name", Domain+"/"+d.ConfigPath)
	d.Config.Set("ports", cfg.StringList("ports"))
	d.Config.Set("dots", cfg.String("dots"))
	d.Config.Set("container", cfg.String("container"))

	if cfg.Has("source") {
		block := cfg.Block("source")

		if block.Has("repo") {
			d.Config.Set("source.repo", block.String("repo"))
		}

		if block.Has("name") && block.Has("email") {
			d.Config.Set("source.name", block.String("name"))
			d.Config.Set("source.email", block.String("email"))
		}

		if block.Has("token") {
			d.Config.Set("source.token", block.String("token"))
		}
	}
}

func (d *DevCube) SetDefaults() {
	d.Config = viper.New()
	d.Config.SetDefault("workspaceFolder", "/workspace")
	d.Config.SetDefault("dotsFolder", "/root/.config")
	d.Config.SetDefault("gitFolder", "/root/.config/git")
	d.Config.SetDefault("gitConfigFile", "/root/.config/git/config")
	d.Config.SetDefault("gitCredentialsFile", "/root/.config/git/credentials")
}

func (d *DevCube) SetEngine() {
	d.Engine = &Podman{}

	if err := d.Engine.Init(d); err != nil {
		log.Fatal().Err(err).Msg("cannot initialize")
	}
}

func (d *DevCube) InstallDepends() (bool, error) {
	dotRepo := d.Config.GetString("dots")
	dotDir := d.Config.GetString("dotsFolder")
	if exists, _ := d.StatContainer(dotDir); !exists {
		if _, err := d.Engine.Exec([]string{
			"git", 
			"clone", 
			"https://github.com/"+dotRepo,
			dotDir,
		}, Hidden); err != nil {
			log.Fatal().Err(err).Msg("cannot create config directory")
		}

		gitDir := d.Config.GetString("gitFolder")
		if _, err := d.Engine.Exec([]string{"mkdir", gitDir}, Hidden); err != nil {
			log.Fatal().Err(err).Msg("cannot create git directory")
		}
		
		gitConfig := d.CreateGitConfig()
		gitConfigFile := d.Config.GetString("gitConfigFile")
		if len(gitConfig) > 0 {
			if _, err := d.Engine.Exec([]string{
				"sh",
				"-c",
				fmt.Sprintf(
					"cat > %s <<EOF\n%s\nEOF", 
					gitConfigFile, 
					gitConfig,
				),
			}, Hidden); err != nil {
				log.Fatal().Err(err).Msg("cannot create git config file")
			}
		}

			gitCredentials := d.CreateGitCredentials()
		gitCredentialsFile := d.Config.GetString("gitCredentialsFile")
		if len(gitCredentials) > 0 {
			if _, err := d.Engine.Exec([]string{
				"sh",
				"-c",
				fmt.Sprintf(
					"cat > %s <<EOF\n%s\nEOF", 
					gitCredentialsFile, 
					gitCredentials,
				),
			}, Hidden); err != nil {
				log.Fatal().Err(err).Msg("cannot create git credentials file")
			}
		}
	} 

	workDir := d.Config.GetString("workspaceFolder")
	if exists, _ := d.StatContainer(workDir); !exists {
		if d.Config.IsSet("source.repo") {
			projectRepo := d.Config.GetString("source.repo")
			if _, err := d.Engine.Exec([]string{
				"git", 
				"clone", 
				projectRepo,
				workDir,
			}, Hidden); err != nil {
				log.Fatal().Err(err).Msg("cannot create workspace directory")
			}

		} else {
			if exists, _ := d.StatContainer(workDir); !exists {
				if _, err := d.Engine.Exec([]string{"mkdir", workDir}, Hidden);
				err != nil {
					log.Fatal().Err(err).Msg("cannot create workspace directory")
				}
			}
		}
	}

	return true, nil
}
