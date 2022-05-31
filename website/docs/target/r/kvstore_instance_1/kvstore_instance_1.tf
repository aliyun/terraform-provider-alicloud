resource "alicloud_kvstore_instance" "example" {
  db_instance_name = "tf-test-basic"
  vswitch_id       = "vsw-123456"
  security_ips = [
  "10.23.12.24"]
  instance_type  = "Redis"
  engine_version = "4.0"
  config = {
    appendonly             = "yes",
    lazyfree-lazy-eviction = "yes",
  }
  tags = {
    Created = "TF",
    For     = "Test",
  }
  resource_group_id = "rg-123456"
  zone_id           = "cn-beijing-h"
  instance_class    = "redis.master.large.default"
  payment_type      = "PrePaid"
  period            = "12"
}
