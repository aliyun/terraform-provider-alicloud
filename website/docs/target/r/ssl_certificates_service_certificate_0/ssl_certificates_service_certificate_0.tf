resource "alicloud_ssl_certificates_service_certificate" "example" {
  certificate_name = "test"
  cert             = file("${path.module}/test.crt")
  key              = file("${path.module}/test.key")
}

