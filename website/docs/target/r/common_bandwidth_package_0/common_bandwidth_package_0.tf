resource "alicloud_common_bandwidth_package" "foo" {
  bandwidth              = "1000"
  internet_charge_type   = "PayByBandwidth"
  bandwidth_package_name = "test-common-bandwidth-package"
  description            = "test-common-bandwidth-package"
}
