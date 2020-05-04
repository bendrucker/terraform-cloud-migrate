package migrate

import "github.com/bendrucker/terraform-cloud-migrate/configwrite"

type Config struct {
	Backend           configwrite.RemoteBackendConfig
	WorkspaceVariable string
	TfvarsFilename    string
	ModulesDir        string
}
