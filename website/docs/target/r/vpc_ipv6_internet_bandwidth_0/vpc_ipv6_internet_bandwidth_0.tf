data "alicloud_instances" "example" {
  name_regex = "ecs_with_ipv6_address"
  status     = "Running"
}

data "alicloud_vpc_ipv6_addresses" "example" {
  associated_instance_id = data.alicloud_instances.example.instances.0.id
  status                 = "Available"
}

resource "alicloud_vpc_ipv6_internet_bandwidth" "example" {
  ipv6_address_id      = data.alicloud_vpc_ipv6_addresses.example.addresses.0.id
  ipv6_gateway_id      = data.alicloud_vpc_ipv6_addresses.example.addresses.0.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
}

