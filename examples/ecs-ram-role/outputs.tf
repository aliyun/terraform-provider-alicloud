// Output the IDs of the ECS instances created
output "ecs_instance_ids" {
  value = "${join(",", alicloud_instance.instances.*.id)}"
}

// Output the IDs of the ECS disks created
output "ecs_disk_ids" {
  value = "${join(",", alicloud_disk.disks.*.id)}"
}

// Output the ID of the new VPC created
output "vpc_id" {
  value = "${join(",", alicloud_vpc.vpc.*.id)}"
}

// Output the ID of the new VSwitch created
output "vswitch_id" {
  value = "${join(",", alicloud_vswitch.vswitch.*.id)}"
}

// Output the ID of the new Security Group created
output "security_group_id" {
  value = "${join(",", alicloud_security_group.group.*.id)}"
}