package migrate

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

type RemoteStateStep struct {
	module        *Module
	Path          string
	RemoteBackend RemoteBackendConfig
}

// Complete checks if any modules in the path are using remote_state
func (b *RemoteStateStep) Complete() bool {
	return false
}

// Description returns a description of the step
func (s *RemoteStateStep) Description() string {
	return `A "remote" backend should be configured for Terraform Cloud (https://www.terraform.io/docs/backends/types/remote.html)`
}

// Changes updates the configured backend
func (s *RemoteStateStep) Changes() (Changes, hcl.Diagnostics) {
	parser := configs.NewParser(nil)
	changes := Changes{}
	diags := hcl.Diagnostics{}

	_ = filepath.Walk(s.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || !parser.IsConfigDir(path) {
			return nil
		}

		sources, sDiags := s.sources(path)
		diags = append(diags, sDiags...)

		for _, source := range sources {
			file, fDiags := s.module.File(source.DeclRange.Filename)
			diags = append(diags, fDiags...)

			// sourceAttrs, sDiags := source.Config.JustAttributes()
			// diags = append(diags, sDiags...)

			// if ws, ok := sourceAttrs["workspace"]; ok {
			// 	workspace := ws.Expr
			// 	diags = append(diags, wDiags...)
			// }

			block := file.Body().FirstMatchingBlock("data", []string{
				source.Type,
				source.Name,
			})

			block.Body().SetAttributeValue("backend", cty.StringVal("remote"))
			block.Body().SetAttributeValue("config", cty.ObjectVal(map[string]cty.Value{
				"hostname":     cty.StringVal(s.RemoteBackend.Hostname),
				"organization": cty.StringVal(s.RemoteBackend.Organization),
				"workspaces": cty.MapVal(map[string]cty.Value{
					"name": cty.StringVal(s.RemoteBackend.Workspaces.Name),
				}),
			}))

			changes[path] = &Change{File: file}
		}

		if diags.HasErrors() {
			return diags
		}

		return nil
	})

	return changes, diags
}

// Changes updates the configured backend
func (s *RemoteStateStep) sources(path string) ([]*configs.Resource, hcl.Diagnostics) {
	mod, diags := NewModule(path)
	sources := make([]*configs.Resource, 0)

Source:
	for _, source := range mod.RemoteStateDataSources() {
		attrs, aDiags := source.Config.JustAttributes()
		diags = append(diags, aDiags...)

		for _, attr := range attrs {
			switch attr.Name {
			case "backend":
				v, vDiags := attr.Expr.Value(nil)
				diags = append(diags, vDiags...)

				if v.AsString() != s.module.Backend().Type {
					continue Source
				}
			case "config":
				remoteStateConfig, rDiags := attr.Expr.Value(nil)
				diags = append(diags, rDiags...)

				remoteBackendConfigAttrs, rDiags := s.module.Backend().Config.JustAttributes()
				diags = append(diags, rDiags...)

				for key, value := range remoteStateConfig.AsValueMap() {
					rbValue, rDiags := remoteBackendConfigAttrs[key].Expr.Value(nil)
					diags = append(diags, rDiags...)

					if value.AsString() != rbValue.AsString() {
						continue Source
					}
				}
			}
		}

		sources = append(sources, source)
	}

	return sources, diags
}

var _ Step = (*RemoteStateStep)(nil)
