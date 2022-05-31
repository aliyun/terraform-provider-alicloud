resource "alicloud_oss_bucket" "bucket-sserule" {
  bucket = "bucket-170309-sserule"
  acl    = "private"

  server_side_encryption_rule {
    sse_algorithm = "AES256"
  }
}

resource "alicloud_oss_bucket" "bucket-sserule" {
  bucket = "bucket-170309-sserule"
  acl    = "private"

  server_side_encryption_rule {
    sse_algorithm     = "KMS"
    kms_master_key_id = "your kms key id"
  }
}
