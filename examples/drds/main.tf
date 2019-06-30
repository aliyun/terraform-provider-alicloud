data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "foo" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "foo" {
  vpc_id            = "${alicloud_vpc.foo.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alicloud_drds_instance" "vpc" {
  description          = "drds vpc"
  zone_id              = "${data.alicloud_zones.default.zones.0.id}"
  instance_series      = "${var.instance_series}"
  instance_charge_type = "PostPaid"
  vswitch_id           = "${alicloud_vswitch.foo.id}"
  specification        = "drds.sn1.4c8g.8C16G"
}
