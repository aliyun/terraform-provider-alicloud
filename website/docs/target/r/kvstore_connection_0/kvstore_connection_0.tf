resource "alicloud_kvstore_connection" "default" {
  connection_string_prefix = "allocatetestupdate"
  instance_id              = "r-abc123456"
  port                     = "6370"
}
