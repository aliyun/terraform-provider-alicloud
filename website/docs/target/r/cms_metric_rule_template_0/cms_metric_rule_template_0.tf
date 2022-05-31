resource "alicloud_cms_metric_rule_template" "example" {
  metric_rule_template_name = "example_value"
  alert_templates {
    category    = "ecs"
    metric_name = "cpu_total"
    namespace   = "acs_ecs_dashboard"
    rule_name   = "tf_testAcc_new"
    escalations {
      critical {
        comparison_operator = "GreaterThanThreshold"
        statistics          = "Average"
        threshold           = "90"
        times               = "3"
      }
    }
  }
}

