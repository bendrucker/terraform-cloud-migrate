package migrate

import (
	"github.com/hashicorp/hcl/v2"
)

func New(path string, config Config) (*Migration, hcl.Diagnostics) {
	module, diags := NewModule(path)
	steps := []Step{
		&RemoteBackendStep{
			module: module,
			Config: config.Backend,
		},
		&TerraformWorkspaceStep{
			module:   module,
			Variable: config.WorkspaceVariable,
		},
		&TfvarsStep{
			module:   module,
			filename: config.TfvarsFilename,
		},
	}

	if config.ModulesDir != "" {
		steps = append(steps, &RemoteStateStep{
			module:        module,
			RemoteBackend: config.Backend,
			Path:          config.ModulesDir,
		})
	}

	return &Migration{steps: steps}, diags
}

type Migration struct {
	steps []Step
}

func (m *Migration) Changes() (Changes, hcl.Diagnostics) {
	changes := make(Changes)
	diags := hcl.Diagnostics{}

	for _, step := range m.steps {
		stepChanges, sDiags := step.Changes()
		diags = append(diags, sDiags...)

		for path, change := range stepChanges {
			if existing, ok := changes[path]; ok {
				if existing.Rename == "" {
					existing.Rename = change.Rename
				}
			}

			changes[path] = change
		}
	}

	return changes, diags
}
