resource "alicloud_mse_cluster" "example" {
  cluster_specification = "MSE_SC_1_2_200_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_ANS_1_2_1"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = "vsw-123456"
  pub_network_flow      = "1"
  acl_entry_list        = ["127.0.0.1/32"]
  cluster_alias_name    = "tf-testAccMseCluster"
}
