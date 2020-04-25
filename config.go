package migrate

import "github.com/bendrucker/terraform-cloud-migrate/steps"

type Config struct {
	Backend           steps.RemoteBackendConfig
	WorkspaceVariable string
	TfvarsFilename    string
	ModulesDir        string
}
