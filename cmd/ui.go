package main

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
)

// https://github.com/hashicorp/terraform/blob/master/command/cli_ui.go

type ColorizeUI struct {
	Colorize    *colorstring.Colorize
	OutputColor string
	InfoColor   string
	ErrorColor  string
	WarnColor   string
	UI          cli.Ui
}

func (u *ColorizeUI) Ask(query string) (string, error) {
	return u.UI.Ask(u.colorize(query, u.OutputColor))
}

func (u *ColorizeUI) AskSecret(query string) (string, error) {
	return u.UI.AskSecret(u.colorize(query, u.OutputColor))
}

func (u *ColorizeUI) Output(message string) {
	u.UI.Output(u.colorize(message, u.OutputColor))
}

func (u *ColorizeUI) Info(message string) {
	u.UI.Info(u.colorize(message, u.InfoColor))
}

func (u *ColorizeUI) Error(message string) {
	u.UI.Error(u.colorize(message, u.ErrorColor))
}

func (u *ColorizeUI) Warn(message string) {
	u.UI.Warn(u.colorize(message, u.WarnColor))
}

func (u *ColorizeUI) colorize(message string, color string) string {
	if color == "" {
		return message
	}

	return u.Colorize.Color(fmt.Sprintf("%s%s[reset]", color, message))
}
