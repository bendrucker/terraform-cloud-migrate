package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	migrate "github.com/bendrucker/terraform-cloud-migrate"
	"github.com/bendrucker/terraform-cloud-migrate/steps"
	"github.com/mitchellh/cli"

	"github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
)

func NewRunCommand(ui cli.Ui) cli.Command {
	rc := &RunCommand{
		Ui:    ui,
		Flags: flag.NewFlagSet("run", flag.ContinueOnError),
	}

	rc.Flags.SortFlags = false
	c := rc.Config
	rc.Flags.StringVarP(&c.WorkspaceName, "workspace-name", "n", "", "The name of the Terraform Cloud workspace (conflicts with --workspace-prefix)")
	rc.Flags.StringVarP(&c.WorkspacePrefix, "workspace-prefix", "p", "", "The prefix of the Terraform Cloud workspaces (conflicts with --workspace-name)")
	rc.Flags.StringVarP(&c.ModulesDir, "modules", "m", "", "A directory where other Terraform modules are stored. If set, it will be scanned recursively for terrafor_remote_state references.")
	rc.Flags.StringVar(&c.WorkspaceVariable, "workspace-variable", "environment", "Variable that will replace terraform.workspace")
	rc.Flags.StringVar(&c.TfvarsFilename, "tfvars-filename", steps.TfvarsAlternateFilename, "New filename for terraform.tfvars")

	rc.Flags.StringVar(&c.Hostname, "hostname", "app.terraform.io", "Hostname for Terraform Cloud")
	rc.Flags.StringVar(&c.Organization, "organization", "", "Organization name in Terraform Cloud")

	rc.Flags.BoolVar(&c.NoInit, "no-init", false, "Disable calling 'terraform init' before and after updating configuration to copy state.")

	return rc
}

type RunCommand struct {
	Flags  *pflag.FlagSet
	Config RunCommandConfig
	Ui     cli.Ui
}

type RunCommandConfig struct {
	Hostname          string
	Organization      string
	WorkspaceName     string
	WorkspacePrefix   string
	WorkspaceVariable string
	TfvarsFilename    string
	ModulesDir        string
	NoInit            bool
}

func (c *RunCommand) Run(args []string) int {
	var hostname, organization, name, prefix, variable, tfvarsName, modules string
	var noInit bool

	if err := c.Flags.Parse(args); err != nil {
		return 1
	}

	if len(c.Flags.Args()) != 1 {
		c.Ui.Error("module path is required")
		return 1
	}

	path := c.Flags.Args()[0]
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
		Backend: steps.RemoteBackendConfig{
			Hostname:     hostname,
			Organization: organization,
			Workspaces: steps.WorkspaceConfig{
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
		fmt.Println()

		if code := c.terraformInit(abspath); code != 0 {
			return code
		}
	}

	if err := changes.WriteFiles(); err != nil {
		return c.fail(err)
	}

	for path, change := range changes {
		str := path
		if change.Rename != "" {
			str = fmt.Sprintf("%s -> %s", path, change.Destination(path))
		}

		fmt.Println(str)
	}

	if !noInit {
		c.Ui.Info("Running 'terraform init' to copy state")
		c.Ui.Info("When prompted, type 'yes' to confirm")
		fmt.Println()

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
	return strings.TrimSpace(`
Usage: terraform-cloud-migrate run [DIR] [options]
  Migrate a Terraform module to Terraform Cloud

Options:
` + c.Flags.FlagUsages())
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
