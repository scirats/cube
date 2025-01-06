package main

import (
	"fmt"
	"strings"
	"reflect"
	"github.com/scirats/exo"
	"github.com/samber/lo"
	ps "github.com/scirats/cube/plugins"
)

const (
	dotsDir = "/root/.config"
)

type Plugin interface {
	Configure(cfg *exo.Block)
}

func (c *Cube) ParseConfig(name string) {
	if cfg, err := exo.ParseFile(name); cfg != nil {
		c.Config.Set("name", Domain+"/"+name)
		
		if cfg.Has("deps") {
			c.Config.Set("pre.deps", cfg.StringList("deps"))
		}

		if cfg.Has("workspace") {
			c.Config.Set("post.workspace", cfg.String("workspace"))
		}
		
		if cfg.Has("dots") {
			dots := cfg.Block("dots")
			if dots.Has("repo") {
				c.Config.Set("post.dots.repo", dots.String("repo"))
			}

			if dots.Has("ext") {
				ext := dots.Block("ext")
				c.Config.Set("post.dots.extMeta", ext)
				c.Config.Set("post.dots.ext", []Plugin{&ps.Git{}})
			}
		}
	} else {
		log.Fatal().Err(err).Msg("cannot parse config file")
	}
}

func (c *Cube) GenerateImage() string {
	deps := c.Config.GetStringSlice("pre.deps")

	mods := strings.Join(deps, " ")

	imageCfg := fmt.Sprintf(`
	FROM alpine:latest

	RUN apk update && \
		apk add --no-cache neovim git %s
		
	CMD ["sleep", "infinity"]
	`, mods)

	return imageCfg
}

func (c *Cube) IsInstalled() bool {
	workDir := c.Config.GetString("workspace")

	cmdArgs := []string{"sh", "-c"}
	cmdArgs = append(cmdArgs, "test -d "+workDir+" && echo 1 || echo 0")
	out, _ := c.Podman.PreExec(cmdArgs...)
	records, _ := GetResults(out)
	test, _ := lo.First(records)

	return test == "1"
}

func (c *Cube) InstallWorkspace() {
	workDir := c.Config.GetString("workspace")
	cmdArgs := []string{}
	if c.Config.IsSet("post.workspace") {
		gitRepo := c.Config.GetString("post.workspace")
		cmdArgs = append(cmdArgs, "git", "clone", gitRepo, workDir)
	} else {
		cmdArgs = append(cmdArgs, "mkdir", workDir)
	}

	if _, err := c.Podman.PreExec(cmdArgs...); err != nil {
		log.Fatal().Err(err).Msg("could not install or create workspace")		
	}
}

func (c *Cube) InstallDots() {
	cmdArgs := []string{}
	if c.Config.IsSet("post.dots.repo") {
		gitRepo := c.Config.GetString("post.dots.repo")
		cmdArgs = append(cmdArgs, "git", "clone", gitRepo, dotsDir)
	} else {
		cmdArgs = append(cmdArgs, "mkdir", dotsDir)
	}

	if _, err := c.Podman.PreExec(cmdArgs...); err != nil {
		log.Fatal().Err(err).Msg("could not install or create dots files")		
	}
}

func (c *Cube) InstallExtDir(folder string)  {
	cmdFolderArgs := []string{"mkdir", dotsDir+"/"+folder}
	if _, err := c.Podman.PreExec(cmdFolderArgs...); err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("could not create %s folder", folder))
	}	
}

func (c *Cube) InstallExtFile(folder string, file string, content string)  {
	cmdFileArgs := []string{
		"sh",
		"-c",
		fmt.Sprintf(
			"cat > %s <<EOF\n%s\nEOF", 
			dotsDir+"/"+folder+"/"+file, 
			content,
		),
	}
	if _, err := c.Podman.PreExec(cmdFileArgs...); err != nil {
		log.Fatal().Err(err).Msg(fmt.Sprintf("could not create %s file", file))
	}
}

func (c *Cube) InstallDeps() bool {
	isInstalled := c.IsInstalled()
	if !isInstalled {
		c.InstallDots()

		if c.Config.IsSet("post.dots.ext") {
			cfg := c.Config.Get("post.dots.extMeta").(*exo.Block)
			plugins := c.Config.Get("post.dots.ext").([]Plugin)
			for _, plugin := range plugins {
				plugin.Configure(cfg)
				
				v := reflect.ValueOf(plugin)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}
				t := v.Type()
				
				folder := strings.ToLower(t.Name())
				c.InstallExtDir(folder)
				
				for i := 0; i < v.NumField(); i++ {
					field := t.Field(i)
					value := v.Field(i)

					c.InstallExtFile(
						folder, 
						strings.ToLower(field.Name),
						value.String(),
					)
				}
			}
		}

		c.InstallWorkspace()
		
		isInstalled = true
	}

	return isInstalled
}
