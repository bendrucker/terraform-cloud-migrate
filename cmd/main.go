package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	migrate "github.com/bendrucker/terraform-cloud-migrate"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "terraform-cloud-migrate",
		Usage: "migrate a Terraform module to Terraform Cloud",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "write",
				Usage: "Writes proposed changes to disk. When --init is set, --write is implicitly set.",
				Value: false,
			},
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
			&cli.StringFlag{
				Name:  "modules-dir",
				Usage: "A directory where other Terraform modules are stored. If set, it will be scanned recursively for terrafor_remote_state references.",
				Value: "",
			},
			&cli.BoolFlag{
				Name:  "init",
				Usage: "When set, 'terraform init' will be called before and after updating configuration. After updating the backend configuration, Terraform will prompt to copy state.",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			if !c.IsSet("workspace-name") && !c.IsSet("workspace-prefix") {
				return errors.New("one of --workspace-name or --workspace-prefix must be set")
			}

			path := c.Args().First()

			if c.Bool("init") {
				if err := terraformInit(path); err != nil {
					return fmt.Errorf("failed to init: %v", err)
				}
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
				ModulesDir:        c.String("modules-dir"),
			})

			if diags.HasErrors() {
				return diags
			}

			changes, diags := migration.Changes()
			if diags.HasErrors() {
				return diags
			}

			if c.Bool("write") || c.Bool("init") {
				for path, change := range changes {
					destination := path
					if change.Rename != "" {
						destination = filepath.Join(filepath.Dir(path), change.Rename)
					}

					file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
					if err != nil {
						return err
					}

					_, err = change.File.WriteTo(file)
					if err != nil {
						return err
					}

					if change.Rename != "" {
						os.Remove(path)
					}
				}
			} else {
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
			}

			if c.Bool("init") {
				if err := terraformInit(path); err != nil {
					return fmt.Errorf("failed to init: %v", err)
				}
			}

			fmt.Println("If your workspace is VCS-enabled, commit these changes and push to trigger a run")
			fmt.Println("If not, you can now call 'terraform plan' and 'terraform apply' locally.")

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func terraformInit(path string) error {
	cmd := exec.Command("terraform", "init")

	fmt.Fprint(os.Stderr, "Running 'terraform init'")

	cmd.Dir = path

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
