package migrate

import (
	"context"
	"fmt"

	"github.com/bendrucker/terraform-cloud-migrate/configwrite"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcl/v2"
)

func New(path string, config Config) (*Migration, hcl.Diagnostics) {
	writer, diags := configwrite.New(path)
	steps := configwrite.NewSteps(writer, configwrite.Steps{
		&configwrite.RemoteBackend{Config: config.Backend},
		&configwrite.TerraformWorkspace{Variable: config.WorkspaceVariable},
		&configwrite.Tfvars{Filename: configwrite.TfvarsFilename},
	})

	if config.ModulesDir != "" {
		step := &configwrite.RemoteState{
			RemoteBackend: config.Backend,
			Path:          config.ModulesDir,
		}
		step.WithWriter(writer)
		steps = steps.Append(step)
	}

	client, err := tfe.NewClient(config.API)
	if err != nil {
		return nil, diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "Terraform Cloud API client error",
			Detail:   fmt.Sprintf("Terraform Cloud API client could not be created: %v", err),
		})
	}

	return &Migration{
		config: config,
		steps:  steps,
		api:    client,
	}, diags
}

type Migration struct {
	config Config
	steps  configwrite.Steps
	api    *tfe.Client
}

func (m *Migration) GetWorkspaces() ([]*tfe.Workspace, error) {
	list, err := m.api.Workspaces.List(context.TODO(), m.config.Backend.Organization, tfe.WorkspaceListOptions{
		Search: tfe.String(m.config.Backend.Workspaces.Prefix + m.config.Backend.Workspaces.Name),
	})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (m *Migration) Changes() (configwrite.Changes, hcl.Diagnostics) {
	return m.steps.Changes()
}
