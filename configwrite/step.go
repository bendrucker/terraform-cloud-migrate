package configwrite

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
)

// Step is a step required to prepare a module to run in Terraform Cloud
type Step interface {
	Name() string

	// Description returns a description of the step
	Description() string

	// Changes returns a list of files changes and diagnostics if errors ocurred. If Complete() returns true, this should be empty.
	Changes() (Changes, hcl.Diagnostics)
}

type Steps []Step

func (s Steps) Append(steps ...Step) Steps {
	return append(s, steps...)
}

func (s Steps) Changes() (Changes, hcl.Diagnostics) {
	changes := make(Changes)
	diags := hcl.Diagnostics{}

	for _, step := range s {
		stepChanges, sDiags := step.Changes()
		diags = append(diags, sDiags...)

		for path, change := range stepChanges {
			if err, ok := changes.Add(path, change).(*renameCollisionError); ok {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagWarning,
					Summary:  "Rename skipped due to conflict",
					Detail:   fmt.Sprintf(`The "%s" step attempted to rename %s to %s, but a previous step already renamed this file to %s.`, step.Name(), path, err.Proposed, err.Existing),
					Subject:  &hcl.Range{Filename: err.Proposed},
				})
			}
		}
	}

	return changes, diags
}
