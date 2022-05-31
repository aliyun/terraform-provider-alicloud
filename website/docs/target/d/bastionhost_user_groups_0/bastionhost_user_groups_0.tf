data "alicloud_bastionhost_user_groups" "ids" {
  instance_id = "bastionhost-cn-xxxx"
  ids         = ["1", "2"]
}
output "bastionhost_user_group_id_1" {
  value = data.alicloud_bastionhost_user_groups.ids.groups.0.id
}

data "alicloud_bastionhost_user_groups" "nameRegex" {
  instance_id = "bastionhost-cn-xxxx"
  name_regex  = "^my-UserGroup"
}
output "bastionhost_user_group_id_2" {
  value = data.alicloud_bastionhost_user_groups.nameRegex.groups.0.id
}

