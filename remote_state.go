package migrate

import (
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/hashicorp/terraform/configs"
	"github.com/zclconf/go-cty/cty"
)

type RemoteStateStep struct {
	module        *Module
	Path          string
	RemoteBackend RemoteBackendConfig
}

// Complete checks if any modules in the path are using remote_state
func (s *RemoteStateStep) Complete() bool {
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
			filepath := source.DeclRange.Filename
			file, fDiags := s.module.File(source.DeclRange.Filename)
			diags = append(diags, fDiags...)

			block := file.Body().FirstMatchingBlock("data", []string{
				source.Type,
				source.Name,
			})

			workspace := block.Body().RemoveAttribute("workspace")

			block.Body().SetAttributeValue("backend", cty.StringVal("remote"))
			block.Body().SetAttributeRaw("config", flattenTokens([]hclwrite.Tokens{
				{
					{
						Type:  hclsyntax.TokenOBrace,
						Bytes: []byte("{"),
					},
					{
						Type:  hclsyntax.TokenNewline,
						Bytes: []byte("\n"),
					},
					{
						Type:  hclsyntax.TokenIdent,
						Bytes: []byte("hostname"),
					},
					{
						Type:  hclsyntax.TokenEqual,
						Bytes: []byte("="),
					},
					{
						Type:  hclsyntax.TokenOQuote,
						Bytes: []byte(`"`),
					},
					{
						Type:  hclsyntax.TokenQuotedLit,
						Bytes: []byte(s.RemoteBackend.Hostname),
					},
					{
						Type:  hclsyntax.TokenCQuote,
						Bytes: []byte(`"`),
					},
					{
						Type:  hclsyntax.TokenNewline,
						Bytes: []byte("\n"),
					},
					{
						Type:  hclsyntax.TokenIdent,
						Bytes: []byte("organization"),
					},
					{
						Type:  hclsyntax.TokenEqual,
						Bytes: []byte("="),
					},
					{
						Type:  hclsyntax.TokenOQuote,
						Bytes: []byte(`"`),
					},
					{
						Type:  hclsyntax.TokenQuotedLit,
						Bytes: []byte(s.RemoteBackend.Organization),
					},
					{
						Type:  hclsyntax.TokenCQuote,
						Bytes: []byte(`"`),
					},
					{
						Type:  hclsyntax.TokenNewline,
						Bytes: []byte("\n\n"),
					},
					{
						Type:  hclsyntax.TokenIdent,
						Bytes: []byte("workspaces"),
					},
					{
						Type:  hclsyntax.TokenEqual,
						Bytes: []byte("="),
					},
					{
						Type:  hclsyntax.TokenOBrace,
						Bytes: []byte("{"),
					},
					{
						Type:  hclsyntax.TokenNewline,
						Bytes: []byte("\n"),
					},
					{
						Type:  hclsyntax.TokenStringLit,
						Bytes: []byte("name"),
					},
					{
						Type:  hclsyntax.TokenEqual,
						Bytes: []byte("="),
					},
				},
				s.workspaceNameTokens(workspace),
				{
					{
						Type:  hclsyntax.TokenNewline,
						Bytes: []byte("\n"),
					},
					{
						Type:  hclsyntax.TokenCBrace,
						Bytes: []byte("}"),
					},
					{
						Type:  hclsyntax.TokenNewline,
						Bytes: []byte("\n"),
					},
					{
						Type:  hclsyntax.TokenCBrace,
						Bytes: []byte("}"),
					},
				},
			}))

			changes[filepath] = &Change{File: file}
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

		backend, bDiags := attrs["backend"].Expr.Value(nil)
		diags = append(diags, bDiags...)

		if backend.AsString() != s.module.Backend().Type {
			continue
		}

		config, cDiags := attrs["config"].Expr.Value(nil)
		diags = append(diags, cDiags...)

		remoteBackendConfigAttrs, rDiags := s.module.Backend().Config.JustAttributes()
		diags = append(diags, rDiags...)

		for key, value := range config.AsValueMap() {
			rbValue, rDiags := remoteBackendConfigAttrs[key].Expr.Value(nil)
			diags = append(diags, rDiags...)

			if value.AsString() != rbValue.AsString() {
				continue Source
			}
		}

		sources = append(sources, source)
	}

	return sources, diags
}

func (s *RemoteStateStep) workspaceNameTokens(workspace *hclwrite.Attribute) hclwrite.Tokens {
	if s.RemoteBackend.Workspaces.Prefix == "" {
		return hclwrite.Tokens{
			{
				Type:  hclsyntax.TokenOQuote,
				Bytes: []byte(`"`),
			},
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte(s.RemoteBackend.Workspaces.Name),
			},
			{
				Type:  hclsyntax.TokenCQuote,
				Bytes: []byte(`"`),
			},
		}
	}

	return flattenTokens([]hclwrite.Tokens{
		{
			{
				Type:  hclsyntax.TokenOQuote,
				Bytes: []byte(`"`),
			},
			{
				Type:  hclsyntax.TokenStringLit,
				Bytes: []byte(s.RemoteBackend.Workspaces.Prefix),
			},
			{
				Type:  hclsyntax.TokenTemplateInterp,
				Bytes: []byte("${"),
			},
		},
		workspace.Expr().BuildTokens(nil),
		{
			{
				Type:  hclsyntax.TokenTemplateSeqEnd,
				Bytes: []byte("}"),
			},
			{
				Type:  hclsyntax.TokenCQuote,
				Bytes: []byte(`"`),
			},
		},
	})
}

func flattenTokens(in []hclwrite.Tokens) hclwrite.Tokens {
	out := make(hclwrite.Tokens, 0)
	for _, tokens := range in {
		out = append(out, tokens...)
	}
	return out
}

var _ Step = (*RemoteStateStep)(nil)
