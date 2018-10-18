output "slb_id" {
  value = "${alicloud_slb.instance.id}"
}

output "slbname" {
  value = "${alicloud_slb.instance.name}"
}

output "slb_acl_name" {
  value = "${alicloud_slb_acl.acl.name}"
}

output "slb_acl_id" {
  value = "${alicloud_slb_acl.acl.id}"
}

output "slb_acl_entry_list" {
  value = "${alicloud_slb_acl.acl.entry_list}"
}
