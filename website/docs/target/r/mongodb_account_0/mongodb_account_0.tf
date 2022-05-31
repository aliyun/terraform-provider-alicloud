data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_name = "subnet-for-local-test"
}
resource "alicloud_mongodb_instance" "default" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  name                = var.name
  vswitch_id          = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

resource "alicloud_mongodb_account" "example" {
  account_name        = "root"
  account_password    = "example_value"
  instance_id         = alicloud_mongodb_instance.default.id
  account_description = "example_value"
}

