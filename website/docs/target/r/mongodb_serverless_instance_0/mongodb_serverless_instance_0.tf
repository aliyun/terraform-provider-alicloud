data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_mongodb_serverless_instance" "example" {
  account_password        = "Abc12345"
  db_instance_description = "example_value"
  db_instance_storage     = 5
  storage_engine          = "WiredTiger"
  capacity_unit           = 100
  engine                  = "MongoDB"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  engine_version          = "4.2"
  period                  = 1
  period_price_type       = "Month"
  vpc_id                  = data.alicloud_vpcs.default.ids.0
  zone_id                 = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  tags = {
    Created = "MongodbServerlessInstance"
    For     = "TF"
  }
  security_ip_groups {
    security_ip_group_attribute = "example_value"
    security_ip_group_name      = "example_value"
    security_ip_list            = "192.168.0.1"
  }
}

