package main

import (
	"fmt"
	"bytes"
	"gopkg.in/ini.v1"
)

func (d *DevCube) CreateGitConfig() string {
	cfg := ini.Empty()

	if !d.Config.IsSet("source.email") || !d.Config.IsSet("source.name") {
		return ""
	}

	cfg.Section("user").Key("email").SetValue(d.Config.GetString("source.email"))
	cfg.Section("user").Key("name").SetValue(d.Config.GetString("source.name"))
	cfg.Section("credential").Key("helper").SetValue("store")
	cfg.Section("init").Key("defaultBranch").SetValue("master")


	var buffer bytes.Buffer

	cfg.WriteToIndent(&buffer, "\t")
	content := buffer.String()

	return content
}

func (d *DevCube) CreateGitCredentials() string {
	if !d.Config.IsSet("source.token") {
		return ""
	}

	user := d.Config.GetString("source.name")
	token := d.Config.GetString("source.token")
	provider := d.Config.GetString("source.provider")

	return fmt.Sprintf("https://%s:%s@%s", user, token, provider)
}
