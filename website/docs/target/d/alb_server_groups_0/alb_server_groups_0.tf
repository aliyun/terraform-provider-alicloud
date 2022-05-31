data "alicloud_alb_server_groups" "ids" {}
output "alb_server_group_id_1" {
  value = data.alicloud_alb_server_groups.ids.groups.0.id
}

data "alicloud_alb_server_groups" "nameRegex" {
  name_regex = "^my-ServerGroup"
}
output "alb_server_group_id_2" {
  value = data.alicloud_alb_server_groups.nameRegex.groups.0.id
}

