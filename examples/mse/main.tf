terraform {
  required_providers {
    alicloud = {
      source = "registry.terraform.io/aliyun/alicloud"
    }
  }
}

data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_mse_cluster" "example" {
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_2_0_0"
  instance_count        = 3
  net_type              = "privatenet"
  pub_network_flow      = "1"
  connection_type       = "slb"
  cluster_alias_name    = "terraform-example"
  mse_version           = "mse_pro"
  vswitch_id            = alicloud_vswitch.example.id
  vpc_id                = alicloud_vpc.example.id
}

resource "alicloud_mse_engine_namespace" "example" {
  instance_id         = alicloud_mse_cluster.example.id
  namespace_show_name = "terraform-example"
  namespace_id        = "terraform-example"
  namespace_desc      = "description"
}

# Declare the data source
data "alicloud_mse_engine_namespaces" "example" {
  instance_id = alicloud_mse_engine_namespace.example.instance_id
}

output "mse_engine_namespace_id_public" {
  value = data.alicloud_mse_engine_namespaces.example.namespaces.0.id
}

output "mse_engine_namespace_id_example" {
  value = data.alicloud_mse_engine_namespaces.example.namespaces.1.id
}