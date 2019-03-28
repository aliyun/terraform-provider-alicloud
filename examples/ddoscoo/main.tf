resource "alicloud_ddoscoo_instance" "instance" {
  business_endpoint = "${var.business_endpoint}"
  name = "${var.name}"
  bandwidth = "${var.bandwidth}"
  base_bandwidth     = "${var.base_bandwidth}"
  service_bandwidth       = "${var.service_bandwidth}"
  port_count  = "${var.port_count}"
  domain_count  = "${var.domain_count}"
}
