package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudPvtzZone_Basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone.foo", "name", "foo.test.com"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone.foo", "id"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZone_update(t *testing.T) {
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
						"alicloud_pvtz_zone.foo", "name", "foo.test.com"),
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
						"alicloud_pvtz_zone.bar_1", "name", "bar1.test.com"),
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.bar_2", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.bar_2", "name", "bar2.test.com"),
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.bar_3", &zone),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone.bar_3", "name", "bar3.test.com"),
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

		client := testAccProvider.Meta().(*AliyunClient)

		instance, err := client.DescribePvtzZoneInfo(rs.Primary.ID)

		if err != nil {
			return err
		}

		*zone = instance
		return nil
	}
}

func testAccAlicloudPvtzZoneDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pvtz_zone" {
			continue
		}

		instance, err := client.DescribePvtzZoneInfo(rs.Primary.ID)

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
	name = "foo.test.com"
}
`
const testAccPvtzZoneConfigUpdate = `
resource "alicloud_pvtz_zone" "foo" {
	name = "foo.test.com"
	remark = "remark-test"
}
`

const testAccPvtzZoneConfigMulti = `
resource "alicloud_pvtz_zone" "bar_1" {
	name = "bar1.test.com"
}
resource "alicloud_pvtz_zone" "bar_2" {
	name = "bar2.test.com"
}
resource "alicloud_pvtz_zone" "bar_3" {
	name = "bar3.test.com"
}
`
