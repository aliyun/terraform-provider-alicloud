resource "alicloud_oss_bucket" "bucket-target" {
  bucket = "bucket-170309-acl"
  acl    = "public-read"
}

resource "alicloud_oss_bucket" "bucket-logging" {
  bucket = "bucket-170309-logging"

  logging {
    target_bucket = alicloud_oss_bucket.bucket-target.id
    target_prefix = "log/"
  }
}
