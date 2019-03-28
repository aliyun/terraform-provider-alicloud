resource "alicloud_ddoscoo_instance" "instance" {
  bssopenapi_endpoint = "${var.bssopenapi_endpoint}"
  name = "${var.name}"
  bandwidth = "${var.bandwidth}"
  base_bandwidth     = "${var.base_bandwidth}"
  service_bandwidth       = "${var.service_bandwidth}"
  port_count  = "${var.port_count}"
  domain_count  = "${var.domain_count}"
}
