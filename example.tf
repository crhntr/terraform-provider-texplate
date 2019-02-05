provider "texplate" {}

data "texplate_execute" "test" {
  template = "Hello, world!"
}

output "test" {
  value = "${data.texplate_execute.test.output}"
}
