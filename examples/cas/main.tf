provider "alicloud" {
	region = "cn-hangzhou"
}

data "alicloud_cas_certificates" "certs" {
  output_file = "./tmp.txt"
}

resource "alicloud_cas_certificate" "cert" {
   name = "test"
   cert = "./test.crt"
   key = "./test.key"
}
