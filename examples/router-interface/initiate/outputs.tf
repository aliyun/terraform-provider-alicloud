// Output the IDs of the ECS instances created
output "vpc_id" {
  value = "${var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id}"
}

output "vswitch_id" {
  value = "${join(",", alicloud_vswitch.vswitches.*.id)}"
}

output "router_id" {
  value = "${join("", alicloud_route_entry.route_entry.*.router_id)}"
}

output "route_table_id" {
  value = "${join("", alicloud_route_entry.route_entry.*.route_table_id)}"
}
