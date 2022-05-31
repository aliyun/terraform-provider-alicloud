data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_mse_cluster" "default" {
  cluster_specification = "MSE_SC_1_2_200_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_ANS_1_2_1"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = data.alicloud_vswitches.default.ids.0
  pub_network_flow      = "1"
  acl_entry_list        = ["127.0.0.1/32"]
  cluster_alias_name    = "example_value"
}
resource "alicloud_mse_engine_namespace" "example" {
  cluster_id          = alicloud_mse_cluster.default.cluster_id
  namespace_show_name = "example_value"
  namespace_id        = "example_value"
}
