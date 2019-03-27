package alicloud

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/acctest"
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

func TestAccAlicloudPvtzZone_update(t *testing.T) {
	var zone pvtz.DescribeZoneInfoResponse
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.foo", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.foo", "name", fmt.Sprintf("tf-testacc%d.test.com", rand)),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone.foo", "name", fmt.Sprintf("tf-testacc%d.test.com", rand)),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "id"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone.foo", "remark", ""),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "creation_time"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "update_time"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "is_ptr"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "record_count"),
				),
			},
			{
				Config: testAccPvtzZoneConfigUpdate(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.foo", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.foo", "name", fmt.Sprintf("tf-testacc%d.test.com", rand)),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone.foo", "name", fmt.Sprintf("tf-testacc%d.test.com", rand)),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "id"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.foo", "remark", "remark-test"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "creation_time"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "update_time"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "is_ptr"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "record_count"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZone_multi(t *testing.T) {
	var zone pvtz.DescribeZoneInfoResponse
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.foo.4", &zone),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone.foo.4", "name", fmt.Sprintf("tf-testacc%d.test.com", rand)),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo.4", "id"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone.foo.4", "remark", ""),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo.4", "creation_time"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo.4", "update_time"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo.4", "is_ptr"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo.4", "record_count"),
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

		instance, err := pvtzService.DescribePvtzZone(rs.Primary.ID)

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

		instance, err := pvtzService.DescribePvtzZone(rs.Primary.ID)

		if err != nil && !NotFoundError(err) {
			return err
		}

		if instance.ZoneId != "" {
			return fmt.Errorf("zone %s still exist", instance.ZoneId)
		}
	}

	return nil
}

func testAccPvtzZoneConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "foo" {
		name = "tf-testacc%d.test.com"
	}
	`, rand)
}
func testAccPvtzZoneConfigUpdate(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "foo" {
		name = "tf-testacc%d.test.com"
		remark = "remark-test"
	}
	`, rand)
}

func testAccPvtzZoneConfigMulti(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "foo" {
		count = 5
		name = "tf-testacc%d.test.com"
	}
	`, rand)
}
