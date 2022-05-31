resource "alicloud_oss_bucket" "bucket-redundancytype" {
  bucket          = "bucket_name"
  redundancy_type = "ZRS"

  # ... other configuration ...
}
