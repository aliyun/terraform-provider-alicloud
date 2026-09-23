resource "alicloud_maxcompute_project" "example" {
  name               = var.project_name
  specification_type = "OdpsStandard"
  order_type         = "PayAsYouGo"
}
