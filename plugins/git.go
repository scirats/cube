package plugins

import (
	"bytes"
	"strings"
	"gopkg.in/ini.v1"
	"github.com/scirats/exo"
)

type Git struct {
	Config		string
	Credentials	string
}

func (g *Git) resolveConfig(cfg *exo.Block) {
	if cfg.Has("config") {
		config := cfg.Block("config")
		root := ini.Empty()

		if config.Has("user") {
			user := config.Block("user") 
			if user.Has("email") && user.Has("name") {
				root.Section("user").Key("email").SetValue(user.String("email"))
				root.Section("user").Key("name").SetValue(user.String("name"))
			}
		}

		root.Section("credential").Key("helper").SetValue("store")
		root.Section("init").Key("defaultBranch").SetValue("master")


		var buffer bytes.Buffer
		root.WriteToIndent(&buffer, "\t")
		g.Config = buffer.String()
	}
}

func (g *Git) resolveCredentials(cfg *exo.Block) {
	if cfg.Has("credentials") {
		credentials := cfg.StringList("credentials")
		g.Credentials = strings.Join(credentials, "\n")
	}
}

func (g *Git) Configure(cfg *exo.Block) {
	if cfg.Has("git") {
		git := cfg.Block("git")
		g.resolveConfig(git)
		g.resolveCredentials(git)
	}
}

