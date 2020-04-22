package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	migrate "github.com/bendrucker/terraform-cloud-migrate"
	"github.com/hashicorp/terraform/configs"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "hostname",
				Usage: "Hostname for Terraform Cloud",
				Value: "app.terraform.io",
			},
			&cli.StringFlag{
				Name:  "organization",
				Usage: "Organization name in Terraform Cloud",
			},
			&cli.StringFlag{
				Name:  "workspace-name",
				Usage: "Variable that will replace terraform.workspace",
			},
			&cli.StringFlag{
				Name:  "workspace-prefix",
				Usage: "Variable that will replace terraform.workspace",
			},
			&cli.StringFlag{
				Name:  "workspace-variable",
				Usage: "Variable that will replace terraform.workspace",
				Value: "environment",
			},
		},
		Action: func(c *cli.Context) error {
			path := c.Args().First()
			if empty, err := configs.IsEmptyDir(path); empty || err != nil {
				return fmt.Errorf("could not load Terraform files from %s", path)
			}

			parser := configs.NewParser(nil)
			mod, _ := parser.LoadConfigDir(path)

			steps := []migrate.Step{
				&migrate.RemoteBackendStep{
					Module: mod,
					Config: migrate.RemoteBackendConfig{
						Hostname:     c.String("hostname"),
						Organization: c.String("organization"),
						Workspaces: migrate.WorkspaceConfig{
							Prefix: c.String("workspace-prefix"),
							Name:   c.String("workspace-name"),
						},
					},
				},
				&migrate.TerraformWorkspaceStep{
					Module:   mod,
					Variable: c.String("workspace-variable"),
				},
				&migrate.TfvarsStep{Module: mod},
			}

			var changed bool
			for _, step := range steps {
				if step.Complete() {
					continue
				}

				changes, err := step.Changes()
				if err != nil {
					return err
				}

				for path, file := range changes {
					changed = true
					fmt.Fprintln(os.Stderr, "file: ", path, "desc: ", step.Description())
					file.WriteTo(os.Stderr)
					fmt.Fprint(os.Stderr, "\n\n")
				}
			}

			if changed {
				return errors.New("updates are required")
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
