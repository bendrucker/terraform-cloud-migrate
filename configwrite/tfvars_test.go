package configwrite

import (
	"testing"
)

func TestTfvars(t *testing.T) {
	step := &Tfvars{
		Filename: "terraform.auto.tfvars",
	}

	testStepChanges(t, step, stepTests{
		{
			name: "incomplete",
			in: map[string]string{
				"main.tf": "",
				"terraform.tfvars": `
					foo = "bar"
					baz = "qux"
				`,
			},
			expected: map[string]string{
				"terraform.auto.tfvars": `
					foo = "bar"
					baz = "qux"
				`,
			},
			diags: nil,
		},
		{
			name: "complete",
			in: map[string]string{
				"main.tf": "",
				"terraform.auto.tfvars": `
					foo = "bar"
					baz = "qux"
				`,
			},
			expected: map[string]string{},
			diags:    nil,
		},
	})
}
