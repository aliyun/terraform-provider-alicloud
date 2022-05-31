resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "example" {
  bandwidth      = 20
  type           = "Basic"
  bandwidth_type = "Basic"
  duration       = 1
  auto_pay       = true
  ratio          = 30
}
resource "alicloud_ga_bandwidth_package_attachment" "example" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.example.id
}
resource "alicloud_ga_ip_set" "example" {
  depends_on           = [alicloud_ga_bandwidth_package_attachment.example]
  accelerate_region_id = "cn-hangzhou"
  bandwidth            = "5"
  accelerator_id       = alicloud_ga_accelerator.example.id
}

