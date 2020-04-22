package migrate

import (
	"fmt"
	"io/ioutil"
	"os"

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
				Summary:  "Not a module directory",
				Detail:   fmt.Sprintf("Directory %s does not contain Terraform configuration files.", path),
			},
		}
	}

	module, diags := parser.LoadConfigDir(path)

	return &Module{
		module: module,
		files:  make(map[string]*hclwrite.File),
	}, diags
}

// Module provides access to information about the Terraform module structure and the ability to update its files
type Module struct {
	module *configs.Module
	files  map[string]*hclwrite.File
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

// RemoteStateDataSources returns a list of remote state data sources defined for the module
func (m *Module) RemoteStateDataSources() []*configs.Resource {
	resources := make([]*configs.Resource, 0)

	for _, resource := range m.module.DataResources {
		if resource.Type == "terraform_remote_state" {
			resources = append(resources, resource)
		}
	}

	return resources
}

// File returns an existing file object or creates and caches one
func (m *Module) File(path string) (*hclwrite.File, hcl.Diagnostics) {
	file, ok := m.files[path]
	if ok {
		return file, hcl.Diagnostics{}
	}

	b, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, hcl.Diagnostics{
			&hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "file read error",
				Detail:   fmt.Sprintf("file %s could not be read: %v", path, err),
			},
		}
	}

	var diags hcl.Diagnostics
	if os.IsNotExist(err) {
		file = hclwrite.NewEmptyFile()
	} else {
		file, diags = hclwrite.ParseConfig(b, path, hcl.InitialPos)
	}

	if file != nil {
		m.files[path] = file
	}

	return file, diags
}
