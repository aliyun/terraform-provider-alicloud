provider "alicloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  # If not set, cn-beijing will be used.
  region = var.region
}
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
  network_type_acl = ["VPC", "INTERNET"]
  network_source_acl = ["TRUST_PROXY"]
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

  defined_column {
    name = var.defined_column_4_name
    type = var.boolean_type
  }

  defined_column {
    name = var.defined_column_5_name
    type = var.double_type
  }

  time_to_live = var.time_to_live
  max_version  = var.max_version
}

