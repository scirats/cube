package main

import (
	"github.com/samber/lo"
)

type ExecType int

const (
	Hidden ExecType = iota
	Capture
	Workspace
)

type Podman struct {
	_ExecCmd func([]string, bool, *string) (string, error)
	Name     	string
	Container	string
	Ports    	[]string
	WorkDir  	string
}

func (p *Podman) createArgs() (cmdArgs []string) {
	for _, port := range p.Ports {
		cmdArgs = append(cmdArgs, "--publish", port)
	}

	return cmdArgs
}

func (p *Podman) Init(c *DevCube) error {
	p._ExecCmd = lo.Ternary(p._ExecCmd != nil, p._ExecCmd, execCmd)
	p.Name = c.Config.GetString("name")
	p.Container = c.Config.GetString("container")
	p.Ports = c.Config.GetStringSlice("ports")
	p.WorkDir = c.Config.GetString("workspaceFolder")

	return nil
}

func (p *Podman) GetImage(args ...string) (string, error) {
	cmdArgs := []string{"image", "ls"}
	cmdArgs = append(cmdArgs, "--filter", "reference="+p.Name)
	cmdArgs = append(cmdArgs, "--format", "{{.ID}}")
	cmdArgs = append(cmdArgs, args...)

	return p._ExecCmd(cmdArgs, true, nil)
} 

func (p *Podman) GetContainer(args ...string) (string, error) {
	cmdArgs := []string{"container", "ls"}
	cmdArgs = append(cmdArgs, "--latest")
	cmdArgs = append(cmdArgs, "--filter", "ancestor="+p.Name)
	cmdArgs = append(cmdArgs, "--format", "{{.ID}}")
	cmdArgs = append(cmdArgs, args...)

	return p._ExecCmd(cmdArgs, true, nil)
}

func (p *Podman) IsBuilt() (bool, error) {
	out, err := p.GetImage()
	_, exists, _ := GetResult(out)

	return exists, err
}

func (p *Podman) IsCreated() (bool, error) {
	out, err := p.GetContainer()
	_, exists, _ := GetResult(out)

	return exists, err
}

func (p *Podman) IsRunning() (bool, error) {
	out, err := p.GetContainer("--filter", "status=running")
	_, exists, _ := GetResult(out)

	return exists, err
}

func (p *Podman) Build() (string, error) {
	cmdArgs := []string{"image", "build"}
	cmdArgs = append(cmdArgs, "--tag", p.Name)
	cmdArgs = append(cmdArgs, "--file", "-")

	return p._ExecCmd(cmdArgs, false, &p.Container)
}

func (p *Podman) Create() (string, error) {
	cmdArgs := []string{"container", "create"}
	cmdArgs = append(cmdArgs, p.createArgs()...)
	cmdArgs = append(cmdArgs, p.Name)

	return p._ExecCmd(cmdArgs, true, nil)
}

func (p *Podman) Start() (string, error) {
	out, _ := p.GetContainer()
	container, _, _ := GetResult(out)
	cmdArgs := []string{"container", "start"}
	cmdArgs = append(cmdArgs, container)

	return p._ExecCmd(cmdArgs, true, nil)
}

func (p *Podman) Exec(command []string, ex ExecType) (string, error) {
	out, _ := p.GetContainer()
	container, _, _ := GetResult(out)
	cmdArgs := []string{"container", "exec"}
	cmdArgs = append(cmdArgs, "--interactive", "--tty")

	if ex == Workspace {
		cmdArgs = append(cmdArgs, "--workdir", p.WorkDir)
	}

	cmdArgs = append(cmdArgs, container)
	cmdArgs = append(cmdArgs, command...)

	return p._ExecCmd(cmdArgs, ex == Capture, nil)
}
