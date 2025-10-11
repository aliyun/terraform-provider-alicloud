// Zones data source for availability_zone
terraform {
  required_providers {
    alicloud = {
      source = "registry.terraform.io/aliyun/alicloud"
    }
  }
}

resource "alicloud_polardb_on_ens_cluster" "mytest" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.medium.c"
  pay_type      = "PrePaid"
  renewal_status = "AutoRenewal"
  period = 1
  auto_renew_period = 1
  db_minor_version = "8.0.2"
  description   = var.name
  ens_region_id = "sg-singapore-9"
  db_node_num = length(var.db_cluster_nodes_configs) >0 ? null : var.db_node_num
  storage_space = 20
  storage_type = "ESSDPL0"
  target_minor_version= "innovate_x86#20250311"
  vpc_id = "n-56tme6lvq1n6s3d1zn78zd12w"
  vswitch_id    = "vsw-56tme6tk908wmeco8kg3nse60"
  db_cluster_nodes_configs = {
    for node, config in var.db_cluster_nodes_configs : node => jsonencode({for k, v in config : k => v if v != null})
  }
}

resource "alicloud_polardb_on_ens_cluster" "mytest" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = "polar.mysql.x4.medium.c"
  pay_type      = "PrePaid"
  renewal_status = "AutoRenewal"
  period = 1
  auto_renew_period = 1
  db_minor_version = "8.0.2"
  description   = var.name
  ens_region_id = "sg-singapore-9"
  db_node_num = length(var.db_cluster_nodes_configs) >0 ? null : var.db_node_num
  storage_space = 20
  storage_type = "ESSDPL0"
  target_minor_version= "innovate_x86#20250311"
  vpc_id = "n-56tme6lvq1n6s3d1zn78zd12w"
  vswitch_id    = "vsw-56tme6tk908wmeco8kg3nse60"
  db_cluster_nodes_configs = {
    for node, config in var.db_cluster_nodes_configs : node => jsonencode({for k, v in config : k => v if v != null})
  }
}

resource "alicloud_polardb_on_ens_endpoint" "endpoint_1" {
  db_cluster_id           = alicloud_polardb_on_ens_cluster.mytest.id
  endpoint_type           = "Custom"
  db_cluster_nodes_ids    = alicloud_polardb_on_ens_cluster.mytest.db_cluster_nodes_ids
  auto_add_new_nodes      = "Enable"
  endpoint_config         = {"MasterAcceptReads":"on"}
  net_type                = "Private"
  db_endpoint_description = ""
  vpc_id                  = "n-56tme6lvq1n6s3d1zn78zd12w"
  vswitch_id              = "vsw-56tme6tk908wmeco8kg3nse60"
  nodes_key               = ["node_reader_1"]
  read_write_mode         = "ReadWrite"
}

resource "alicloud_polardb_on_ens_account" "account_1" {
  db_cluster_id          = alicloud_polardb_on_ens_cluster.mytest.id
  account_name           = "xdang_terraform1"
  account_password       = "Ali123789"
  account_type           = "Normal"
  account_description    = "from_terraform_modify"
}

resource "alicloud_polardb_on_ens_database" "database_1" {
  db_cluster_id         = alicloud_polardb_on_ens_cluster.mytest.id
  db_description        = "test terraforms"
  db_name               = "xdang_terraform"
}

resource "alicloud_polardb_on_ens_account_privilege" "privilege_1" {
  db_cluster_id         = alicloud_polardb_on_ens_cluster.mytest.id
  account_name          = alicloud_polardb_on_ens_account.account_1.account_name
  account_privilege     = "ReadWrite"
  db_names              = [alicloud_polardb_on_ens_database.database_1.db_name]
}
