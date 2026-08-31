output "cluster_id" {
  value = alicloud_cassandra_cluster.default.id
}
output "dc1_id" {
  value = alicloud_cassandra_cluster.default.zone_id
}
output "dc2_id" {
  value = alicloud_cassandra_data_center.dc_2.zone_id
}