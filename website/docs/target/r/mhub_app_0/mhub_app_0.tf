variable "name" {
  default = "example_value"
}
resource "alicloud_mhub_app" "default" {
  app_name     = var.name
  product_id   = alicloud_mhub_product.default.id
  package_name = "com.test.android"
  type         = "Android"
}

