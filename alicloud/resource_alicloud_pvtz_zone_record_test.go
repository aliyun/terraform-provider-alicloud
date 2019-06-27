package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPvtzZoneRecord_update(t *testing.T) {
	var record pvtz.Record
	rand := acctest.RandIntRange(10000, 999999)
	resourceId := "alicloud_pvtz_zone_record.foo"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneRecordConfig(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			}, {
				Config: testAccPvtzZoneRecordConfigUpdateType(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "TXT"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			{
				Config: testAccPvtzZoneRecordConfigUpdateValue(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.3"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "TXT"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			{
				Config: testAccPvtzZoneRecordConfigUpdateTTl(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.3"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "TXT"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "30"),
				),
			},
			{
				Config: testAccPvtzZoneRecordConfigUpdatePriority(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.3"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "TXT"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "30"),
				),
			},
			{
				Config: testAccPvtzZoneRecordConfigUpdateAll(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "bbb.test.com"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "10"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_multi(t *testing.T) {
	var record pvtz.Record
	rand := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPvtzZoneRecordConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo.1", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.foo.1", "resource_record", "aaa"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.foo.1", "type", "A"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.foo.1", "priority", "0"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.foo.1", "value", "2.2.2.2"),
				),
			},
		},
	})
}

func testAccAlicloudPvtzZoneRecordExists(n string, record *pvtz.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return WrapError(fmt.Errorf("Not found: %s", n))
		}

		if rs.Primary.ID == "" {
			return WrapError(fmt.Errorf("No Record ID is set"))
		}

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		pvtzService := PvtzService{client}

		instance, err := pvtzService.DescribeZoneRecord(rs.Primary.ID)

		if err != nil {
			return WrapError(err)
		}

		*record = instance
		return nil
	}
}

func testAccAlicloudPvtzZoneRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	pvtzService := PvtzService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_pvtz_zone_record" {
			continue
		}
		zoneRecord, err := pvtzService.DescribeZoneRecord(rs.Primary.ID)

		if err != nil && !NotFoundError(err) {
			return WrapError(err)
		}

		if zoneRecord.Rr != "" {
			return WrapError(fmt.Errorf("recordId %s still exist", zoneRecord.Rr))
		}
	}

	return nil
}

func testAccPvtzZoneRecordConfig(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "www"
		type = "A"
		value = "2.2.2.2"
		ttl = "60"
	}
	`, rand)
}

func testAccPvtzZoneRecordConfigUpdateType(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "www"
		type = "TXT"
		value = "2.2.2.2"
		ttl = "60"
	}
	`, rand)
}
func testAccPvtzZoneRecordConfigUpdateValue(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "www"
		type = "TXT"
		value = "2.2.2.3"
		ttl = "60"
	}
	`, rand)
}

func testAccPvtzZoneRecordConfigUpdatePriority(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "www"
		type = "TXT"
		value = "2.2.2.3"
		priority = "50"
		ttl = "30"
	}
	`, rand)
}

func testAccPvtzZoneRecordConfigUpdateTTl(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "www"
		type = "TXT"
		value = "2.2.2.3"
		ttl = "30"
	}
	`, rand)
}

func testAccPvtzZoneRecordConfigUpdateAll(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "www"
		type = "MX"
		value = "bbb.test.com"
		priority = "10"
		ttl = "60"
	}
	`, rand)
}

func testAccPvtzZoneRecordConfigMulti(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		count="2"
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "aaa"
		type = "A"
		value = "2.2.2.2"
		priority = "10"
	}
	`, rand)
}
