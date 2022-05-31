data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "example_value"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/21"
  zone_id      = data.alicloud_zones.example.zones.0.id
  vswitch_name = "example_value"
}

resource "alicloud_nat_gateway" "example" {
  vpc_id               = alicloud_vpc.example.id
  internet_charge_type = "PayByLcu"
  nat_gateway_name     = "example_value"
  description          = "example_value"
  nat_type             = "Enhanced"
  vswitch_id           = alicloud_vswitch.example.id
  network_type         = "intranet"
}

resource "alicloud_vpc_nat_ip_cidr" "example" {
  nat_ip_cidr             = "192.168.0.0/16"
  nat_gateway_id          = alicloud_nat_gateway.example.id
  nat_ip_cidr_description = "example_value"
  nat_ip_cidr_name        = "example_value"
}

resource "alicloud_vpc_nat_ip" "example" {
  nat_ip             = "192.168.0.37"
  nat_gateway_id     = alicloud_nat_gateway.example.id
  nat_ip_description = "example_value"
  nat_ip_name        = "example_value"
  nat_ip_cidr        = alicloud_vpc_nat_ip_cidr.example.nat_ip_cidr
}

