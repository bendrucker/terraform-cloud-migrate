package migrate

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformWorkspaceStep_incomplete(t *testing.T) {
	path := "./fixtures/terraform-workspace/incomplete"
	mod, diags := NewModule(path)

	if diags.HasErrors() {
		assert.Fail(t, diags.Error())
	}

	step := TerraformWorkspaceStep{
		module:   mod,
		Variable: "environment",
	}

	assert.False(t, step.Complete())

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)

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
	assert.Equal(t, expectedOutputs, strings.TrimSpace(string(changes[filepath.Join(path, "outputs.tf")].File.Bytes())))
	assert.Equal(t, expectedVariables, strings.TrimSpace(string(changes[filepath.Join(path, "variables.tf")].File.Bytes())))

}

func TestTerraformWorkspaceStep_complete(t *testing.T) {
	path := "./fixtures/terraform-workspace/complete"
	mod, diags := NewModule(path)

	if diags.HasErrors() {
		assert.Fail(t, diags.Error())
	}

	step := TerraformWorkspaceStep{
		module: mod,
	}

	assert.True(t, step.Complete())
}
