resource "alicloud_maxcompute_project" "example" {
  project_name       = "tf_maxcompute_project"
  specification_type = "OdpsStandard"
  order_type         = "PayAsYouGo"
}
