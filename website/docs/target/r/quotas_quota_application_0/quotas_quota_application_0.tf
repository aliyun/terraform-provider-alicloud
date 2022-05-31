resource "alicloud_quotas_quota_application" "example" {
  notice_type       = "0"
  desire_value      = "100"
  product_code      = "ess"
  quota_action_code = "q_db_instance"
  reason            = "For Terraform Test"
  dimensions {
    key   = "regionId"
    value = "cn-hangzhou"
  }
}

