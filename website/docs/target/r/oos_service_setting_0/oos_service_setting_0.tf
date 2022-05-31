variable "name" {
  default = "tf-testaccoossetting"
}

resource "alicloud_oss_bucket" "default" {
  bucket = var.name
  acl    = "public-read-write"
}

resource "alicloud_log_project" "default" {
  name = var.name
}

resource "alicloud_oos_service_setting" "default" {
  delivery_oss_enabled      = true
  delivery_oss_key_prefix   = "path1/"
  delivery_oss_bucket_name  = alicloud_oss_bucket.default.bucket
  delivery_sls_enabled      = true
  delivery_sls_project_name = alicloud_log_project.default.name
}
