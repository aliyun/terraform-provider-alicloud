provider "alicloud" {
  endpoints {
    bssopenapi = "business.aliyuncs.com"
  }
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.description}"
  cidr_block = "${var.vpc_cidr_block}"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "${var.vswitch_cidr_block}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.description}"
}

resource "alicloud_security_group" "default" {
  name   = "${var.security_name}"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_yundun_bastionhost_instance" "instance" {
  description       = "${var.description}"
  license_code      = "${var.license_code}"
  period            = "${var.period}"
  vswitch_id        = "${alicloud_vswitch.default.id}"
  security_group_ids = ["${alicloud_security_group.default.id}"]
}
