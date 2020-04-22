package migrate

import (
	"strings"
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
		Path: "./fixtures/remote-state/incomplete",
	}

	changes, diags := step.Changes()
	assert.Len(t, diags, 0)
	assert.Len(t, changes, 1)

	expected := strings.TrimSpace(`
data "terraform_remote_state" "match" {
  backend = "remote"

  config = { hostname = "host.name", organization = "org", workspaces = { name = "ws" } }
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

	assert.Equal(t, expected+"\n", string(changes["./fixtures/remote-state/incomplete"].File.Bytes()))
}