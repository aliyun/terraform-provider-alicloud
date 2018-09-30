package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenBandwidthPackagesDataSource_instance_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig,
			},
			{
				Config: testAccCheckAlicloudCenBandwidthPackagesDataSourceInstanceIdReadConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_a_id", "Asia-Pacific"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_b_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.status", "InUse"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth", "5"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.business_status", "Normal"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth_package_charge_type", "POSTPAY"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.name", "tf-testAccCenBwpName"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.instance_id"),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackagesDataSource_bandwidth_package_nameRegex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenBandwidthPackagesDataSourceNameRegexConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_a_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_b_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.status", "Idle"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth", "5"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.business_status", "Normal"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth_package_charge_type", "POSTPAY"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.name", "tf-testAccCenBwpName1"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.instance_id", ""),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackagesDataSource_multi_bandwith_packages(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenBandwidthPackagesDataSourceMultiBwpIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.#", "6"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_a_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_b_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.status", "Idle"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth", "5"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.business_status", "Normal"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth_package_charge_type", "POSTPAY"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.name", "tf-testAccCenBwpName"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.creation_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.instance_id", ""),
				),
			},
		},
	})
}

const testAccCheckAlicloudCenBandwidthPackagesDataSourceConfig = `
resource "alicloud_cen_instance" "cen" {
	name = "tf-testAccCenConfig"
	description = "tf-testAccCenConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
    bandwidth = 5
	name = "tf-testAccCenBwpName"
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidth_package_attachment" "foo" {
	instance_id = "${alicloud_cen_instance.cen.id}"
	bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}
`

const testAccCheckAlicloudCenBandwidthPackagesDataSourceInstanceIdReadConfig = `
resource "alicloud_cen_instance" "cen" {
	name = "tf-testAccCenConfig"
	description = "tf-testAccCenConfigDescription"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
    bandwidth = 5
	name = "tf-testAccCenBwpName"
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidth_package_attachment" "foo" {
	instance_id = "${alicloud_cen_instance.cen.id}"
	bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}

data "alicloud_cen_bandwidth_packages" "tf-testAccCenBwp" {
	instance_id = "${alicloud_cen_instance.cen.id}"
}`

const testAccCheckAlicloudCenBandwidthPackagesDataSourceNameRegexConfig = `
resource "alicloud_cen_bandwidth_package" "tf-testAccCenBwp1" {
	name = "tf-testAccCenBwpName1"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"China"]
}

data "alicloud_cen_bandwidth_packages" "tf-testAccCenBwp" {
	name_regex = "${alicloud_cen_bandwidth_package.tf-testAccCenBwp1.name}"
}
`

const testAccCheckAlicloudCenBandwidthPackagesDataSourceMultiBwpIdConfig = `
resource "alicloud_cen_bandwidth_package" "bwp" {
	name = "tf-testAccCenBwpName"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"China"]
	count = 6
}

data "alicloud_cen_bandwidth_packages" "tf-testAccCenBwp" {
	ids = ["${alicloud_cen_bandwidth_package.bwp.*.id}"]
}
`
