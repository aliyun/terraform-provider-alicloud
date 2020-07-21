resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  name = var.name
}
data "alicloud_zones" "default" {
}
resource "alicloud_vswitch" "default" {
  cidr_block        = "192.168.0.0/24"
  availability_zone = data.alicloud_zones.default.zones[0].id
  vpc_id            = alicloud_vpc.default.id
}
resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_network_interface" "default" {
  vswitch_id      = alicloud_vswitch.default.id
  security_groups = [ alicloud_security_group.default.id ]
  private_ip        = "192.168.0.2"
  private_ips_count = 3
}
resource "alicloud_log_project" "default"{
  name = lower(var.name)
  description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  project = alicloud_log_project.default.name
  name = lower(var.name)
  retention_period = 3650
  shard_count = 3
  auto_split = true
  max_split_shard_count = 60
  append_meta = true
}
resource "alicloud_vpc_flowlog" "default" {
  resource_id = alicloud_vpc.default.id
  resource_type ="VPC"
  traffic_type = "All"
  log_store_name = alicloud_log_store.default.name
  project_name = alicloud_log_project.default.name
  flow_log_name = var.name
  description = var.description
  status = "Inactive"
}

data "alicloud_vpc_flowlogs" "describe" {
  log_store_name = alicloud_vpc_flowlog.default.log_store_name
  project_name = alicloud_vpc_flowlog.default.project_name
  resource_id = alicloud_vpc_flowlog.default.resource_id
  resource_type = alicloud_vpc_flowlog.default.resource_type
  traffic_type = alicloud_vpc_flowlog.default.traffic_type
  status = "Inactive"
  description = alicloud_vpc_flowlog.default.description
  depends_on = [alicloud_vpc_flowlog.default]
}
