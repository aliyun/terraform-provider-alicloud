variable "name" {
  default = "forward-entry-example-name"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_nat_gateway" "default" {
  vpc_id        = alicloud_vpc.default.id
  specification = "Small"
  name          = var.name
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

resource "alicloud_eip_association" "default" {
  allocation_id = alicloud_eip_address.default.id
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_forward_entry" "default" {
  forward_table_id = alicloud_nat_gateway.default.forward_table_ids
  external_ip      = alicloud_eip_address.default.ip_address
  external_port    = "80"
  ip_protocol      = "tcp"
  internal_ip      = "172.16.0.3"
  internal_port    = "8080"
}
