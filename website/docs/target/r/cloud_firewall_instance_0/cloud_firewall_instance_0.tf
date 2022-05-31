resource "alicloud_cloud_firewall_instance" "example" {
  payment_type    = "Subscription"
  spec            = "premium_version"
  ip_number       = 20
  band_width      = 10
  cfw_log         = false
  cfw_log_storage = 1000
  cfw_service     = false
  period          = 6
}
