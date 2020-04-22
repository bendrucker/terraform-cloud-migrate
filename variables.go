package migrate

import (
	"bytes"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
)

type TerraformWorkspaceStep struct {
	Module *configs.Module
}

// Complete checks if any terraform.workspace replaces are proposed
func (s *TerraformWorkspaceStep) Complete() bool {
	changes, _ := s.Changes()
	return len(changes) == 0
}

// Description returns a description of the step
func (s *TerraformWorkspaceStep) Description() string {
	return `terraform.workpace will always be set to default and should not be used with Terraform Cloud (https://www.terraform.io/docs/state/workspaces.html#current-workspace-interpolation)`
}

// Changes determines changes required to remove terraform.workspace
func (s *TerraformWorkspaceStep) Changes() (Changes, error) {
	parser := configs.NewParser(nil)
	primary, _, _ := parser.ConfigDirFiles(s.Module.SourceDir)

	files := make(Changes)
	for _, path := range primary {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		file, _ := hclwrite.ParseConfig(bytes, path, hcl.InitialPos)
		removeTerraformWorkspace(file.Body())
		files[path] = file
	}

	return changedFiles(parser.Sources(), files)
}

func removeTerraformWorkspace(body *hclwrite.Body) {
	for _, attr := range body.Attributes() {
		attr.Expr().RenameVariablePrefix(
			[]string{"terraform", "workspace"},
			[]string{"var", "environment"},
		)
	}

	for _, block := range body.Blocks() {
		removeTerraformWorkspace(block.Body())
	}
}

func changedFiles(sources map[string][]byte, files Changes) (Changes, error) {
	changed := make(Changes)

	for path, file := range files {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if bytes.Equal(b, file.Bytes()) {
			continue
		}

		changed[path] = file
	}

	return changed, nil
}

var _ Step = (*TerraformWorkspaceStep)(nil)
