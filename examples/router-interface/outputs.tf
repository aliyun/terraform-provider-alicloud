// Output the IDs of the ECS instances created
output "vpc_id" {
  value = "${var.init_vpc_id == "" ? join("", alicloud_vpc.init.*.id) : var.init_vpc_id}"
}

output "accepting_vpc_id" {
  value = "${var.accept_vpc_id == "" ? join("", alicloud_vpc.accept.*.id) : var.accept_vpc_id}"
}

output "interface_id" {
  value = "${join("", alicloud_router_interface.init.*.id)}"
}

output "accepting_interface_id" {
  value = "${join("", alicloud_router_interface.accept.*.id)}"
}
