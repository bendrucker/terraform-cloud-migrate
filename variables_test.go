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
		Module: mod,
	}

	assert.False(t, step.Complete())

	changes, err := step.Changes()
	if !assert.NoError(t, err) {
		return
	}

	expected := strings.TrimSpace(`
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

	assert.Len(t, changes, 1)
	assert.Equal(t, expected, strings.TrimSpace(string(changes[filepath.Join(path, "outputs.tf")].Bytes())))
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
