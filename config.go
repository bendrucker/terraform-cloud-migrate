package migrate

type Config struct {
	Backend           RemoteBackendConfig
	WorkspaceVariable string
	TfvarsFilename    string
}
