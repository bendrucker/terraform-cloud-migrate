package migrate

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
)

func NewModule(path string) (*Module, hcl.Diagnostics) {
	parser := configs.NewParser(nil)

	if !parser.IsConfigDir(path) {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary: "Not a module directory",
				Detail: fmt.Sprintf("Directory %s does not contain Terraform configuration files.", path),
			},
		}
	}

	module, diags := parser.LoadConfigDir(path)
	primary, _, fDiags := parser.ConfigDirFiles(path)
	diags = append(diags, fDiags...)

	files := make(map[string]*hclwrite.File, len(primary))
	for _, filename := range primary {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			diags = diags.Append(&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary: "file read error",
				Detail: fmt.Sprintf("file %s could not be read: %v", filename, err),
			})
		}
		file, fDiag := hclwrite.ParseConfig(b, filename, hcl.InitialPos)
		diags = append(diags, fDiag...)
		files[filename] = file
	}

	return &Module{
		module: module,
		files: files,
	}, nil
}

// Module provides access to information about the Terraform module structure and the ability to update its files
type Module struct {
	module *configs.Module
	files map[string]*hclwrite.File
}

// Dir returns the module directory
func (m *Module) Dir() string {
	return m.module.SourceDir
}

// Backend returns the backend, or nil if none is defined
func (m *Module) Backend() *configs.Backend {
	return m.module.Backend
}

// HasBackend returns true if the module has a backend configuration
func (m *Module) HasBackend() bool {
	return m.Backend() != nil
}

// Variables returns the declared variables for the module
func (m *Module) Variables() map[string]*configs.Variable {
	return m.module.Variables
}

// File returns an existing file object or creates and caches one
func (m *Module) File(path string) *hclwrite.File {
	file, ok := m.files[path]
	if !ok {
		file = hclwrite.NewEmptyFile()
		m.files[path] = file
		return m.File(path)
	}

	return file
}