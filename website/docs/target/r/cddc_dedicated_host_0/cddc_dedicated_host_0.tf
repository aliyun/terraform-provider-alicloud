
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type      = "mysql"
  zone_id      = data.alicloud_cddc_zones.default.ids.0
  storage_type = "cloud_essd"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  engine                    = "MySQL"
  vpc_id                    = data.alicloud_vpcs.default.ids.0
  cpu_allocation_ratio      = 101
  mem_allocation_ratio      = 50
  disk_allocation_ratio     = 200
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = "example_value"
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = "example_value"
  dedicated_host_group_id = alicloud_cddc_dedicated_host_group.default.id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  payment_type            = "Subscription"
  tags = {
    Created = "TF"
    For     = "CDDC_DEDICATED"
  }
}

