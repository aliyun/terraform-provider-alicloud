data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "example-ots-table"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = "example-ots-table"
  cidr_block   = "172.16.1.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_ots_instance" "default" {
  name        = var.ots_instance_name
  description = "TF ots instance example"
}

resource "alicloud_ots_instance_attachment" "default" {
  instance_name = alicloud_ots_instance.default.id
  vswitch_id    = alicloud_vswitch.default.id
  vpc_name      = "table"
}

resource "alicloud_ots_table" "table" {
  instance_name = alicloud_ots_instance.default.name
  table_name    = var.table_name

  primary_key {
    name = var.primary_key_1_name
    type = var.integer_type
  }
  primary_key {
    name = var.primary_key_2_name
    type = var.integer_type
  }
  primary_key {
    name = var.primary_key_3_name
    type = var.integer_type
  }
  primary_key {
    name = var.primary_key_4_name
    type = var.string_type
  }

  defined_column {
    name = var.defined_column_1_name
    type = var.integer_type
  }

  defined_column {
    name = var.defined_column_2_name
    type = var.string_type
  }

  defined_column {
    name = var.defined_column_3_name
    type = var.binary_type
  }

  time_to_live = var.time_to_live
  max_version  = var.max_version
}

resource "alicloud_ots_secondary_index" "index1" {
  instance_name = alicloud_ots_instance.default.name
  table_name = alicloud_ots_table.table.table_name

  # required
  index_name = var.secondary_index_name
  # required [Global|Local]
  index_type = var.secondary_index_type
  # required
  include_base_data = var.secondary_index_include_base_data
  # required
  primary_keys = var.secondary_index_pks
  # optional
  defined_columns = var.index_defined_cols
}