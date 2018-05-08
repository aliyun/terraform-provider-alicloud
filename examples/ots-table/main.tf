provider "alicloud" {
  region = "cn-hangzhou"
  ots_instance_name = "${var.ots_instance_name}"
}

resource "alicloud_ots_table" "table" {
  provider = "alicloud"
  table_name = "${var.table_name}"
  primary_key = [
    {
      name = "${var.primary_key_1_name}"
      type = "${var.primary_key_integer_type}"
    },
    {
      name = "${var.primary_key_2_name}"
      type = "${var.primary_key_integer_type}"
    },
    {
      name = "${var.primary_key_3_name}"
      type = "${var.primary_key_integer_type}"
    },
    {
      name = "${var.primary_key_4_name}"
      type = "${var.primary_key_string_type}"
    },
  ]
  time_to_live = "${var.time_to_live}"
  max_version = "${var.max_version}"
}
