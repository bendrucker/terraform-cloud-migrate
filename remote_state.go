package migrate

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/configs"
)

type RemoteStateStep struct {
	module        *Module
	Path          string
	RemoteBackend RemoteBackendConfig
}

// Complete checks if any modules in the path are using remote_state
func (b *RemoteStateStep) Complete() bool {
	return false
}

// Description returns a description of the step
func (s *RemoteStateStep) Description() string {
	return `A "remote" backend should be configured for Terraform Cloud (https://www.terraform.io/docs/backends/types/remote.html)`
}

// Changes updates the configured backend
func (s *RemoteStateStep) Changes() (Changes, hcl.Diagnostics) {
	parser := configs.NewParser(nil)
	changes := Changes{}

	err := filepath.Walk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || !parser.IsConfigDir(path) {
			return nil
		}

		_, diags := s.moduleChanges(path)

		return diags
	})

	if err != nil {

	}

	return changes, nil

	// return Changes{path: &Change{File: file}}, diags
}

// Changes updates the configured backend
func (s *RemoteStateStep) moduleChanges(path string) (Changes, hcl.Diagnostics) {
	mod, diags := NewModule(path)
	mod.RemoteStateDataSources()
	return Changes{}, diags
}

var _ Step = (*RemoteStateStep)(nil)
