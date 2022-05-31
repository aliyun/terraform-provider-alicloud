resource "alicloud_ecs_auto_snapshot_policy_attachment" "example" {
  auto_snapshot_policy_id = "s-ge465xxxx"
  disk_id                 = "d-gw835xxxx"
}

