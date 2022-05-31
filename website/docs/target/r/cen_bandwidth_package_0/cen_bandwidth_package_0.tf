resource "alicloud_cen_bandwidth_package" "example" {
  bandwidth                  = 5
  cen_bandwidth_package_name = "tf-testAccCenBandwidthPackageConfig"
  geographic_region_a_id     = "China"
  geographic_region_b_id     = "China"
}
