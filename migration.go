package migrate

import (
	filesteps "github.com/bendrucker/terraform-cloud-migrate/steps"
	"github.com/hashicorp/hcl/v2"
)

func New(path string, config Config) (*Migration, hcl.Diagnostics) {
	writer, diags := filesteps.NewWriter(path)
	steps := filesteps.Steps{
		&filesteps.RemoteBackend{
			Writer: writer,
			Config: config.Backend,
		},
		&filesteps.TerraformWorkspace{
			Writer:   writer,
			Variable: config.WorkspaceVariable,
		},
		&filesteps.Tfvars{
			Writer:   writer,
			Filename: config.TfvarsFilename,
		},
	}

	if config.ModulesDir != "" {
		steps = steps.Append(&filesteps.RemoteState{
			Writer:        writer,
			RemoteBackend: config.Backend,
			Path:          config.ModulesDir,
		})
	}

	return &Migration{steps}, diags
}

type Migration struct {
	steps filesteps.Steps
}

func (m *Migration) Changes() (filesteps.Changes, hcl.Diagnostics) {
	return m.steps.Changes()
}
