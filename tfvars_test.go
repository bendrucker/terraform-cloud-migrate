package migrate

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTfvarsStep_incomplete(t *testing.T) {
	path := "./fixtures/tfvars/incomplete"
	mod, diags := NewModule(path)
	if diags.HasErrors() {
		assert.Fail(t, diags.Error())
	}

	step := TfvarsStep{
		module:   mod,
		filename: "terraform.auto.tfvars",
	}

	assert.False(t, step.Complete())

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)

	change := changes[filepath.Join(path, "terraform.tfvars")]
	assert.Equal(t, "terraform.auto.tfvars", change.Rename)
}

func TestTfvarsStep_complete(t *testing.T) {
	path := "./fixtures/tfvars/complete"
	mod, diags := NewModule(path)
	assert.Len(t, diags, 0)

	step := TerraformWorkspaceStep{
		module: mod,
	}

	assert.True(t, step.Complete())
}
