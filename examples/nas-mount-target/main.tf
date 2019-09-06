resource "alicloud_nas_file_system" "main" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = "Create_FileSystem"
}

resource "alicloud_nas_access_group" "main" {
  name        = "tf-testAccNasConfigName"
  type        = "Classic"
  description = "tf-testAccNasConfig"
}

resource "alicloud_nas_mount_target" "main" {
  file_system_id    = alicloud_nas_file_system.main.id
  access_group_name = alicloud_nas_access_group.main.id
}

