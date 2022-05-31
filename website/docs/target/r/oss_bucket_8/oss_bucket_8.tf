resource "alicloud_oss_bucket" "bucket-tags" {
  bucket = "bucket-170309-tags"
  acl    = "private"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
