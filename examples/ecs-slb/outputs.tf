output "slb_id" {
  value = alicloud_slb_load_balancer.instance.id
}

output "slbname" {
  value = alicloud_slb_load_balancer.instance.name
}

output "hostname_list" {
  value = join(",", alicloud_instance.instance.*.instance_name)
}

output "ecs_ids" {
  value = join(",", alicloud_instance.instance.*.id)
}

output "slb_backendserver" {
  value = alicloud_slb_attachment.default.backend_servers
}

