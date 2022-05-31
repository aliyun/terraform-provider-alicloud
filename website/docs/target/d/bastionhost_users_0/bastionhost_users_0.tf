data "alicloud_bastionhost_users" "ids" {
  instance_id = "example_value"
  ids         = ["1", "10"]
}
output "bastionhost_user_id_1" {
  value = data.alicloud_bastionhost_users.ids.users.0.id
}

data "alicloud_bastionhost_users" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-User"
}
output "bastionhost_user_id_2" {
  value = data.alicloud_bastionhost_users.nameRegex.users.0.id
}

