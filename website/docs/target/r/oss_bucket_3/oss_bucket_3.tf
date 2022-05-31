resource "alicloud_oss_bucket" "bucket-referer" {
  bucket = "bucket-170309-referer"
  acl    = "private"

  referer_config {
    allow_empty = false
    referers    = ["http://www.aliyun.com", "https://www.aliyun.com"]
  }
}
