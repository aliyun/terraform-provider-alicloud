resource "alicloud_fc_function_async_invoke_config" "example" {
  service_name  = alicloud_fc_service.example.name
  function_name = alicloud_fc_function.example.name
  qualifier     = "LATEST"

  # ... other configuration ...
}
