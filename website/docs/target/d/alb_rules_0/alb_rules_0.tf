data "alicloud_alb_rules" "ids" {
  ids = ["example_id"]
}
output "alb_rule_id_1" {
  value = data.alicloud_alb_rules.ids.rules.0.id
}

data "alicloud_alb_rules" "nameRegex" {
  name_regex = "^my-Rule"
}
output "alb_rule_id_2" {
  value = data.alicloud_alb_rules.nameRegex.rules.0.id
}

