resource "alicloud_nas_auto_snapshot_policy" "example" {
  auto_snapshot_policy_name = "example_value"
  repeat_weekdays           = ["3", "4", "5"]
  retention_days            = 30
  time_points               = ["3", "4", "5"]
}
