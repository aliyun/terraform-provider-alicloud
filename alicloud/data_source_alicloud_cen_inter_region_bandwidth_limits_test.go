package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudCenInterRegionBandwidthLimitsDataSource_cen_id(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceCenIdConfigSet,
			},
			resource.TestStep{
				Config: testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceCenIdConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cen_inter_region_bandwidth_limits.limit"),
					resource.TestCheckResourceAttr("data.alicloud_cen_inter_region_bandwidth_limits.limit",
						"cen_inter_region_bandwidth_limits.0.local_region_id", "cn-beijing"),
					resource.TestCheckResourceAttr("data.alicloud_cen_inter_region_bandwidth_limits.limit",
						"cen_inter_region_bandwidth_limits.0.opposite_region_id", "cn-shanghai"),
					resource.TestCheckResourceAttr("data.alicloud_cen_inter_region_bandwidth_limits.limit",
						"cen_inter_region_bandwidth_limits.0.status", "Active"),
					resource.TestCheckResourceAttr("data.alicloud_cen_inter_region_bandwidth_limits.limit",
						"cen_inter_region_bandwidth_limits.0.bandwidth_limit", "15"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceCenIdConfigSet = `
provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "sh"
  region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
  provider = "alicloud.bj"
  name = "terraform-01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider = "alicloud.sh"
  name = "terraform-02"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen" "cen" {
     name = "terraform-01"
     description = "terraform01"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"China"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "bwp_attach" {
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidthlimit" "foo" {
     cen_id = "${alicloud_cen.cen.id}"
     regions_id = [
                    "cn-beijing",
                    "cn-shanghai"]
     bandwidth_limit = 15
     depends_on = [
        "alicloud_cen_bandwidthpackage_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}
`

const testAccCheckAlicloudCenInterRegionBandwidthLimitsDataSourceCenIdConfig = `
provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "sh"
  region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
  provider = "alicloud.bj"
  name = "terraform-01"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  provider = "alicloud.sh"
  name = "terraform-02"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen" "cen" {
     name = "terraform-01"
     description = "terraform01"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"China"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "bwp_attach" {
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidthlimit" "foo" {
     cen_id = "${alicloud_cen.cen.id}"
     regions_id = [
                    "cn-beijing",
                    "cn-shanghai"]
     bandwidth_limit = 15
     depends_on = [
        "alicloud_cen_bandwidthpackage_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}

data "alicloud_cen_inter_region_bandwidth_limits" "limit" {
	cen_id = "${alicloud_cen.cen.id}"
}
`
