resource "alicloud_oss_bucket" "bucket-storageclass" {
  bucket        = "bucket-170309-storageclass"
  storage_class = "IA"
}
