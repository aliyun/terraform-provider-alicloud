data "alicloud_cdn_service" "ci" {
  enable               = "On"
  internet_charge_type = "PayByTraffic"
}