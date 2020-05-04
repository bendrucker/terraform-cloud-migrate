package configwrite

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTfvars_incomplete(t *testing.T) {
	path := "./fixtures/tfvars/incomplete"
	mod, diags := New(path)
	if diags.HasErrors() {
		assert.Fail(t, diags.Error())
	}

	step := Tfvars{
		Writer:   mod,
		Filename: "terraform.auto.tfvars",
	}

	assert.False(t, step.Complete())

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)

	change := changes[filepath.Join(path, "terraform.tfvars")]
	assert.Equal(t, "terraform.auto.tfvars", change.Rename)
}

func TestTfvars_complete(t *testing.T) {
	path := "./fixtures/tfvars/complete"
	mod, diags := New(path)
	assert.Len(t, diags, 0)

	step := Tfvars{
		Writer: mod,
	}

	assert.True(t, step.Complete())
}
