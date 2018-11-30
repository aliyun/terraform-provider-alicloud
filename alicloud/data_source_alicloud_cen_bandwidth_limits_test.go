package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenBandwidthLimitsDataSource_instance_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceCenIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_limits.limit"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.local_region_id", "cn-beijing"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.opposite_region_id", "us-west-1"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.status", "Active"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.bandwidth_limit", "15"),
					resource.TestCheckResourceAttrSet("data.alicloud_cen_bandwidth_limits.limit", "limits.0.instance_id"),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthLimitsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_bandwidth_limits.limit"),
					resource.TestCheckResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.#", "0"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.local_region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.opposite_region_id"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.status"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.bandwidth_limit"),
					resource.TestCheckNoResourceAttr("data.alicloud_cen_bandwidth_limits.limit", "limits.0.instance_id"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceCenIdConfig = `
provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "us"
  region = "us-west-1"
}

resource "alicloud_vpc" "vpc1" {
  provider = "alicloud.bj"
  name = "tf-testAccTerraform-01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider = "alicloud.us"
  name = "tf-testAccTerraform-02"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "cen" {
     name = "tf-testAccTerraform-01"
     description = "tf-testAccTerraform01"
}

resource "alicloud_cen_bandwidth_package" "bwp" {
    name = "tf-testAccCenBandwidthPackage"
    bandwidth = 20
    geographic_region_ids = [
		"China",
		"North-America"]
}

resource "alicloud_cen_bandwidth_package_attachment" "bwp_attach" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    bandwidth_package_id = "${alicloud_cen_bandwidth_package.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_region_id = "us-west-1"
}

resource "alicloud_cen_bandwidth_limit" "foo" {
     instance_id = "${alicloud_cen_instance.cen.id}"
     region_ids = ["cn-beijing",
                   "us-west-1"]
     bandwidth_limit = 15
     depends_on = [
        "alicloud_cen_bandwidth_package_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}

data "alicloud_cen_bandwidth_limits" "limit" {
	instance_ids = ["${alicloud_cen_bandwidth_limit.foo.instance_id}"]
}
`
const testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceEmpty = `
data "alicloud_cen_bandwidth_limits" "limit" {
	instance_ids = ["cen-cidwnnenc"]
}
`
