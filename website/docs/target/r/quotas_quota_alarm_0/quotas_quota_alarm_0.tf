resource "alicloud_quotas_quota_alarm" "example" {
  quota_alarm_name  = "tf-testAcc"
  product_code      = "ecs"
  quota_action_code = "q_prepaid-instance-count-per-once-purchase"
  threshold         = "100"
  quota_dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
}

