resource "alicloud_nas_file_system" "foo" {
  file_system_type = "extreme"
  protocol_type    = "NFS"
  zone_id          = "cn-hangzhou-f"
  storage_type     = "standard"
  description      = "tf-testAccNasConfig"
  capacity         = "100"
}
