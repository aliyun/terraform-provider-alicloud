data "alicloud_arms_dispatch_rules" "ids" {}
output "arms_dispatch_rule_id_1" {
  value = data.alicloud_arms_dispatch_rules.ids.rules.0.id
}

data "alicloud_arms_dispatch_rules" "nameRegex" {
  name_regex = "^my-DispatchRule"
}
output "arms_dispatch_rule_id_2" {
  value = data.alicloud_arms_dispatch_rules.nameRegex.rules.0.id
}

