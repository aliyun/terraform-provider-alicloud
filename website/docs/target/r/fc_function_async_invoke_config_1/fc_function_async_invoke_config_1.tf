resource "alicloud_fc_function_async_invoke_config" "example" {
  service_name                 = alicloud_fc_service.example.name
  function_name                = alicloud_fc_function.example.name
  maximum_event_age_in_seconds = 60
  maximum_retry_attempts       = 0
}
