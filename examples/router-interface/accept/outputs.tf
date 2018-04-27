// Output the IDs of the ECS instances created
output "vpc_id" {
  value = "${var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id}"
}

output "vswitch_ids" {
  value = "${join(",", alicloud_vswitch.vswitches.*.id)}"
}

output "availability_zones" {
  value = "${join(",", alicloud_vswitch.vswitches.*.availability_zone)}"
}

output "router_id" {
  value = "${join("", alicloud_route_entry.route_entry.*.router_id)}"
}

output "route_table_id" {
  value = "${join("", alicloud_route_entry.route_entry.*.route_table_id)}"
}

output "interface_id" {
  value = "${join("", alicloud_router_interface.interface.*.id)}"
}

output "router_type" {
  value = "${join("", alicloud_router_interface.interface.*.router_type)}"
}
