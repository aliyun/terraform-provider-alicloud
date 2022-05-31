resource "random_uuid" "this" {}

resource "alicloud_cms_group_metric_rule" "this" {
  group_id = "539****"
  rule_id  = random_uuid.this.id

  category    = "ecs"
  namespace   = "acs_ecs_dashboard"
  metric_name = "cpu_total"
  period      = "60"

  group_metric_rule_name = "tf-testacc-rule-name"
  email_subject          = "tf-testacc-rule-name-warning"
  interval               = "3600"
  silence_time           = 85800
  no_effective_interval  = "00:00-05:30"
  webhook                = "http://www.aliyun.com"
  escalations {
    warn {
      comparison_operator = "GreaterThanOrEqualToThreshold"
      statistics          = "Average"
      threshold           = "90"
      times               = 3
    }
    info {
      comparison_operator = "LessThanLastWeek"
      statistics          = "Average"
      threshold           = "90"
      times               = 5
    }
  }
}
