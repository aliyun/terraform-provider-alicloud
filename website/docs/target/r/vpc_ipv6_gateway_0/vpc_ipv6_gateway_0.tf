resource "alicloud_vpc" "default" {
  vpc_name    = "example_value"
  enable_ipv6 = "true"
}

resource "alicloud_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = "example_value"
  vpc_id            = alicloud_vpc.default.id
}

