output "cen_id" {
  value = alicloud_cen_instance.cen.id
}

output "vpc_1_id" {
  value = alicloud_vpc.vpc_1.id
}

output "vpc_2_id" {
  value = alicloud_vpc.vpc_2.id
}

output "bwp_id" {
  value = alicloud_cen_bandwidth_package.bwp.id
}

output "vswitch_id" {
  value = alicloud_vswitch.default.id
}

output "security_group_id" {
  value = alicloud_security_group.default.id
}

output "instance_id" {
  value = alicloud_instance.default.id
}

output "bwp_attach_id" {
  value = alicloud_cen_bandwidth_package_attachment.bwp_attach.id
}

output "vpc_attach_1_id" {
  value = alicloud_cen_instance_attachment.vpc_attach_1.id
}

output "vpc_attach_2_id" {
  value = alicloud_cen_instance_attachment.vpc_attach_2.id
}

output "limit_id" {
  value = alicloud_cen_bandwidth_limit.limit.id
}

output "route_entry_id" {
  value = alicloud_route_entry.route_entry.id
}

output "cen_route_entry_id" {
  value = alicloud_cen_route_entry.cen_route_entry.id
}

