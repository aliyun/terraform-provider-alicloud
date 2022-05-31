resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "test file system"
}

resource "alicloud_nas_access_group" "example" {
  access_group_name = "test_name"
  access_group_type = "Classic"
  description       = "test access group"
}

resource "alicloud_nas_mount_target" "example" {
  file_system_id    = alicloud_nas_file_system.example.id
  access_group_name = alicloud_nas_access_group.example.access_group_name
}
