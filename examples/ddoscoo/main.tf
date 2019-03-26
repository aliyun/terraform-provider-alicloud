resource "alicloud_ddoscoo_instance" "instance" {
  business_endpoint = "${var.business_endpoint}"
  band_width = "${var.band_width}"
  base_band_width     = "${var.base_band_width}"
  service_band_width       = "${var.service_band_width}"
  port_count  = "${var.port_count}"
  domain_count  = "${var.domain_count}"
}
