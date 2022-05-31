resource "alicloud_ecs_auto_snapshot_policy" "example" {
  name            = "tf-testAcc"
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}

