provider "alicloud" {
  endpoints {
    bssopenapi = "business.aliyuncs.com"
  }
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = var.description
  cidr_block = var.vpc_cidr_block
}

resource "alicloud_vswitch" "default" {
  vpc_id = alicloud_vpc.default.id
  cidr_block = var.vswitch_cidr_block
  availability_zone = data.alicloud_zones.default.zones.0.id
  name = var.description
}

resource "alicloud_yundun_dbaudit_instance" "instance" {
  description       = var.description
  plan_code         = var.plan_code
  period            = var.period
  vswitch_id        = alicloud_vswitch.default.id
}
