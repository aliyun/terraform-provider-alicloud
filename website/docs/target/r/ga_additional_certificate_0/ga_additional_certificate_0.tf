variable "name" {
  default = "tf-testacc-ga"
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth       = 20
  type            = "Basic"
  bandwidth_type  = "Basic"
  duration        = 1
  ratio           = 30
  auto_pay        = true
  auto_use_coupon = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ssl_certificates_service_certificate" "default" {
  count            = 2
  certificate_name = var.name
  cert             = file("${path}/test.crt")
  key              = file("${path}/test.key")
}

resource "alicloud_ga_listener" "default" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.default]
  accelerator_id = alicloud_ga_accelerator.default.id
  name           = var.name
  protocol       = "HTTPS"
  port_ranges {
    from_port = 8080
    to_port   = 8080
  }
  certificates {
    id = join("-", [alicloud_ssl_certificates_service_certificate.default.0.id, "cn-hangzhou"])
  }
}

resource "alicloud_ga_additional_certificate" "default" {
  certificate_id = join("-", [alicloud_ssl_certificates_service_certificate.default.1.id, "cn-hangzhou"])
  domain         = "test"
  accelerator_id = alicloud_ga_accelerator.default.id
  listener_id    = alicloud_ga_listener.default.id
}
