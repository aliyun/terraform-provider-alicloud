resource "alicloud_kvstore_audit_log_config" "example" {
  instance_id = "r-abc123455"
  db_audit    = true
  retention   = 1
}

