package migrate

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
)

const (
	TfvarsFilename          = "terraform.tfvars"
	TfvarsAlternateFilename = "default.auto.tfvars"
)

type TfvarsStep struct {
	Module *configs.Module
}

// Complete checks if a terraform.tfvars file exists and returns false if it does
func (s *TfvarsStep) Complete() bool {
	_, err := ioutil.ReadFile(s.path(TfvarsFilename))
	return err != nil && os.IsNotExist(err)
}

// Description returns a description of the step
func (s *TfvarsStep) Description() string {
	return `Terraform Cloud passes workspace variables by writing to terraform.tfvars and will overwrite existing content (terraform.workpace will always be set to default and should not be used with Terraform Cloud (https://www.terraform.io/docs/state/workspaces.html#current-workspace-interpolation)`
}

func (s *TfvarsStep) path(filename string) string {
	return filepath.Join(s.Module.SourceDir, filename)
}

// Changes determines changes required to remove terraform.workspace
func (s *TfvarsStep) Changes() (Changes, error) {
	if s.Complete() {
		return Changes{}, nil
	}

	existing := s.path(TfvarsFilename)

	bytes, err := ioutil.ReadFile(existing)
	if err != nil {
		return nil, err
	}

	file, _ := hclwrite.ParseConfig(bytes, existing, hcl.InitialPos)

	return Changes{
		existing:                        nil,
		s.path(TfvarsAlternateFilename): file,
	}, nil
}

var _ Step = (*TfvarsStep)(nil)
