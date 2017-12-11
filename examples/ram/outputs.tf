output "alias" {
  value = "${alicloud_ram_alias.alias.id}"
}

output "user" {
  value = "${alicloud_ram_user.user.id}"
}

output "group" {
  value = "${alicloud_ram_group.group.id}"
}

output "role" {
  value = "${alicloud_ram_role.role.id}"
}

output "policy" {
  value = "${alicloud_ram_policy.policy.id}"
}
