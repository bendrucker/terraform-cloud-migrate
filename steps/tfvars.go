package steps

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
)

const (
	TfvarsFilename          = "terraform.tfvars"
	TfvarsAlternateFilename = "terraform.auto.tfvars"
)

type Tfvars struct {
	Writer   *Writer
	Filename string
}

func (s *Tfvars) Name() string {
	return "Rename terraform.tfvars"
}

// Complete checks if a terraform.tfvars file exists and returns false if it does
func (s *Tfvars) Complete() bool {
	_, err := ioutil.ReadFile(s.path(TfvarsFilename))
	return err != nil && os.IsNotExist(err)
}

// Description returns a description of the step
func (s *Tfvars) Description() string {
	return `Terraform Cloud passes workspace variables by writing to terraform.tfvars and will overwrite existing content (terraform.workpace will always be set to default and should not be used with Terraform Cloud (https://www.terraform.io/docs/state/workspaces.html#current-workspace-interpolation)`
}

func (s *Tfvars) path(filename string) string {
	return filepath.Join(s.Writer.Dir(), filename)
}

// Changes determines changes required to remove terraform.workspace
func (s *Tfvars) Changes() (Changes, hcl.Diagnostics) {
	if s.Complete() {
		return Changes{}, nil
	}

	existing := s.path(TfvarsFilename)
	file, diags := s.Writer.File(existing)

	return Changes{
		existing: &Change{
			File:   file,
			Rename: s.Filename,
		},
	}, diags
}

var _ Step = (*Tfvars)(nil)
