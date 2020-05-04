package migrate

import (
	"github.com/bendrucker/terraform-cloud-migrate/configwrite"
	"github.com/hashicorp/hcl/v2"
)

func New(path string, config Config) (*Migration, hcl.Diagnostics) {
	writer, diags := configwrite.New(path)
	steps := configwrite.Steps{
		&configwrite.RemoteBackend{
			Writer: writer,
			Config: config.Backend,
		},
		&configwrite.TerraformWorkspace{
			Writer:   writer,
			Variable: config.WorkspaceVariable,
		},
		&configwrite.Tfvars{
			Writer:   writer,
			Filename: configwrite.TfvarsFilename,
		},
	}

	if config.ModulesDir != "" {
		steps = steps.Append(&configwrite.RemoteState{
			Writer:        writer,
			RemoteBackend: config.Backend,
			Path:          config.ModulesDir,
		})
	}

	return &Migration{steps}, diags
}

type Migration struct {
	steps configwrite.Steps
}

func (m *Migration) Changes() (configwrite.Changes, hcl.Diagnostics) {
	return m.steps.Changes()
}
