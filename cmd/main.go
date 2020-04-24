package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

func main() {
	app := cli.NewCLI("terraform-cloud-migrate", "0.0.0")
	app.Args = os.Args[1:]

	ui := &cli.ColoredUi{
		OutputColor: cli.UiColorNone,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
		WarnColor:   cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	app.Commands = map[string]cli.CommandFactory{
		"run": func() (cli.Command, error) {
			return &RunCommand{
				Ui: ui,
			}, nil
		},
	}

	status, err := app.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}
