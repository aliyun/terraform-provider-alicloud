package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCen_BandwidthPackage_basic(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_bandwidthpackage.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthPackageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidthpackage.foo", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.foo", "geographic_region_id.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.foo", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.foo", "status"),
				),
			},
		},
	})
}

func TestAccAlicloudCen_BandwidthPackage_update(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthPackageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidthpackage.foo", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.foo", "bandwidth", "20"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.foo", "geographic_region_id.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.foo", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.foo", "status"),
				),
			},
			resource.TestStep{
				Config: testAccCenBandwidthPackageConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidthpackage.foo", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.foo", "bandwidth", "25"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.foo", "geographic_region_id.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.foo", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.foo", "status"),
				),
			},
		},
	})
}

func TestAccAlicloudCen_BandwidthPackage_muti(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthPackageConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidthpackage.bar1", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.bar1", "bandwidth", "11"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.bar1", "geographic_region_id.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.bar1", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.bar1", "status"),

					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidthpackage.bar2", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.bar2", "bandwidth", "13"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidthpackage.bar2", "geographic_region_id.#", "1"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "China"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.bar2", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidthpackage.bar2", "status"),
				),
			},
		},
	})
}

func testAccCheckCenBandwidthPackageExists(n string, cenBwp *cbn.CenBandwidthPackage) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CenBandwidthPackage ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)
		instance, err := client.DescribeCenBandwidthPackage(rs.Primary.ID)

		if err != nil {
			return err
		}

		*cenBwp = instance
		return nil
	}
}

func testAccCheckCenBandwidthPackageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_bandwidthpackage" {
			continue
		}

		// Try to find the CEN
		instance, err := client.DescribeCenBandwidthPackage(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.CenBandwidthPackageId != "" {
			return fmt.Errorf("CEN Bandwidth Package %s still exist", instance.CenBandwidthPackageId)
		}
	}

	return nil
}

func testAccCheckCenBandwidthPackageRegionId(cenBwp *cbn.CenBandwidthPackage, regionAId string, regionBId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		responseRegionAId := convertGeographicRegionId(cenBwp.GeographicRegionAId)
		responseRegionBId := convertGeographicRegionId(cenBwp.GeographicRegionBId)
		if (responseRegionAId == regionAId && responseRegionBId == regionBId) ||
			(responseRegionAId == regionBId && responseRegionBId == regionAId) {
			return nil
		} else {
			return fmt.Errorf("CEN Bandwidth Package %s geographic region ID error", cenBwp.CenBandwidthPackageId)
		}
	}
}

const testAccCenBandwidthPackageConfig = `
resource "alicloud_cen_bandwidthpackage" "foo" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"Asia-Pacific"]
}
`

const testAccCenBandwidthPackageConfigUpdate = `
resource "alicloud_cen_bandwidthpackage" "foo" {
    bandwidth = 25
    geographic_region_id = [
		"China",
		"Asia-Pacific"]
}
`

const testAccCenBandwidthPackageConfigMulti = `
resource "alicloud_cen_bandwidthpackage" "bar1" {
    bandwidth = 11
    geographic_region_id = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidthpackage" "bar2" {
    bandwidth = 13
    geographic_region_id = [
		"China",
		"China"]
}
`
