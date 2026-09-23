output "name" {
  description = "Name of ONS Instance."
  value       = alicloud_ons_instance.instance.name
}

output "instance_id" {
  description = "The ID used to identify ons.instance resource."
  value       = alicloud_ons_instance.instance.id
}

output "group_id" {
  description = "Name of ONS Group."
  value       = alicloud_ons_group.default.group_id
}

output "group_remark" {
  description = "This attribute is a concise description of group."
  value       = alicloud_ons_group.default.remark
}

