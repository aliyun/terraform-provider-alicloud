package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func init() {
	resource.AddTestSweepers("alicloud_pvtz_zone", &resource.Sweeper{
		Name: "alicloud_pvtz_zone",
		F:    testSweepPvtzZones,
	})
}

func testSweepPvtzZones(region string) error {
	rawClient, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	prefixes := []string{
		"tf-testAcc",
		"tf_testAcc",
		"tf_test_",
		"tf-test-",
		"tftest",
	}

	var zones []pvtz.Zone
	req := pvtz.CreateDescribeZonesRequest()
	req.RegionId = client.RegionId
	req.PageSize = requests.NewInteger(PageSizeLarge)
	req.PageNumber = requests.NewInteger(1)
	for {
		raw, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DescribeZones(req)
		})
		if err != nil {
			return fmt.Errorf("Error retrieving Private Zones: %s", err)
		}
		resp, _ := raw.(*pvtz.DescribeZonesResponse)
		if resp == nil || len(resp.Zones.Zone) < 1 {
			break
		}
		zones = append(zones, resp.Zones.Zone...)

		if len(resp.Zones.Zone) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(req.PageNumber); err != nil {
			return err
		} else {
			req.PageNumber = page
		}
	}
	sweeped := false

	for _, v := range zones {
		name := v.ZoneName
		id := v.ZoneId
		skip := true
		for _, prefix := range prefixes {
			if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
				skip = false
				break
			}
		}
		if skip {
			log.Printf("[INFO] Skipping Private Zone: %s (%s)", name, id)
			continue
		}
		sweeped = true
		log.Printf("[INFO] Unbinding VPC from Private Zone: %s (%s)", name, id)
		request := pvtz.CreateBindZoneVpcRequest()
		request.ZoneId = id
		vpcs := make([]pvtz.BindZoneVpcVpcs, 0)
		request.Vpcs = &vpcs

		_, err := client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.BindZoneVpc(request)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to unbind VPC from Private Zone (%s (%s)): %s ", name, id, err)
		}

		log.Printf("[INFO] Deleting Private Zone: %s (%s)", name, id)
		req := pvtz.CreateDeleteZoneRequest()
		req.ZoneId = id
		_, err = client.WithPvtzClient(func(pvtzClient *pvtz.Client) (interface{}, error) {
			return pvtzClient.DeleteZone(req)
		})
		if err != nil {
			log.Printf("[ERROR] Failed to delete Private Zone (%s (%s)): %s", name, id, err)
		}
	}
	if sweeped {
		time.Sleep(5 * time.Second)
	}
	return nil
}

func TestAccAlicloudPvtzZone_Basic(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

	var zone pvtz.DescribeZoneInfoResponse
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_pvtz_zone.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.foo", &zone),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone.foo", "name", "tf-testacc.test.com"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "id"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZone_update(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

	var zone pvtz.DescribeZoneInfoResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.foo", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.foo", "name", "tf-testacc.test.com"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.foo", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.foo", "remark", "remark-test"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZone_multi(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

	var zone pvtz.DescribeZoneInfoResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.bar_1", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.bar_1", "name", "tf-testacc1.test.com"),
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.bar_2", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.bar_2", "name", "tf-testacc2.test.com"),
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.bar_3", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.bar_3", "name", "tf-testacc3.test.com"),
				),
			},
		},
	})
}

func testAccAlicloudPvtzZoneExists(n string, zone *pvtz.DescribeZoneInfoResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ZONE ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		pvtzService := PvtzService{client}

		instance, err := pvtzService.DescribePvtzZoneInfo(rs.Primary.ID)

		if err != nil {
			return err
		}

		*zone = instance
		return nil
	}
}

func testAccAlicloudPvtzZoneDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pvtz_zone" {
			continue
		}

		instance, err := pvtzService.DescribePvtzZoneInfo(rs.Primary.ID)

		if err != nil && !NotFoundError(err) {
			return err
		}

		if instance.ZoneId != "" {
			return fmt.Errorf("zone %s still exist", instance.ZoneId)
		}
	}

	return nil
}

const testAccPvtzZoneConfig = `
resource "alicloud_pvtz_zone" "foo" {
	name = "tf-testacc.test.com"
}
`
const testAccPvtzZoneConfigUpdate = `
resource "alicloud_pvtz_zone" "foo" {
	name = "tf-testacc.test.com"
	remark = "remark-test"
}
`

const testAccPvtzZoneConfigMulti = `
resource "alicloud_pvtz_zone" "bar_1" {
	name = "tf-testacc1.test.com"
}
resource "alicloud_pvtz_zone" "bar_2" {
	name = "tf-testacc2.test.com"
}
resource "alicloud_pvtz_zone" "bar_3" {
	name = "tf-testacc3.test.com"
}
`
