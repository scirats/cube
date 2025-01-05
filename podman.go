package main

import (
	"github.com/samber/lo"
)

type Podman struct {
	Name	string
	Engine	string
	WorkDir string
}

func (p *Podman) GetImages(args ...string) ([]string, error) {
	cmdArgs := []string{p.Engine, "image", "ls"}
	cmdArgs = append(cmdArgs, "--filter", "reference="+p.Name)
	cmdArgs = append(cmdArgs, "--format", "{{.ID}}")
	cmdArgs = append(cmdArgs, args...)

	out, _ := ExecCmd(cmdArgs, true)
	
	return GetResults(out)
}

func (p *Podman) GetContainers(args ...string) ([]string, error) {
	cmdArgs := []string{p.Engine, "container", "ls"}
	cmdArgs = append(cmdArgs, "--latest")
	cmdArgs = append(cmdArgs, "--filter", "ancestor="+p.Name)
	cmdArgs = append(cmdArgs, "--format", "{{.ID}}")
	cmdArgs = append(cmdArgs, args...)

	out, _ := ExecCmd(cmdArgs, true)

	return GetResults(out)
}

func (p *Podman) IsBuilt() (bool, error) {
	images, err := p.GetImages()
	_, exists := lo.First(images)

	return exists, err
}

func (p *Podman) IsCreated() (bool, error) {
	containers, err := p.GetContainers()
	_, exists := lo.First(containers)

	return exists, err
}

func (p *Podman) IsRunning() (bool, error) {
	containers, err := p.GetContainers("--filter", "status=running")
	_, exists := lo.First(containers)

	return exists, err
}


func (p *Podman) Build(imageConfig string) (string, error) {
	cmdArgs := []string{p.Engine, "image", "build"}
	cmdArgs = append(cmdArgs, "--tag", p.Name)
	cmdArgs = append(cmdArgs, "--file", "-")
	
	return ExecCmdWithBuffer(cmdArgs, imageConfig, false)
}

func (p *Podman) Create() (string, error) {
	cmdArgs := []string{p.Engine, "container", "create"}
	cmdArgs = append(cmdArgs, "--net", "host")
	cmdArgs = append(cmdArgs, p.Name)

	return ExecCmd(cmdArgs, false)
}

func (p *Podman) Start() (string, error) {
	containers, _ := p.GetContainers()
	container, _ := lo.First(containers)

	cmdArgs := []string{p.Engine, "container", "start"}
	cmdArgs = append(cmdArgs, container)

	return ExecCmd(cmdArgs, false)
}

func (p *Podman) PreExec(args ...string) (string, error) {
	containers, _ := p.GetContainers()
	container, _ := lo.First(containers)

	cmdArgs := []string{p.Engine, "container", "exec"}
	cmdArgs = append(cmdArgs, container)
	cmdArgs = append(cmdArgs, args...)

	return ExecCmd(cmdArgs, true)
}

func (p *Podman) Exec(args ...string) (string, error) {
	containers, _ := p.GetContainers()
	container, _ := lo.First(containers)

	cmdArgs := []string{p.Engine, "container", "exec"}
	cmdArgs = append(cmdArgs, "--interactive", "--tty")
	cmdArgs = append(cmdArgs, "--workdir", p.WorkDir)
	cmdArgs = append(cmdArgs, container)
	cmdArgs = append(cmdArgs, args...)

	return ExecCmd(cmdArgs, false)
}
