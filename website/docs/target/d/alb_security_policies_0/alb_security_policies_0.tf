data "alicloud_alb_security_policies" "ids" {}
output "alb_security_policy_id_1" {
  value = data.alicloud_alb_security_policies.ids.policies.0.id
}

data "alicloud_alb_security_policies" "nameRegex" {
  name_regex = "^my-SecurityPolicy"
}
output "alb_security_policy_id_2" {
  value = data.alicloud_alb_security_policies.nameRegex.policies.0.id
}

