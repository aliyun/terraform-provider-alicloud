// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name              = "${var.name}"
}
resource "alicloud_polardb_cluster" "default" {
  db_type              = "MySQL"
  db_version           = "8.0"
  db_node_class        = "polar.mysql.x4.large"
  pay_type             = "PostPaid"
  description          = "${var.name}"
  vswitch_id           = "${alicloud_vswitch.default.id}"
}
