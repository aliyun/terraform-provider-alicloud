data "alicloud_arms_prometheus_alert_rules" "ids" {
  cluster_id = "example_value"
  ids        = ["example_value-1", "example_value-2"]
}
output "arms_prometheus_alert_rule_id_1" {
  value = data.alicloud_arms_prometheus_alert_rules.ids.rules.0.id
}

data "alicloud_arms_prometheus_alert_rules" "nameRegex" {
  cluster_id = "example_value"
  name_regex = "^my-PrometheusAlertRule"
}
output "arms_prometheus_alert_rule_id_2" {
  value = data.alicloud_arms_prometheus_alert_rules.nameRegex.rules.0.id
}

