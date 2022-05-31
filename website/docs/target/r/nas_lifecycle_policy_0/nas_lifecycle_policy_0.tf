resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Capacity"
}

resource "alicloud_nas_lifecycle_policy" "example" {
  file_system_id        = alicloud_nas_file_system.example.id
  lifecycle_policy_name = "my-LifecyclePolicy"
  lifecycle_rule_name   = "DEFAULT_ATIME_14"
  storage_type          = "InfrequentAccess"
}
