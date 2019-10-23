// Zones data source for availability_zone
data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = var.vpc_name
  cidr_block = var.vpc_cidr
}

resource "alicloud_vswitch" "default" {
  name = var.vswitch_name
  vpc_id = alicloud_vpc.default.id
  cidr_block = var.vswitch_cidr
  availability_zone = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_serverless_kubernetes" "serverless" {
  name = var.serverless_cluster_name
  vpc_id = alicloud_vpc.default.id
  vswitch_id = alicloud_vswitch.default.id
  new_nat_gateway = true
  enndpoint_public_access_enabled = true
  private_zone = false
  deletion_protection = false
  tags = {
    "k-aa":"v-aa"
    "k-bb":"v-aa",
  }
}

