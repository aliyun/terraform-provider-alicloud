variable "name" {
  default = "natGatewayExampleName"
}

resource "alicloud_vpc" "enhanced" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

resource "alicloud_vswitch" "enhanced" {
  vswitch_name = var.name
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cidr_block   = "10.10.0.0/20"
  vpc_id       = alicloud_vpc.enhanced.id
}

resource "alicloud_nat_gateway" "enhanced" {
  depends_on       = [alicloud_vswitch.enhanced]
  vpc_id           = alicloud_vpc.enhanced.id
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = alicloud_vswitch.enhanced.id
  nat_type         = "Enhanced"
}
