package migrate

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteBackendStep_incomplete(t *testing.T) {
	path := "./fixtures/backend/incomplete"
	mod, diags := NewModule(path)

	if diags.HasErrors() {
		assert.Fail(t, diags.Error())
	}

	step := RemoteBackendStep{
		module: mod,
		Config: RemoteBackendConfig{
			Hostname:     "host.name",
			Organization: "org",
			Workspaces: WorkspaceConfig{
				Name: "ws",
			},
		},
	}

	assert.False(t, step.Complete())

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)

	assert.Len(t, changes, 1)

	expected := strings.TrimSpace(`
terraform {
  backend "remote" {
    hostname     = "host.name"
    organization = "org"

    workspaces {
      name = "ws"
    }
  }
}
`)

	assert.Equal(t, expected+"\n", string(changes[filepath.Join(path, "backend.tf")].File.Bytes()))
}

func TestRemoteBackendStep_incomplete_prefix(t *testing.T) {
	path := "fixtures/backend/incomplete"
	mod, diags := NewModule(path)

	if diags.HasErrors() {
		assert.Error(t, diags)
	}

	step := RemoteBackendStep{
		module: mod,
		Config: RemoteBackendConfig{
			Hostname:     "host.name",
			Organization: "org",
			Workspaces: WorkspaceConfig{
				Prefix: "ws-",
			},
		},
	}

	assert.False(t, step.Complete())

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)

	assert.Len(t, changes, 1)

	expected := strings.TrimSpace(`
terraform {
  backend "remote" {
    hostname     = "host.name"
    organization = "org"

    workspaces {
      prefix = "ws-"
    }
  }
}
`)

	assert.Equal(t, expected+"\n", string(changes[filepath.Join(path, "backend.tf")].File.Bytes()))
}

func TestRemoteBackendStep_complete(t *testing.T) {
	mod, diags := NewModule("./fixtures/backend/complete")

	if diags.HasErrors() {
		assert.Error(t, diags)
	}

	step := RemoteBackendStep{
		module: mod,
	}

	assert.True(t, step.Complete())
}
