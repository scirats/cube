package main

import (
	"os"
	"strings"

	"github.com/samber/lo"
)

type Podman struct {
	_ExecCmd func([]string, bool) (string, error)
	Name     string
	Image    string
	Envs     []string
	User     string
	Mount    string
	Ports    []string
	WorkDir  string
}

func (p *Podman) Init(c *DevCube) error {
	if err := os.Chdir(c.ConfigDir); err != nil {
		return err
	}

	p._ExecCmd = lo.Ternary(p._ExecCmd != nil, p._ExecCmd, execCmd)
	p.Name = c.Config.GetString("name")
	p.Image = c.Config.GetString("image")
	p.Envs = lo.MapToSlice(
		c.Config.GetStringMapString("environments"),
		func(k string, v string) string { return k + "=" + v },
	)
	p.User = c.Config.GetString("user")
	p.Mount = c.Config.GetString("workspaceMount")
	p.Ports = c.Config.GetStringSlice("forwardPorts")
	p.WorkDir = c.Config.GetString("workspaceFolder")

	return nil
}

func (p *Podman) IsBuilt() (bool, error) {
	cmdArgs := []string{"podman", "image", "ls"}
	cmdArgs = append(cmdArgs, "--quiet")
	cmdArgs = append(cmdArgs, "--format", "{{ .Repository }}")
	out, err := p._ExecCmd(cmdArgs, true)
	built := strings.TrimSpace(out) == p.Image

	return built, err
}

func (p *Podman) GetContainer(args ...string) (string, error) {
	cmdArgs := []string{"podman", "container", "ls"}
	cmdArgs = append(cmdArgs, "--quiet")
	cmdArgs = append(cmdArgs, "--latest")
	cmdArgs = append(cmdArgs, "--filter", "label=devcube.name="+p.Name)
	cmdArgs = append(cmdArgs, args...)

	return p._ExecCmd(cmdArgs, true)
}

func (d *Podman) Build() (string, error) {
	cmdArgs := []string{"podman", "image", "build"}
	cmdArgs = append(cmdArgs, "--tag", d.Name)
	cmdArgs = append(cmdArgs, "--file", d.Image)
	for _, arg := range d.Envs {
		cmdArgs = append(cmdArgs, "--build-arg", arg)
	}

	return d._ExecCmd(cmdArgs, false)
}

func (p *Podman) Start() (string, error) {
	container, _ := p.GetContainer()
	cmdArgs := []string{"podman", "container", "start"}
	cmdArgs = append(cmdArgs, container)

	return p._ExecCmd(cmdArgs, true)
}

func (p *Podman) Exec(command []string) (string, error) {
	container, _ := p.GetContainer()
	cmdArgs := []string{"podman", "container", "exec"}
	cmdArgs = append(cmdArgs, "--interactive", "--tty")
	cmdArgs = append(cmdArgs, "--workdir", p.WorkDir)
	cmdArgs = append(cmdArgs, "--user", p.User)
	cmdArgs = append(cmdArgs, container)
	cmdArgs = append(cmdArgs, command...)

	return p._ExecCmd(cmdArgs, false)
}
