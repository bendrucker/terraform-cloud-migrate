# terraform-cloud-migrate

> Migrate a Terraform module to [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html)

**Status**: WIP

The `terraform-cloud-migrate` CLI automates the process of adapting a Terraform [root module](https://www.terraform.io/docs/glossary.html#root-module) for [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html) (including Terraform Enterprise). It performs a series of required configuration changes (described below) and runs `terraform init` which will prompt you to copy state from your existing backend to workspaces in Terraform Cloud.

## Usage

```sh
terraform-cloud-migrate ./path/to/module --organization my-org --workspace-name my-ws ~/src/tf/my-module

# update terraform_remote_state with --modules
terraform-cloud-migrate ./path/to/module --organization my-org --workspace-name my-ws --modules ~/src/tf/my-module
```

## Steps

Steps include a link to the Terraform docs where available.

* Configures a remote backend. ([?](https://www.terraform.io/docs/cloud/migrate/index.html#step-5-edit-the-backend-configuration)).
* Updates any [`terraform_remote_state`](https://www.terraform.io/docs/providers/terraform/d/remote_state.html) data sources that match the previous backend configuration.
* Replaces `terraform.workspace` with a variable of your choice, `var.environment` by default. ([?](https://www.terraform.io/docs/state/workspaces.html#current-workspace-interpolation))
* Renames `terraform.tfvars` to a name of your choice, `default.auto.tfvars` by default. ([?](https://www.terraform.io/docs/cloud/workspaces/variables.html#terraform-variables))

## License

MIT © [Ben Drucker](http://bendrucker.me)
