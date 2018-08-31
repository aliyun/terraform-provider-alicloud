package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCen_BandwidthLimit_basic(t *testing.T) {
	var cenBwpLimit cbn.CenInterRegionBandwidthLimit

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_bandwidthlimit.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenBandwidthLimitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExists("alicloud_cen_bandwidthlimit.foo", &cenBwpLimit),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.foo", "bandwidth_limit", "15"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.foo", "regions_id.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "cn-beijing", "cn-shanghai"),
				),
			},
		},
	})
}

func TestAccAlicloudCen_BandwidthLimit_update(t *testing.T) {
	var cenBwpLimit cbn.CenInterRegionBandwidthLimit

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenBandwidthLimitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthLimitConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExists("alicloud_cen_bandwidthlimit.foo", &cenBwpLimit),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.foo", "bandwidth_limit", "15"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.foo", "regions_id.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "cn-beijing", "cn-shanghai"),
				),
			},
			resource.TestStep{
				Config: testAccCenBandwidthLimitUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExists("alicloud_cen_bandwidthlimit.foo", &cenBwpLimit),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.foo", "bandwidth_limit", "17"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.foo", "regions_id.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "cn-beijing", "cn-shanghai"),
				),
			},
		},
	})
}

func TestAccAlicloudCen_BandwidthLimit_multi(t *testing.T) {
	var cenBwpLimit cbn.CenInterRegionBandwidthLimit

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenBandwidthLimitDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthLimitMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthLimitExists("alicloud_cen_bandwidthlimit.bar1", &cenBwpLimit),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.bar1", "bandwidth_limit", "12"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.bar1", "regions_id.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "cn-beijing", "cn-shanghai"),

					testAccCheckCenBandwidthLimitExists("alicloud_cen_bandwidthlimit.bar2", &cenBwpLimit),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.bar2", "bandwidth_limit", "8"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthlimit.bar2", "regions_id.#", "2"),
					testAccCheckCenBandwidthLimitRegionId(&cenBwpLimit, "cn-beijing", "cn-hangzhou"),
				),
			},
		},
	})
}

func testAccCheckCenBandwidthLimitExists(n string, cenBwpLimit *cbn.CenInterRegionBandwidthLimit) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CenBandwidthPackage ID is set")
		}

		params, err := getParaForCenBandwidthLimit(rs.Primary.ID)
		if err != nil {
			return err
		}

		cenId := params[0]
		localRegionId := params[1]
		oppositeRegionId := params[2]
		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)

		if err != nil {
			return err
		}

		*cenBwpLimit = instance
		return nil
	}
}

func testAccCheckCenBandwidthLimitDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_bandwidthlimit" {
			continue
		}

		params, err := getParaForCenBandwidthLimit(rs.Primary.ID)
		if err != nil {
			return err
		}
		// Try to find the Bandwidth Limit
		cenId := params[0]
		localRegionId := params[1]
		oppositeRegionId := params[2]
		instance, err := client.DescribeCenBandwidthLimit(cenId, localRegionId, oppositeRegionId)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.CenId != "" {
			return fmt.Errorf("CEN Bandwidth Limit still exist, CEN ID %s", instance.CenId)
		}
	}

	return nil
}

func testAccCheckCenBandwidthLimitRegionId(cenBwpLimit *cbn.CenInterRegionBandwidthLimit, regionAId string, regionBId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if (cenBwpLimit.LocalRegionId == regionAId && cenBwpLimit.OppositeRegionId == regionBId) ||
			(cenBwpLimit.OppositeRegionId == regionBId && cenBwpLimit.LocalRegionId == regionAId) {
			return nil
		} else {
			return fmt.Errorf("CEN %s BandwidthLimit Region ID error", cenBwpLimit.CenId)
		}
	}
}

const testAccCenBandwidthLimitConfig = `
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

const testAccCenBandwidthLimitUpdate = `
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
	provider = "alicloud.bj"
	name = "terraform-yl-01"
	description = "terraform01"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
	provider = "alicloud.bj"
    bandwidth = 20
    geographic_region_id = [
		"China",
		"China"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "bwp_attach" {
	provider = "alicloud.bj"
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
	provider = "alicloud.bj"
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
	provider = "alicloud.bj"
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidthlimit" "foo" {
	provider = "alicloud.bj"
     cen_id = "${alicloud_cen.cen.id}"
     regions_id = [
                    "cn-beijing",
                    "cn-shanghai"]
     bandwidth_limit = 17
     depends_on = [
        "alicloud_cen_bandwidthpackage_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}
`

const testAccCenBandwidthLimitMulti = `
provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "sh"
  region = "cn-shanghai"
}

provider "alicloud" {
  alias = "hz"
  region = "cn-hangzhou"
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

resource "alicloud_vpc" "vpc3" {
  provider = "alicloud.hz"
  name = "terraform-03"
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen" "cen" {
	provider = "alicloud.bj"
     name = "terraform-yl-01"
     description = "terraform01"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
	provider = "alicloud.bj"
    bandwidth = 20
    geographic_region_id = [
		"China",
		"China"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "bwp_attach" {
	provider = "alicloud.bj"
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
	provider = "alicloud.bj"
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
	provider = "alicloud.bj"
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_3" {
	provider = "alicloud.bj"
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc3.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-hangzhou"
}

resource "alicloud_cen_bandwidthlimit" "bar1" {
	provider = "alicloud.bj"
     cen_id = "${alicloud_cen.cen.id}"
     regions_id = [
                    "cn-beijing",
                    "cn-shanghai"]
     bandwidth_limit = 12
     depends_on = [
        "alicloud_cen_bandwidthpackage_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}

resource "alicloud_cen_bandwidthlimit" "bar2" {
	provider = "alicloud.bj"
     cen_id = "${alicloud_cen.cen.id}"
     regions_id = [
                    "cn-beijing",
                    "cn-hangzhou"]
     bandwidth_limit = 8
     depends_on = [
        "alicloud_cen_bandwidthpackage_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_3"]
}
`
