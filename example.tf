provider "texplate" {}

data "texplate_execute" "greeting" {
  template = "Hello, world!"
}

output "greeting" {
  value = "${data.texplate_execute.greeting.output}"
}

data "local_file" "control_plane_config_template" {
  filename = "${path.module}/testdata/configure-control-plane.yml"
}

data "texplate_execute" "control_plane_config" {
  template = "${data.local_file.control_plane_config_template.content}"

  vars = {
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

output "control_plane_config" {
  value = "${data.texplate_execute.control_plane_config.output}"
}
