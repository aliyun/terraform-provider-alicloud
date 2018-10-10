output "cluster_id" {
  value = "${alicloud_container_cluster.cs_vpc.id}"
}

output "vpc_id" {
  value = "${alicloud_vpc.main.id}"
}

output "vswitch_id" {
  value = "${alicloud_vswitch.main.id}"
}

output "availability_zone" {
  value = "${alicloud_vswitch.main.availability_zone}"
}
