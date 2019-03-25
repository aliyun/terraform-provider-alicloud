package alicloud

import (
	"testing"

	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenBandwidthPackagesDataSource_instance_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
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
	rand := time.Now().UnixNano()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenBandwidthPackagesDataSourceNameRegexConfig(defaultRegionToTest, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_a_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_b_id", "China"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.status", "Idle"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth", "5"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.business_status", "Normal"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth_package_charge_type", "POSTPAY"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.description", ""),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.name",
						fmt.Sprintf("tf-testAccCenBwpName1-%s-%d", defaultRegionToTest, rand)),
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
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
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

func TestAccAlicloudCenBandwidthPackagesDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, DomesticSite)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenBandwidthPackagesDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_a_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.geographic_region_b_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.business_status"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.bandwidth_package_charge_type"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.description"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.name"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.creation_time"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_packages.tf-testAccCenBwp", "packages.0.instance_id"),
				),
			},
		},
	})
}

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
	instance_id = "${alicloud_cen_bandwidth_package_attachment.foo.instance_id}"
}
`

func testAccCheckAlicloudCenBandwidthPackagesDataSourceNameRegexConfig(region string, rand int64) string {
	return fmt.Sprintf(`
		resource "alicloud_cen_bandwidth_package" "tf-testAccCenBwp1" {
			name = "tf-testAccCenBwpName1-%s-%d"
    		bandwidth = 5
    		geographic_region_ids = [
				"China",
				"China"]
		}

		data "alicloud_cen_bandwidth_packages" "tf-testAccCenBwp" {
			name_regex = "${alicloud_cen_bandwidth_package.tf-testAccCenBwp1.name}"
		}
		`, region, rand)
}

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

const testAccCheckAlicloudCenBandwidthPackagesDataSourceEmpty = `
data "alicloud_cen_bandwidth_packages" "tf-testAccCenBwp" {
	name_regex = "^tf-testacc-fake-name"
}
`
