# terraform-cloud-migrate [![tests](https://github.com/bendrucker/terraform-cloud-migrate/workflows/tests/badge.svg?branch=master)](https://github.com/bendrucker/terraform-cloud-migrate/actions?query=workflow%3Atests) [![Project Status: WIP](https://www.repostatus.org/badges/latest/wip.svg)](https://www.repostatus.org/#wip)

> Migrate a Terraform module to [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html)

The `terraform-cloud-migrate` CLI automates the process of adapting a Terraform [root module](https://www.terraform.io/docs/glossary.html#root-module) for [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html) (including Terraform Enterprise). It performs a series of required configuration changes (described below) and runs `terraform init` which will prompt you to copy state from your existing backend to workspaces in Terraform Cloud.

Versioning your Terraform configuration with `git` is **strongly** encouraged so you can recover in the event of unwanted changes.

## Installing

Binaries are available for each [tagged release](https://github.com/bendrucker/terraform-cloud-migrate/releases). Download an appropriate binary for your operating system and install it into `$PATH`.

## Usage

```
Usage: terraform-cloud-migrate [--version] [--help] <command> [<args>]

Available commands are:
    run    Run Terraform Cloud migration
```

### `run`

```
Usage: terraform-cloud-migrate run [DIR] [options]
  Migrate a Terraform module to Terraform Cloud

Options:
  -n, --workspace-name string       The name of the Terraform Cloud workspace (conflicts with --workspace-prefix)
  -p, --workspace-prefix string     The prefix of the Terraform Cloud workspaces (conflicts with --workspace-name)
  -m, --modules string              A directory where other Terraform modules are stored. If set, it will be scanned recursively for terrafor_remote_state references.
      --workspace-variable string   Variable that will replace terraform.workspace (default "environment")
      --tfvars-filename string      New filename for terraform.tfvars (default "terraform.auto.tfvars")
      --hostname string             Hostname for Terraform Cloud (default "app.terraform.io")
      --organization string         Organization name in Terraform Cloud
      --no-init                     Disable calling 'terraform init' before and after updating configuration to copy state.
```

The `run` command performs the following file updates and runs `terraform init` to trigger Terraform to copy state to the new

* Configures a remote backend. ([?](https://www.terraform.io/docs/cloud/migrate/index.html#step-5-edit-the-backend-configuration)).
* Updates any [`terraform_remote_state`](https://www.terraform.io/docs/providers/terraform/d/remote_state.html) data sources that match the previous backend configuration.
* Replaces `terraform.workspace` with a variable of your choice, `var.environment` by default. ([?](https://www.terraform.io/docs/state/workspaces.html#current-workspace-interpolation))
* Renames `terraform.tfvars` to a name of your choice, `terraform.auto.tfvars` by default. ([?](https://www.terraform.io/docs/cloud/workspaces/variables.html#terraform-variables))

#### Examples

##### Basic

```sh
terraform-cloud-migrate run --organization my-org --workspace-name my-ws ./path/to/module
```

##### Remote State

Updates `terraform_remote_state` data sources in `~/src/tf`:

```sh
terraform-cloud-migrate run --modules ~/src/tf # ...
```

##### Terraform Enterprise

By default, `terraform-cloud-migrate` connects to Terraform Cloud at `app.terraform.io`. Terraform Enterprise users can set a custom hostname:

```sh
terraform-cloud-migrate run --hostname terraform.enterprise.host # ...
```


## License

MIT © [Ben Drucker](http://bendrucker.me)
