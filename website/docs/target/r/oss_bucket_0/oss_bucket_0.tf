resource "alicloud_oss_bucket" "bucket-acl" {
  bucket = "bucket-170309-acl"
  acl    = "private"
}
