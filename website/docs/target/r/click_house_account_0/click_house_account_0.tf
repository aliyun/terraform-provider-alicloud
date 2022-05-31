variable "name" {
  default = "testaccountname"
}
variable "pwd" {
  default = "Tf-testpwd"
}
data "alicloud_click_house_regions" "default" {
  current = true
}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_click_house_regions.default.regions.0.zone_ids.0.zone_id
}
resource "alicloud_click_house_db_cluster" "default" {
  db_cluster_version      = "20.3.10.75"
  category                = "Basic"
  db_cluster_class        = "S8"
  db_cluster_network_type = "vpc"
  db_cluster_description  = var.name
  db_node_group_count     = 1
  payment_type            = "PayAsYouGo"
  db_node_storage         = "500"
  storage_type            = "cloud_essd"
  vswitch_id              = data.alicloud_vswitches.default.vswitches.0.id
}
resource "alicloud_click_house_account" "default" {
  db_cluster_id       = alicloud_click_house_db_cluster.default.id
  account_description = "your_description"
  account_name        = var.name
  account_password    = var.pwd
}
