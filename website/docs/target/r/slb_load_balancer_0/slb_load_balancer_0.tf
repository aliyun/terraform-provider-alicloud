# Create a intranet SLB instance
variable "name" {
  default = "terraformtestslbconfig"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name = var.name
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  address_type       = "intranet"
  load_balancer_spec = "slb.s2.small"
  vswitch_id         = alicloud_vswitch.default.id
  tags = {
    info = "create for internet"
  }
}
