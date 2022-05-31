data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_oos_state_configuration" "default" {
  template_name       = "ACS-ECS-InventoryDataCollection"
  configure_mode      = "ApplyOnly"
  description         = var.name
  schedule_type       = "rate"
  schedule_expression = "1 hour"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.ids.0
  targets             = "{\"Filters\": [{\"Type\": \"All\", \"Parameters\": {\"InstanceChargeType\": \"PrePaid\"}}], \"ResourceType\": \"ALIYUN::ECS::Instance\"}"
  parameters          = "{\"policy\": {\"ACS:Application\": {\"Collection\": \"Enabled\"}}}"
  tags = {
    Created = "TF"
    For     = "Test"
  }
}
