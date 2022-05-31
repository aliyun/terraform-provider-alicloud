resource "alicloud_vpc" "default" {
  vpc_name    = "example_value"
  enable_ipv6 = "true"
}

resource "alicloud_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = "example_value"
  vpc_id            = alicloud_vpc.default.id
}

data "alicloud_instances" "default" {
  name_regex = "ecs_with_ipv6_address"
  status     = "Running"
}

data "alicloud_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alicloud_instances.default.instances.0.id
  status                 = "Available"
}

resource "alicloud_vpc_ipv6_egress_rule" "example" {
  instance_id           = data.alicloud_vpc_ipv6_addresses.default.ids.0
  ipv6_egress_rule_name = "example_value"
  description           = "example_value"
  ipv6_gateway_id       = alicloud_vpc_ipv6_gateway.example.id
  instance_type         = "Ipv6Address"
}

