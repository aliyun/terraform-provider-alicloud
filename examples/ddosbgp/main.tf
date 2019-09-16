provider "alicloud" {
  endpoints {
    bssopenapi = "business.aliyuncs.com"
  }
}

resource "alicloud_ddosbgp_instance" "instance" {
  name              = "${var.name}"
  region            = "${var.region}"
  bandwidth         = "${var.bandwidth}"
  base_bandwidth    = "${var.base_bandwidth}"
  ip_count          = "${var.ip_count}"
  ip_type           = "${var.ip_type}"
}
