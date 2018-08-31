output "cen_id" {
  value = "${alicloud_cen.cen.id}"
}

output "vpc_1_id" {
  value = "${alicloud_vpc.vpc_1.id}"
}

output "vpc_2_id" {
  value = "${alicloud_vpc.vpc_2.id}"
}

output "bwp_id" {
  value = "${alicloud_cen_bandwidthpackage.bwp.id}"
}

output "instance_id" {
  value = "${alicloud_instance.default.id}"
}

output "vswitch_id" {
  value = "${alicloud_vswitch.default.id}"
}

output "security_group_id" {
  value = "${alicloud_security_group.default.id}"
}