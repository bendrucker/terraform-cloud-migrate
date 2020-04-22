package migrate

import "github.com/hashicorp/hcl/v2/hclwrite"

// Step is a step required to prepare a module to run in Terraform Cloud
type Step interface {
	// Complete returns whether the step has been completed
	Complete() bool

	// Description returns a description of the step
	Description() string

	// Changes returns a list of files changes that will complete the step or any error if one ocurred. If Complete() returns true, this should be empty.
	Changes() ([]*hclwrite.File, error)
}
