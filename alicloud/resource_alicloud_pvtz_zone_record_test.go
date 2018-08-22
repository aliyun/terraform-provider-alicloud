package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudPvtzZoneRecord_Basic(t *testing.T) {
	var record pvtz.Record
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: "alicloud_pvtz_zone_record.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttrSet("alicloud_pvtz_zone_record.foo", "value"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_update(t *testing.T) {
	var record pvtz.Record

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.foo", "value", "bbb.test.com"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_multi(t *testing.T) {
	var record pvtz.Record

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_1", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "value", "2.2.2.2"),
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_2", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "value", "c.test.com"),
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_3", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_3", "value", "3.3.3.3"),
				),
			},
		},
	})
}

func testAccAlicloudPvtzZoneRecordExists(n string, record *pvtz.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		recordIdStr, zoneId, _ := splitRecordIdAndZoneId(rs.Primary.ID)
		recordId, convErr := strconv.Atoi(recordIdStr)
		if convErr != nil {
			return convErr
		}
		client := testAccProvider.Meta().(*AliyunClient)

		instance, err := client.DescribeZoneRecord(recordId, zoneId)

		if err != nil {
			return err
		}

		*record = instance
		return nil
	}
}

func testAccAlicloudPvtzZoneRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pvtz_zone_record" {
			continue
		}
		recordIdStr, zoneId, _ := splitRecordIdAndZoneId(rs.Primary.ID)
		recordId, err := strconv.Atoi(recordIdStr)
		if err != nil {
			return err
		}

		zoneRecord, err := client.DescribeZoneRecord(recordId, zoneId)

		if err != nil && !NotFoundError(err) {
			return err
		}

		if zoneRecord.Rr != "" {
			return fmt.Errorf("recordId %s still exist", zoneRecord.Rr)
		}
	}

	return nil
}

const testAccPvtzZoneRecordConfig = `
resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "www"
	type = "A"
	value = "2.2.2.2"
	priority = "5"
	ttl = "60"
}
`
const testAccPvtzZoneRecordConfigUpdate = `
resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "www"
	type = "CNAME"
	value = "bbb.test.com"
	priority = "6"
}
`

const testAccPvtzZoneRecordConfigMulti = `
resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_pvtz_zone_record" "bar_1" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "aaa"
	type = "A"
	value = "2.2.2.2"
}
resource "alicloud_pvtz_zone_record" "bar_2" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "bbb"
	type = "CNAME"
	value = "c.test.com"
}
resource "alicloud_pvtz_zone_record" "bar_3" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "ccc"
	type = "A"
	value = "3.3.3.3"
}
`
