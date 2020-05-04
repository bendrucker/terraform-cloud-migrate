package configwrite

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteState_incomplete(t *testing.T) {
	path := "./fixtures/backend/incomplete"
	mod, diags := New(path)
	if diags.HasErrors() {
		assert.Error(t, diags)
	}

	step := RemoteState{
		writer: mod,
		RemoteBackend: RemoteBackendConfig{
			Hostname:     "host.name",
			Organization: "org",
			Workspaces: WorkspaceConfig{
				Name: "ws",
			},
		},
		Path: "./fixtures/remote-state/incomplete",
	}

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)
	assert.Len(t, changes, 1)

	expected := strings.TrimSpace(`
data "terraform_remote_state" "match" {
  backend = "remote"

  config = {
    hostname     = "host.name"
    organization = "org"

    workspaces = {
      name = "ws"
    }
  }
}

data "terraform_remote_state" "wrong_type" {
  backend = "remote"

  config = {}
}

data "terraform_remote_state" "wrong_config" {
  backend = "s3"

  config = {
    key    = "a-different-terraform.tfstate"
    bucket = "terraform-state"
    region = "us-east-1"
  }
}
`)

	assert.Equal(t, expected+"\n", string(changes["fixtures/remote-state/incomplete/main.tf"].File.Bytes()))
}

func TestRemoteState_incomplete_prefix(t *testing.T) {
	path := "./fixtures/backend/incomplete"
	mod, diags := New(path)
	if diags.HasErrors() {
		assert.Error(t, diags)
	}

	step := RemoteState{
		writer: mod,
		RemoteBackend: RemoteBackendConfig{
			Hostname:     "host.name",
			Organization: "org",
			Workspaces: WorkspaceConfig{
				Prefix: "ws-",
			},
		},
		Path: "./fixtures/remote-state/incomplete-prefix",
	}

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)
	assert.Len(t, changes, 1)

	expected := strings.TrimSpace(`
data "terraform_remote_state" "match" {
  backend = "remote"

  config = {
    hostname     = "host.name"
    organization = "org"

    workspaces = {
      name = "ws-${terraform.workspace}"
    }
  }
}

data "terraform_remote_state" "wrong_type" {
  backend = "remote"

  config = {}
}

data "terraform_remote_state" "wrong_config" {
  backend = "s3"

  config = {
    key    = "a-different-terraform.tfstate"
    bucket = "terraform-state"
    region = "us-east-1"
  }
}
`)

	assert.Equal(t, expected+"\n", string(changes["fixtures/remote-state/incomplete-prefix/main.tf"].File.Bytes()))
}
