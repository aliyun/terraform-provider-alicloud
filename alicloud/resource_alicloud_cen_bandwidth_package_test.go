package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_cen_bandwidth_package", &resource.Sweeper{
		Name: "alicloud_cen_bandwidth_package",
		F:    testSweepCenBandwidthPackage,
	})
}

func testSweepCenBandwidthPackage(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"testAcc",
	}

	var insts []cbn.CenBandwidthPackage
	req := cbn.CreateDescribeCenBandwidthPackagesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DescribeCenBandwidthPackages(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving CEN Bandwidth Package: %s", err)
		}
		resp, _ := raw.(*cbn.DescribeCenBandwidthPackagesResponse)
		if resp == nil || len(resp.CenBandwidthPackages.CenBandwidthPackage) < 1 {
			break
		}
		insts = append(insts, resp.CenBandwidthPackages.CenBandwidthPackage...)

		if len(resp.CenBandwidthPackages.CenBandwidthPackage) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}

	sweeped := false
	for _, v := range insts {
		name := v.Name
		id := v.CenBandwidthPackageId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping CEN bandwidth package: %s (%s)", name, id)
			continue
		}
		sweeped = true
		if v.Status == string(InUse) {
			log.Printf("[INFO] Deleting CEN bandwidth package attachment: %s (%s)", name, id)
			req := cbn.CreateUnassociateCenBandwidthPackageRequest()
			req.CenId = v.CenIds.CenId[0]
			req.CenBandwidthPackageId = id
			_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
				return cbnClient.UnassociateCenBandwidthPackage(req)
			})
			if err != nil {
				log.Printf("[ERROR] Failed to delete CEN bandwidth package attachment (%s (%s)): %s", name, id, err)
			}
		}
		log.Printf("[INFO] Deleting CEN bandwidth package: %s (%s)", name, id)
		req := cbn.CreateDeleteCenBandwidthPackageRequest()
		req.CenBandwidthPackageId = id
		_, err := client.WithCenClient(func(cbnClient *cbn.Client) (interface{}, error) {
			return cbnClient.DeleteCenBandwidthPackage(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete CEN bandwidth package (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		// Waiting 5 seconds to eusure these instances have been deleted.
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudCenBandwidthPackage_basic(t *testing.T) {
	var cenBwp cbn.CenBandwidthPackage

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_bandwidth_package.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenBandwidthPackageDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenBandwidthPackageConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidth_package.foo", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "name", "tf-testAccCenBandwidthPackageConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "bandwidth", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "geographic_region_ids.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.foo", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.foo", "status"),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackage_update(t *testing.T) {
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
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidth_package.foo", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "name", "tf-testAccCenBandwidthPackageConfig"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "bandwidth", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "geographic_region_ids.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.foo", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.foo", "status"),
				),
			},
			resource.TestStep{
				Config: testAccCenBandwidthPackageConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidth_package.foo", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "name", "tf-testAccCenBandwidthPackageConfigUpdate"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "bandwidth", "10"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.foo", "geographic_region_ids.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.foo", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.foo", "status"),
				),
			},
		},
	})
}

func TestAccAlicloudCenBandwidthPackage_muti(t *testing.T) {
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
					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidth_package.bar1", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.bar1", "name", "tf-testAccCenBandwidthPackageConfigMulti"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.bar1", "bandwidth", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.bar1", "geographic_region_ids.#", "2"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "Asia-Pacific"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.bar1", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.bar1", "status"),

					testAccCheckCenBandwidthPackageExists("alicloud_cen_bandwidth_package.bar2", &cenBwp),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.bar2", "name", "tf-testAccCenBandwidthPackageConfigMulti"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.bar2", "bandwidth", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_cen_bandwidth_package.bar2", "geographic_region_ids.#", "1"),
					testAccCheckCenBandwidthPackageRegionId(&cenBwp, "China", "China"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.bar2", "expired_time"),
					resource.TestCheckResourceAttrSet(
						"alicloud_cen_bandwidth_package.bar2", "status"),
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
			return fmt.Errorf("No CEN bandwidth package ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cenService := CenService{client}
		instance, err := cenService.DescribeCenBandwidthPackage(rs.Primary.ID)
		if err != nil {
			return err
		}

		*cenBwp = instance
		return nil
	}
}

func testAccCheckCenBandwidthPackageDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cenService := CenService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_bandwidth_package" {
			continue
		}

		// Try to find the CEN
		instance, err := cenService.DescribeCenBandwidthPackage(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		return fmt.Errorf("CEN Bandwidth Package %s still exist", instance.CenBandwidthPackageId)
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
resource "alicloud_cen_bandwidth_package" "foo" {
    name = "tf-testAccCenBandwidthPackageConfig"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}
`

const testAccCenBandwidthPackageConfigUpdate = `
resource "alicloud_cen_bandwidth_package" "foo" {
    name = "tf-testAccCenBandwidthPackageConfigUpdate"
    bandwidth = 10
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}
`

const testAccCenBandwidthPackageConfigMulti = `
resource "alicloud_cen_bandwidth_package" "bar1" {
    name = "tf-testAccCenBandwidthPackageConfigMulti"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"Asia-Pacific"]
}

resource "alicloud_cen_bandwidth_package" "bar2" {
    name = "tf-testAccCenBandwidthPackageConfigMulti"
    bandwidth = 5
    geographic_region_ids = [
		"China",
		"China"]
}
`
