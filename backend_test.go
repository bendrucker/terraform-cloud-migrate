package migrate

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform/configs"
)

func TestRemoteBackendStep_incomplete(t *testing.T) {
	parser := configs.NewParser(nil)
	path := "./fixtures/backend/incomplete"
	mod, _ := parser.LoadConfigDir(path)

	step := RemoteBackendStep{
		Module: mod,
		Config: RemoteBackendConfig{
			Hostname:     "host.name",
			Organization: "org",
			Workspaces: WorkspaceConfig{
				Name: "ws",
			},
		},
	}

	assert.False(t, step.Complete())

	changes, err := step.Changes()
	if !assert.NoError(t, err) {
		return
	}

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

	assert.Equal(t, expected, string(changes[filepath.Join(path, "backend.tf")].Bytes()))
}

func TestRemoteBackendStep_incomplete_prefix(t *testing.T) {
	parser := configs.NewParser(nil)
	path := "fixtures/backend/incomplete"
	mod, _ := parser.LoadConfigDir(path)

	step := RemoteBackendStep{
		Module: mod,
		Config: RemoteBackendConfig{
			Hostname:     "host.name",
			Organization: "org",
			Workspaces: WorkspaceConfig{
				Prefix: "ws-",
			},
		},
	}

	assert.False(t, step.Complete())

	changes, err := step.Changes()
	if !assert.NoError(t, err) {
		return
	}

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

	assert.Equal(t, expected, string(changes[filepath.Join(path, "backend.tf")].Bytes()))
}

func TestRemoteBackendStep_complete(t *testing.T) {
	parser := configs.NewParser(nil)
	mod, _ := parser.LoadConfigDir("./fixtures/backend/complete")

	step := RemoteBackendStep{
		Module: mod,
	}

	assert.True(t, step.Complete())
}
