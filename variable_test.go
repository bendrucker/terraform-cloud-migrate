package migrate

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform/configs"
)

func TestTerraforkWorkspaceStep_incomplete(t *testing.T) {
	parser := configs.NewParser(nil)
	path := "./fixtures/terraform-workspace/incomplete"
	mod, _ := parser.LoadConfigDir(path)

	step := TerraformWorkspaceStep{
		Module:   mod,
		Variable: "environment",
	}

	assert.False(t, step.Complete())

	changes, err := step.Changes()
	if !assert.NoError(t, err) {
		return
	}

	expectedOutputs := strings.TrimSpace(`
output "attribute" {
  value = var.environment
}

output "interpolated" {
  value = "The workspace is ${var.environment}"
}

output "function" {
  value = lookup({}, var.environment, false)
}	
`)

	expectedVariables := strings.TrimSpace(`
variable "environment" {
  type        = string
  description = "The environment where the module will be deployed"
}

variable "foo" {}
`)

	assert.Len(t, changes, 2)
	assert.Equal(t, expectedOutputs, strings.TrimSpace(string(changes[filepath.Join(path, "outputs.tf")].Bytes())))
	assert.Equal(t, expectedVariables, strings.TrimSpace(string(changes[filepath.Join(path, "variables.tf")].Bytes())))

}

func TestTerraforkWorkspaceStep_complete(t *testing.T) {
	parser := configs.NewParser(nil)
	path := "./fixtures/terraform-workspace/complete"
	mod, _ := parser.LoadConfigDir(path)

	step := TerraformWorkspaceStep{
		Module: mod,
	}

	assert.True(t, step.Complete())
}
