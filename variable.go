package migrate

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

type TerraformWorkspaceStep struct {
	module   *Module
	Variable string
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
	primary, _, _ := parser.ConfigDirFiles(s.module.Dir())

	files := make(Changes)
	for _, path := range primary {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		file, _ := hclwrite.ParseConfig(bytes, path, hcl.InitialPos)
		replaceTerraformWorkspace(file.Body(), s.Variable)
		files[path] = &Change{File: file}
	}

	changes, err := changedFiles(parser.Sources(), files)
	if err != nil {
		return nil, err
	}

	if len(changes) == 0 {
		return changes, nil
	}

	if _, ok := s.module.Variables()[s.Variable]; !ok {
		path := filepath.Join(s.module.Dir(), "variables.tf")
		b, err := ioutil.ReadFile(path)

		var file *hclwrite.File
		if os.IsNotExist(err) {
			file = hclwrite.NewEmptyFile()
		} else {
			file, _ = hclwrite.ParseConfig(b, path, hcl.InitialPos)
		}

		changes[path] = &Change{
			File: addWorkspaceVariable(file, s.Variable),
		}
	}

	return changes, nil
}

func replaceTerraformWorkspace(body *hclwrite.Body, variable string) {
	for _, attr := range body.Attributes() {
		attr.Expr().RenameVariablePrefix(
			[]string{"terraform", "workspace"},
			[]string{"var", variable},
		)
	}

	for _, block := range body.Blocks() {
		replaceTerraformWorkspace(block.Body(), variable)
	}
}

func changedFiles(sources map[string][]byte, changes Changes) (Changes, error) {
	changed := make(Changes)

	for path, change := range changes {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		if bytes.Equal(b, change.File.Bytes()) {
			continue
		}

		changed[path] = &Change{File: change.File}
	}

	return changed, nil
}

func addWorkspaceVariable(file *hclwrite.File, name string) *hclwrite.File {
	blocks := file.Body().Blocks()
	file = hclwrite.NewEmptyFile()
	body := file.Body()

	variable := body.AppendBlock(hclwrite.NewBlock("variable", []string{name})).Body()
	body.AppendNewline()

	variable.SetAttributeRaw("type", hclwrite.Tokens{
		{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte("string"),
		},
	})
	variable.SetAttributeValue("description", cty.StringVal(fmt.Sprintf("The %s where the module will be deployed", name)))

	for _, block := range blocks {
		body.AppendBlock(block)
	}

	return file
}

var _ Step = (*TerraformWorkspaceStep)(nil)
