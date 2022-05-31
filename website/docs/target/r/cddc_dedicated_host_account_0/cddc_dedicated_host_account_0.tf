variable "name" {
  default = "tftestacc"
}

data "alicloud_cddc_zones" "default" {}

data "alicloud_cddc_host_ecs_level_infos" "default" {
  db_type        = "mssql"
  zone_id        = data.alicloud_cddc_zones.default.ids.0
  storage_type   = "cloud_essd"
  image_category = "WindowsWithMssqlStdLicense"

}

data "alicloud_cddc_dedicated_host_groups" "default" {
  name_regex = "default-NODELETING"
  engine     = "mssql"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alicloud_cddc_dedicated_host_group" "default" {
  count                     = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? 0 : 1
  engine                    = "SQLServer"
  vpc_id                    = data.alicloud_vpcs.default.ids.0
  allocation_policy         = "Evenly"
  host_replace_policy       = "Manual"
  dedicated_host_group_desc = var.name
  open_permission           = true
}

data "alicloud_vswitches" "default" {
  vpc_id  = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.groups[0].vpc_id : data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_vswitch" "default" {
  count      = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id     = data.alicloud_vpcs.default.ids.0
  cidr_block = data.alicloud_vpcs.default.vpcs[0].cidr_block
  zone_id    = data.alicloud_cddc_zones.default.ids.0
}

resource "alicloud_cddc_dedicated_host" "default" {
  host_name               = var.name
  dedicated_host_group_id = length(data.alicloud_cddc_dedicated_host_groups.default.ids) > 0 ? data.alicloud_cddc_dedicated_host_groups.default.ids.0 : alicloud_cddc_dedicated_host_group.default[0].id
  host_class              = data.alicloud_cddc_host_ecs_level_infos.default.infos.0.res_class_code
  zone_id                 = data.alicloud_cddc_zones.default.ids.0
  vswitch_id              = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids.0 : alicloud_vswitch.default[0].id
  payment_type            = "Subscription"
  image_category          = "WindowsWithMssqlStdLicense"
}
resource "alicloud_cddc_dedicated_host_account" "example" {
  account_name      = var.name
  account_password  = "yourpassword"
  dedicated_host_id = alicloud_cddc_dedicated_host.default.dedicated_host_id
  account_type      = "Normal"
}
