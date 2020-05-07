package migrate

import (
	"fmt"

	"github.com/bendrucker/terraform-cloud-migrate/configwrite"
	"github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/command/cliconfig"
)

type Config struct {
	API               *tfe.Config
	Backend           configwrite.RemoteBackendConfig
	WorkspaceVariable string
	TfvarsFilename    string
	ModulesDir        string
}

type RemoteBackendConfig = configwrite.RemoteBackendConfig
type WorkspaceConfig = configwrite.WorkspaceConfig

func LoadAPIConfig(host string) (*tfe.Config, hcl.Diagnostics) {
	// load environment variables
	config := tfe.DefaultConfig()
	config.Address = fmt.Sprintf("https://%s", host)

	if config.Token == "" {
		cfg, diags := cliconfig.LoadConfig()
		vDiags := cfg.Validate()
		diags = append(diags, vDiags...)

		if diags.HasErrors() {
			return config, hcl.Diagnostics{
				&hcl.Diagnostic{
					Summary: "Failed to load Terraform credentials",
					Detail:  diags.Err().Error(),
				},
			}
		}

		if api, ok := cfg.Credentials[host]; ok {
			if token, ok := api["token"]; ok {
				config.Token = token.(string)
			}
		}
	}

	return config, nil
}
