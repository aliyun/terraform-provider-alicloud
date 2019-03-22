// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "MongoDB"
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

resource "alicloud_mongodb_instance" "example" {
  engine_version      = "${var.engine_version}"
  db_instance_class   = "${var.instance_type}"
  db_instance_storage = "${var.storage}"
  name                = "tf-mongodb_instance-example"
  vswitch_id          = "${alicloud_vswitch.default.id}"
}
