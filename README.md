# terraform-provider-texplate

**Please do not use this. It was a fun experiment. Consider Terraform's [template\_file](https://www.terraform.io/docs/providers/template/d/file.html) instead.**

Go templates syntax in terraform.
Inspired by [texplate](https://github.com/pivotal-cf/texplate).

### Included Template Helper Functions
- [sprig](http://masterminds.github.io/sprig/)
- [cidrhost](https://www.terraform.io/docs/configuration/interpolation.html#cidrhost-iprange-hostnum-)

### Example
*Example "director-config" yaml copied from [pivotal-cf/terraforming-azure](https://github.com/pivotal-cf/terraforming-azure).*

```hcl
provider "texplate" {}

data "texplate_execute" "greeting" {
  template = "Hello, world!"
}

output "greeting" {
  value = "${data.texplate_execute.greeting.output}"
}

data "local_file" "director_config_template" {
  filename = "${path.module}/testdata/configure-director.yml"
}

data "texplate_execute" "director_config" {
  template = "${data.local_file.director_config_template.content}"

  vars {
    "subscription_id"               = "some-sub-id"
    "tenant_id"                     = "some-tenant-id"
    "client_id"                     = "some-client-id"
    "client_secret"                 = "some-client-secret"
    "pcf_resource_group_name"       = "floating-pods"
    "bosh_root_storage_account"     = "bosh-root-storage-account"
    "ops_manager_ssh_public_key"    = "some-ssh-public-key"
    "ops_manager_ssh_private_key"   = "-----BEGIN RSA PRIVATE KEY-----\n701\n-----END RSA PRIVATE KEY-----\n"
    "infrastructure_subnet_name"    = "floating-pods-infra-subnet"
    "infrastructure_subnet_cidr"    = "10.0.8.0/26"
    "infrastructure_subnet_gateway" = "10.0.8.1"
    "network_name"                  = "floating-pods-virtual-network"
    "control_plane_subnet_cidr"     = "10.0.10.0/28"
    "control_plane_subnet_name"     = "floating-pods-plane-subnet"
    "control_plane_subnet_gateway"  = "10.0.10.0"
  }
}

output "director_config" {
  value = "${data.texplate_execute.director_config.output}"
}
```
## Community

### Requesting Some Helper Function

Given this is an early experiment, please send me any helper functions you would like included.

Please PR a new `my_func.go` and `my_func_test.go` and include an example in the example.tf file using your function.
See `main.go` for how to add the function to the defaultTemplate config.
