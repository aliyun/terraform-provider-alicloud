data "alicloud_bastionhost_host_groups" "ids" {
  instance_id = "bastionhost-cn-tl3xxxxxxx"
  ids         = ["example_value-1", "example_value-2"]
}
output "bastionhost_host_group_id_1" {
  value = data.alicloud_bastionhost_host_groups.ids.groups.0.id
}

data "alicloud_bastionhost_host_groups" "nameRegex" {
  instance_id = "bastionhost-cn-tl3xxxxxxx"
  name_regex  = "^my-HostGroup"
}
output "bastionhost_host_group_id_2" {
  value = data.alicloud_bastionhost_host_groups.nameRegex.groups.0.id
}

