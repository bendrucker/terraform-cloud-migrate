package migrate

import (
	"github.com/hashicorp/hcl/v2"
)

func New(path string, config Config) (*Migration, hcl.Diagnostics) {
	module, diags := NewModule(path)
	return &Migration{
		steps: []Step{
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
		},
	}, diags
}

type Migration struct {
	steps []Step
}

func (m *Migration) Changes() (Changes, hcl.Diagnostics) {
	changes := make(Changes)

	for _, step := range m.steps {
		stepChanges, err := step.Changes()
		if err != nil {
			return nil, hcl.Diagnostics{
				&hcl.Diagnostic{
					Detail: err.Error(),
				},
			}
		}

		for path, change := range stepChanges {
			if existing, ok := changes[path]; ok {
				if existing.Rename == "" {
					existing.Rename = change.Rename
				}
			}

			changes[path] = change
		}
	}

	return changes, nil
}
