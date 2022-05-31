resource "alicloud_ga_accelerator" "default" {
  duration         = 1
  spec             = "1"
  accelerator_name = var.name
  auto_use_coupon  = true
  description      = var.name
}
resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth              = 100
  type                   = "Basic"
  bandwidth_type         = "Basic"
  payment_type           = "PayAsYouGo"
  billing_type           = "PayBy95"
  ratio                  = 30
  bandwidth_package_name = var.name
  auto_pay               = true
  auto_use_coupon        = true
}
resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_accelerator_spare_ip_attachment" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  spare_ip       = "127.0.0.1"
}
