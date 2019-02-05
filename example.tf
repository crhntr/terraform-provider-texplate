provider "texplate" {}

data "texplate_template" "test" {
  template = "Hello, world!"
}

output "test" {
  value = "${data.texplate_template.test.output}"
}
