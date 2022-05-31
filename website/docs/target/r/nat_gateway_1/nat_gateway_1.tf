variable "name" {
  default = "nat-transform-to-enhanced"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

resource "alicloud_vpc" "foo" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "foo1" {
  depends_on   = [alicloud_vpc.foo]
  vswitch_name = var.name
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones[1].zone_id
  cidr_block   = "10.10.0.0/20"
  vpc_id       = alicloud_vpc.foo.id
}

resource "alicloud_nat_gateway" "main" {
  depends_on       = [alicloud_vpc.foo, alicloud_vswitch.foo1]
  vpc_id           = alicloud_vpc.foo.id
  nat_gateway_name = var.name
  nat_type         = "Enhanced"
  vswitch_id       = alicloud_vswitch.foo1.id
}
