resource "alicloud_graph_database_db_instance" "example" {
  db_node_class            = "gdb.r.2xlarge"
  db_instance_network_type = "vpc"
  db_version               = "1.0"
  db_instance_category     = "HA"
  db_instance_storage_type = "cloud_ssd"
  db_node_storage          = "example_value"
  payment_type             = "PayAsYouGo"
  db_instance_description  = "example_value"
}

