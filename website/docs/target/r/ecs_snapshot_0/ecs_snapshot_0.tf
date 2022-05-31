resource "alicloud_ecs_snapshot" "default" {
  category       = "standard"
  description    = "Test For Terraform"
  disk_id        = "d-gw8csgxxxxxxxxx"
  retention_days = "20"
  snapshot_name  = "tf-test"
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}

