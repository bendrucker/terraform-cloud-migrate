package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	migrate "github.com/bendrucker/terraform-cloud-migrate"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "terraform-cloud-migrate",
		Usage: "migrate a Terraform module to Terraform Cloud",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "hostname",
				Usage: "Hostname for Terraform Cloud",
				Value: "app.terraform.io",
			},
			&cli.StringFlag{
				Name:     "organization",
				Usage:    "Organization name in Terraform Cloud",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "workspace-name",
				Usage: "The name for the workspace",
			},
			&cli.StringFlag{
				Name:  "workspace-prefix",
				Usage: "The prefix for the workspaces",
			},
			&cli.StringFlag{
				Name:  "workspace-variable",
				Usage: "Variable that will replace terraform.workspace",
				Value: "environment",
			},
			&cli.StringFlag{
				Name:  "tfvars-filename",
				Usage: "New filename for terraform.tfvars",
				Value: migrate.TfvarsAlternateFilename,
			},
		},
		Action: func(c *cli.Context) error {
			if !c.IsSet("workspace-name") && !c.IsSet("workspace-prefix") {
				return errors.New("one of --workspace-name or --workspace-prefix must be set")
			}

			migration, diags := migrate.New(c.Args().First(), migrate.Config{
				Backend: migrate.RemoteBackendConfig{
					Hostname:     c.String("hostname"),
					Organization: c.String("organization"),
					Workspaces: migrate.WorkspaceConfig{
						Prefix: c.String("workspace-prefix"),
						Name:   c.String("workspace-name"),
					},
				},
				WorkspaceVariable: c.String("workspace-variable"),
				TfvarsFilename:    c.String("tfvars-filename"),
			})

			if diags.HasErrors() {
				return diags
			}

			changes, diags := migration.Changes()
			if diags.HasErrors() {
				return diags
			}

			for path, change := range changes {
				var rename string
				if change.Rename != "" {
					rename = fmt.Sprintf("(moved to %s)", change.Rename)
				}

				fmt.Fprintln(os.Stderr, "# file: ", path, rename)
				change.File.WriteTo(os.Stderr)
				fmt.Fprint(os.Stderr, "\n")
			}

			if len(changes) != 0 {
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
