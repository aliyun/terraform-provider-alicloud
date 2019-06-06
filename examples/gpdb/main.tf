// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "Gpdb"
}

// VPC Resource for Module
resource "alicloud_vpc" "default" {
  name       = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "${var.vswitch_cidr}"
  availability_zone = "${var.availability_zone == "" ? data.alicloud_zones.default.zones.0.id : var.availability_zone}"
  name              = "${var.vswitch_name}"
}

resource "alicloud_gpdb_instance" "example" {
  description            = "terraform-xueqian-test"
  instance_class         = "${var.instance_class}"
  instance_group_count   = "${var.instance_group_count}"
  vswitch_id             = "${alicloud_vswitch.default.id}"
  engine                 = "gpdb"
  engine_version         = "4.3"
}
