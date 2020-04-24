package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	migrate "github.com/bendrucker/terraform-cloud-migrate"
	"github.com/mitchellh/cli"

	flag "github.com/spf13/pflag"
)

type RunCommand struct {
	Ui cli.Ui
}

func (c *RunCommand) Run(args []string) int {
	var hostname, organization, name, prefix, variable, tfvarsName, modules string
	var noInit bool

	flags := flag.NewFlagSet("run", flag.ContinueOnError)

	flags.StringVar(&hostname, "hostname", "app.terraform.io", "Hostname for Terraform Cloud")
	flags.StringVar(&organization, "organization", "", "Organization name in Terraform Cloud")
	flags.StringVar(&name, "workspace-name", "", "The name of the Terraform Cloud workspace (conflicts with --workspace-prefix)")
	flags.StringVar(&prefix, "workspace-prefix", "", "The prefix of the Terraform Cloud workspaces (conflicts with --workspace-name)")
	flags.StringVar(&variable, "workspace-variable", "environment", "Variable that will replace terraform.workspace")
	flags.StringVar(&tfvarsName, "tfvars-filename", migrate.TfvarsAlternateFilename, "New filename for terraform.tfvars")
	flags.StringVar(&modules, "modules", "", "A directory where other Terraform modules are stored. If set, it will be scanned recursively for terrafor_remote_state references.")
	flags.BoolVar(&noInit, "no-init", false, "Disable calling 'terraform init' before and after updating configuration to copy state.")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if len(flags.Args()) != 1 {
		c.Ui.Error("module path is required")
		return 1
	}

	path := flags.Args()[0]
	abspath, err := filepath.Abs(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to resolve path: %s", path))
		return 1
	}

	c.Ui.Info(fmt.Sprintf("Upgrading Terraform module %s", abspath))

	if name == "" && prefix == "" {
		c.Ui.Error("workspace name or prefix is required")
		return 1
	}

	if name != "" && prefix != "" {
		c.Ui.Error("workspace cannot have a name and prefix")
		return 1
	}

	migration, diags := migrate.New(path, migrate.Config{
		Backend: migrate.RemoteBackendConfig{
			Hostname:     hostname,
			Organization: organization,
			Workspaces: migrate.WorkspaceConfig{
				Prefix: prefix,
				Name:   name,
			},
		},
		WorkspaceVariable: variable,
		TfvarsFilename:    tfvarsName,
		ModulesDir:        modules,
	})

	if diags.HasErrors() {
		return c.fail(diags)
	}

	changes, diags := migration.Changes()
	if diags.HasErrors() {
		return c.fail(diags)
	}

	if !noInit {
		c.Ui.Info("Running 'terraform init' prior to updating backend")
		c.Ui.Info("This ensures that Terraform has persisted the existing backend configuration to local state")

		if code := c.terraformInit(abspath); code != 0 {
			return code
		}
	}

	for path, change := range changes {
		destination := path
		if change.Rename != "" {
			destination = filepath.Join(filepath.Dir(path), change.Rename)
		}

		file, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return c.fail(err)
		}

		_, err = change.File.WriteTo(file)
		if err != nil {
			return c.fail(err)
		}

		if change.Rename != "" {
			os.Remove(path)
		}
	}

	if !noInit {
		c.Ui.Info("Running 'terraform init' to copy state")
		c.Ui.Info("When prompted, type 'yes' to confirm")

		if code := c.terraformInit(abspath); code != 0 {
			return code
		}
	}

	c.Ui.Info("Migration complete!")
	c.Ui.Info("If your workspace is VCS-enabled, commit these changes and push to trigger a run.")
	c.Ui.Info("If not, you can now call 'terraform plan' and 'terraform apply' locally.")

	return 0
}

func (c *RunCommand) Help() string {
	return "Run Terraform Cloud migration"
}

func (c *RunCommand) Synopsis() string {
	return "Run Terraform Cloud migration"
}

func (c *RunCommand) fail(err error) int {
	c.Ui.Error(err.Error())
	return 1
}

func (c *RunCommand) terraformInit(path string) int {
	cmd := exec.Command("terraform", "init", path)

	cmd.Dir = path

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return err.ExitCode()
		}

		c.Ui.Error(fmt.Sprintf("failed to terraform init: %v", err))
	}

	return 0
}
