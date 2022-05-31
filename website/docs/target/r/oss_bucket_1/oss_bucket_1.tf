resource "alicloud_oss_bucket" "bucket-website" {
  bucket = "bucket-170309-website"

  website {
    index_document = "index.html"
    error_document = "error.html"
  }
}
