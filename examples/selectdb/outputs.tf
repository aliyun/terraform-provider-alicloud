output "db_cluster" {
  value = data.alicloud_selectdb_db_clusters.default.ids.0
}

output "db_instance" {
  value = data.alicloud_selectdb_db_instances.default.ids.0
}
