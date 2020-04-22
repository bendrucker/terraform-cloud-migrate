package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteStateStep_incomplete(t *testing.T) {
	path := "./fixtures/backend/incomplete"
	mod, diags := NewModule(path)
	if diags.HasErrors() {
		assert.Error(t, diags)
	}

	step := RemoteStateStep{
		module: mod,
		RemoteBackend: RemoteBackendConfig{
			Hostname:     "host.name",
			Organization: "org",
			Workspaces: WorkspaceConfig{
				Name: "ws",
			},
		},
		Path: "./fixtures/remote-state",
	}

	_, diags = step.Changes()
	assert.Len(t, diags, 0)

	// assert.Equal(t, expected, string(changes[filepath.Join(path, "backend.tf")].Bytes()))
}
