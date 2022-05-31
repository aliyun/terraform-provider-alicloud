resource "alicloud_ecs_disk" "example" {
  zone_id     = "cn-beijing-b"
  disk_name   = "tf-test"
  description = "Hello ecs disk."
  category    = "cloud_efficiency"
  size        = "30"
  encrypted   = true
  kms_key_id  = "2a6767f0-a16c-4679-a60f-13bf*****"
  tags = {
    Name = "TerraformTest"
  }
}

