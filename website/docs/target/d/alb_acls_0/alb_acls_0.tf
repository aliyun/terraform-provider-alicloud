data "alicloud_alb_acls" "ids" {}
output "alb_acl_id_1" {
  value = data.alicloud_alb_acls.ids.acls.0.id
}

data "alicloud_alb_acls" "nameRegex" {
  name_regex = "^my-Acl"
}
output "alb_acl_id_2" {
  value = data.alicloud_alb_acls.nameRegex.acls.0.id
}

