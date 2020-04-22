package migrate

import (
	"github.com/hashicorp/hcl/v2"
	"io/ioutil"
	"github.com/zclconf/go-cty/cty"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
)

const (
	BackendTypeRemote = "remote"
)

type RemoteBackendStep struct {
	Module *configs.Module
	Config RemoteBackendConfig
}

type RemoteBackendConfig struct {
	Hostname string
	Organization string
	Workspaces WorkspaceConfig
}

type WorkspaceConfig struct {
	Name string
	Prefix string
}

// Complete checks if the module is using a remote backend
func (b *RemoteBackendStep) Complete() bool {
	return b.Module.Backend.Type == BackendTypeRemote
}

// Description returns a description of the step
func (b *RemoteBackendStep) Description() string {
	return `A "remote" backend should be configured for Terraform Cloud (https://www.terraform.io/docs/backends/types/remote.html)`
}

// MultipleWorkspaces returns whether the remote backend will be configured for multiple prefixed workspaces
func (b *RemoteBackendStep) MultipleWorkspaces() bool {
	return b.Config.Workspaces.Prefix != ""
}

// Changes updates the configured backend
func (b *RemoteBackendStep) Changes() (Changes, error) {
	path := b.Module.Backend.DeclRange.Filename
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	file, _ := hclwrite.ParseConfig(bytes, path, hcl.InitialPos)
	for _, block := range file.Body().Blocks() {
		if block.Type() != "terraform" {
			continue
		}

		for _, child := range block.Body().Blocks() {
			if child.Type() != "backend" {
				continue
			}

			block.Body().RemoveBlock(child)

			remote := block.Body().AppendBlock(hclwrite.NewBlock("backend", []string{"remote"})).Body()
			remote.SetAttributeValue("hostname", cty.StringVal(b.Config.Hostname))
			remote.SetAttributeValue("organization", cty.StringVal(b.Config.Organization))
			remote.AppendNewline()

			workspaces := remote.AppendBlock(hclwrite.NewBlock("workspaces", nil)).Body()
			if b.MultipleWorkspaces() {
				workspaces.SetAttributeValue("prefix", cty.StringVal(b.Config.Workspaces.Prefix))
			} else {
				workspaces.SetAttributeValue("name", cty.StringVal(b.Config.Workspaces.Name))
			}
		}

	}

	return Changes{path: file}, nil
}

var _ Step = (*RemoteBackendStep)(nil)