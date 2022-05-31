resource "alicloud_imp_app_template" "example" {
  app_template_name = "example_value"
  component_list    = ["component.live", "component.liveRecord"]
  integration_mode  = "paasSDK"
  scene             = "business"
}

