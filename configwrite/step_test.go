package configwrite

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/lithammer/dedent"
	"github.com/stretchr/testify/assert"
)

type stepTest struct {
	name     string
	in       map[string]string
	expected map[string]string
	diags    hcl.Diagnostics
}

type stepTests []stepTest

func testStepChanges(t *testing.T, step Step, tests stepTests) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := newTestModule(t, test.in)
			changes, diags := step.WithWriter(writer).Changes()

			out := make(map[string]string)
			for path, change := range changes {
				out[change.Destination(path)] = string(change.File.Bytes())
			}

			expected := make(map[string]string)
			for path, content := range test.expected {
				expected[path] = dedent.Dedent(content)
			}

			assert.Equal(t, expected, out)
			assert.Equal(t, test.diags, diags)
		})
	}
}
