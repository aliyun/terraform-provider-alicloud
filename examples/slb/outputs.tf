output "slb_id" {
  value = alicloud_slb_load_balancer.instance.id
}

output "slbname" {
  value = alicloud_slb_load_balancer.instance.name
}

output "slb_tags" {
  value = alicloud_slb_load_balancer.instance.tags
}

output "slb_acl_name" {
  value = alicloud_slb_acl.acl.name
}

output "slb_acl_id" {
  value = alicloud_slb_acl.acl.id
}

output "slb_acl_entry_list" {
  value = alicloud_slb_acl.acl.entry_list
}

output "slb_listener_tcp_established_timeout" {
  value = alicloud_slb_listener.tcp.established_timeout
}

