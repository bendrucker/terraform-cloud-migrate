package steps

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteBackend_incomplete(t *testing.T) {
	path := "./fixtures/backend/incomplete"
	mod, diags := NewWriter(path)

	if diags.HasErrors() {
		assert.Fail(t, diags.Error())
	}

	step := RemoteBackend{
		Writer: mod,
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

func TestRemoteBackend_incomplete_prefix(t *testing.T) {
	path := "fixtures/backend/incomplete"
	mod, diags := NewWriter(path)

	if diags.HasErrors() {
		assert.Error(t, diags)
	}

	step := RemoteBackend{
		Writer: mod,
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

func TestRemoteBackend_complete(t *testing.T) {
	mod, diags := NewWriter("./fixtures/backend/complete")

	if diags.HasErrors() {
		assert.Error(t, diags)
	}

	step := RemoteBackend{
		Writer: mod,
	}

	assert.True(t, step.Complete())
}
