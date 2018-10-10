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
  value = "${var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id}"
}

// Output the ID of the new VSwitch created
output "vswitch_id" {
  value = "${var.vswitch_id == "" ? join("", alicloud_vswitch.vswitch.*.id) : var.vswitch_id}"
}

// Output the ID of the new Security Group created
output "security_group_id" {
  value = "${var.sg_id == "" ? join("", alicloud_security_group.group.*.id) : var.sg_id}"
}
