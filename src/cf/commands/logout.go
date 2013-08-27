package commands

import (
	"cf/configuration"
	term "cf/terminal"
	"github.com/codegangsta/cli"
)

type Logout struct {
	ui     term.UI
	config *configuration.Configuration
}

func NewLogout(ui term.UI, config *configuration.Configuration) (l Logout) {
	l.ui = ui
	l.config = config
	return
}

func (l Logout) Run(c *cli.Context) {
	l.ui.Say("Logging out...")
	err := l.config.ClearSession()

	if err != nil {
		l.ui.Failed("Failed logging out", err)
		return
	}

	l.ui.Ok()
}