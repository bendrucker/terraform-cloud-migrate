# terraform-cloud-migrate

**Status**: WIP

Migrates a Terraform module to [Terraform Cloud](https://www.terraform.io/docs/cloud/index.html). Automatically detects and optionally writes required changes:

* Configures the ["remote" backend](https://www.terraform.io/docs/backends/types/remote.html)
* Replaces `terraform.workspace` interpolations with a variable of your choice (default: `var.environment`)
* Moves `terraform.tfvars` to an alternate location

## Usage

```sh
terraform-cloud-migrate --organization takescoop --workspace-name ws ./my/tf/module
```

By default, each proposed file change will be printed with the file name and proposed content. To write these changes, add `--write`. 
