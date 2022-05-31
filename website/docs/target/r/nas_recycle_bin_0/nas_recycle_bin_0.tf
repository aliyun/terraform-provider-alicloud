resource "alicloud_nas_file_system" "example" {
  protocol_type = "NFS"
  storage_type  = "Performance"
  description   = var.name
  encrypt_type  = "1"
}

resource "alicloud_nas_recycle_bin" "example" {
  file_system_id = alicloud_nas_file_system.example.id
  reserved_days  = 3
}
