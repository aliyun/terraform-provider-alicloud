data "alicloud_adb_db_clusters" "example" {
  description_regex = "example"
}

output "first_adb_db_cluster_id" {
  value = data.alicloud_adb_db_clusters.example.clusters.0.id
}
