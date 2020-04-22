package migrate

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform/configs"
)

func TestTfvarsStep_incomplete(t *testing.T) {
	parser := configs.NewParser(nil)
	path := "./fixtures/tfvars/incomplete"
	mod, _ := parser.LoadConfigDir(path)

	step := TfvarsStep{
		Module: mod,
	}

	assert.False(t, step.Complete())

	changes, err := step.Changes()
	if !assert.NoError(t, err) {
		return
	}

	file, ok := changes[filepath.Join(path, "terraform.tfvars")]
	assert.True(t, ok)
	assert.Nil(t, file)

	assert.NotNil(t, changes[filepath.Join(path, TfvarsAlternateFilename)])

}

func TestTfvarsStep_complete(t *testing.T) {
	parser := configs.NewParser(nil)
	path := "./fixtures/tfvars/complete"
	mod, _ := parser.LoadConfigDir(path)

	step := TerraformWorkspaceStep{
		Module: mod,
	}

	assert.True(t, step.Complete())
}
