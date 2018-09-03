package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenBandwidthPackageDataSource_bandwidth_package_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenBandwidthPackageDataSourceBwpIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_packages.bwp"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.geographic_region_a_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.geographic_region_b_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.status", "Idle"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.bandwidth", "5"),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackageDataSource_bandwidth_package_nameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenBandwidthPackageDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_packages.bwp"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.geographic_region_a_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.geographic_region_b_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.status", "Idle"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.bwp", "cen_bandwidth_packages.0.bandwidth", "5"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCenBandwidthPackageDataSourceBwpIdConfig = `
resource "alicloud_cen_bandwidthpackage" "bwp" {
    bandwidth = 5
    geographic_region_id = [
		"China",
		"China"]
}

data "alicloud_cen_bandwidth_packages" "bwp" {
	cen_bandwidth_package_ids = ["${alicloud_cen_bandwidthpackage.bwp.id}"]
}
`

const testAccCheckAlicloudCenBandwidthPackageDataSourceNameRegexConfig = `
resource "alicloud_cen_bandwidthpackage" "bwp" {
	name = "terraformTestAccName"
    bandwidth = 5
    geographic_region_id = [
		"China",
		"China"]
}

data "alicloud_cen_bandwidth_packages" "bwp" {
	name_regex = "${alicloud_cen_bandwidthpackage.bwp.name}"
}
`
