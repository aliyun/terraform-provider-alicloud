provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_cas_certificates" "certs" {
  output_file = "${path.module}/cas_certificates.json"
}

resource "alicloud_cas_certificate" "cert" {
  name = "test"
  cert = file("${path.module}/test.crt")
  key  = file("${path.module}/test.key")
}

