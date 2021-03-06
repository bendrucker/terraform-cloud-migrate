module github.com/bendrucker/terraform-cloud-migrate

go 1.14

require (
	github.com/hashicorp/hcl/v2 v2.3.0
	github.com/hashicorp/terraform v0.12.24
	github.com/lithammer/dedent v1.1.0
	github.com/mitchellh/cli v1.0.0
	github.com/spf13/afero v1.2.1
	github.com/spf13/pflag v1.0.3
	github.com/stretchr/testify v1.3.0
	github.com/zclconf/go-cty v1.2.1
)

replace github.com/hashicorp/hcl/v2 => github.com/bendrucker/hcl/v2 v2.4.1-0.20200429040843-8e720e092f94
