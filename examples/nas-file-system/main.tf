resource "alicloud_nas_file_system" "foo" {
  protocol_type    = "NFS"
  storage_type     = "standard"
  file_system_type = "extreme"
  capacity         = "100"
  zone_id          = "cn-hangzhou-f"
}