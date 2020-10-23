variable "name" {
  default = "natGatewayExampleName"
}

resource "alicloud_vpc" "enhanced" {
  name       = var.name
  cidr_block = "10.0.0.0/8"
}

data "alicloud_enhanced_nat_available_zones" "enhanced"{
}

resource "alicloud_vswitch" "enhanced" {
  name              = var.name
  availability_zone = data.alicloud_enhanced_nat_available_zones.enhanced.zones[0].zone_id
  cidr_block        = "10.10.0.0/20"
  vpc_id            = alicloud_vpc.enhanced.id
}

resource "alicloud_nat_gateway" "enhanced" {
  depends_on           = [alicloud_vswitch.enhanced]
  vpc_id               = alicloud_vpc.enhanced.id
  specification        = "Small"
  name                 = var.name
  instance_charge_type = "PostPaid"
  vswitch_id           = alicloud_vswitch.enhanced.id
  nat_type             = "Enhanced"
}
