data "alicloud_account" "default" {}

data "alicloud_regions" "default" {
  current = true
}
resource "alicloud_cdn_fc_trigger" "example" {
  event_meta_name    = "LogFileCreated"
  event_meta_version = "1.0.0"
  notes              = "example_value"
  role_arn           = "acs:ram::${data.alicloud_account.default.id}:role/aliyuncdneventnotificationrole"
  source_arn         = "acs:cdn:*:${data.alicloud_account.default.id}:domain/example.com"
  trigger_arn        = "acs:fc:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:services/FCTestService/functions/printEvent/triggers/testtrigger"
}
