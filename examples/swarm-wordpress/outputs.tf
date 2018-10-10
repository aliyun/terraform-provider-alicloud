output "cluster_id" {
  value = "${alicloud_cs_swarm.cs_vpc.id}"
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

output "default_domain" {
  value = "${alicloud_cs_application.wordpress.default_domain}"
}
