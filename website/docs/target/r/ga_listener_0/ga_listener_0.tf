resource "alicloud_ga_accelerator" "example" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}
resource "alicloud_ga_bandwidth_package" "de" {
  bandwidth      = "100"
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}
resource "alicloud_ga_bandwidth_package_attachment" "de" {
  accelerator_id       = alicloud_ga_accelerator.example.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.de.id
}
resource "alicloud_ga_listener" "example" {
  depends_on     = [alicloud_ga_bandwidth_package_attachment.de]
  accelerator_id = alicloud_ga_accelerator.example.id
  port_ranges {
    from_port = 60
    to_port   = 70
  }
}

