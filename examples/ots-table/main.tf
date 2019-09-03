data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = "example-ots-table"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  name              = "example-ots-table"
  cidr_block        = "172.16.1.0/24"
  vpc_id            = "${alicloud_vpc.default.id}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_ots_instance" "default" {
  name        = "${var.ots_instance_name}"
  description = "TF ots instance example"
}

resource "alicloud_ots_instance_attachment" "default" {
  instance_name = "${alicloud_ots_instance.default.id}"
  vswitch_id    = "${alicloud_vswitch.default.id}"
  vpc_name      = "table"
}

resource "alicloud_ots_table" "table" {
  instance_name = "${alicloud_ots_instance.default.name}"
  table_name    = "${var.table_name}"

  primary_key {
      name = "${var.primary_key_1_name}"
      type = "${var.primary_key_integer_type}"
  }
  primary_key {
      name = "${var.primary_key_2_name}"
      type = "${var.primary_key_integer_type}"
  }
  primary_key {
      name = "${var.primary_key_3_name}"
      type = "${var.primary_key_integer_type}"
  }
  primary_key {
      name = "${var.primary_key_4_name}"
      type = "${var.primary_key_string_type}"
  }

  time_to_live = "${var.time_to_live}"
  max_version  = "${var.max_version}"
}
