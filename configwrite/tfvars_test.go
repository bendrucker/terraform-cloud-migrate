package configwrite

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTfvars_incomplete(t *testing.T) {
	writer := newTestModule(t, map[string]string{
		"main.tf": "",
		"terraform.tfvars": `
			foo = "bar"
			baz = "qux"
		`,
	})

	step := Tfvars{
		Writer:   writer,
		Filename: "terraform.auto.tfvars",
	}

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)
	assert.Len(t, changes, 1)

	change := changes["terraform.tfvars"]
	assert.Equal(t, "terraform.auto.tfvars", change.Rename)
}

func TestTfvars_complete(t *testing.T) {
	writer := newTestModule(t, map[string]string{
		"main.tf": "",
		"terraform.auto.tfvars": `
			foo = "bar"
			baz = "qux"
		`,
	})

	step := Tfvars{
		Writer: writer,
	}

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)
	assert.Len(t, changes, 0)
}
