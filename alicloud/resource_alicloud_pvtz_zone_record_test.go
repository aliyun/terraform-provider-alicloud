package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

// Only MX supports priority

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
			{
				Config: testAccPvtzZoneRecordConfig(acctest.RandIntRange(10000, 999999)),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_updateRr(t *testing.T) {
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
				Config: testAccPvtzZoneRecordConfigUpdateResourceRecord(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "@"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_updateType(t *testing.T) {
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
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_updateValue(t *testing.T) {
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
				Config: testAccPvtzZoneRecordConfigUpdateValue(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.3"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}
func TestAccAlicloudPvtzZoneRecord_updatePriority(t *testing.T) {
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
				Config: testAccPvtzZoneRecordConfigPriority(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "aaa.test.com"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "10"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			{
				Config: testAccPvtzZoneRecordConfigUpdatePriority(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "aaa.test.com"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "20"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}
func TestAccAlicloudPvtzZoneRecord_updateTTL(t *testing.T) {
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
				Config: testAccPvtzZoneRecordConfigUpdateTTl(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "0"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "30"),
				),
			},
		},
	})

}
func TestAccAlicloudPvtzZoneRecord_updateAll(t *testing.T) {
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
				Config: testAccPvtzZoneRecordConfigUpdateAll(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "@"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "bbb.test.com"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "MX"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "10"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "30"),
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
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_1", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "resource_record", "aaa"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "type", "A"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "priority", "0"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "value", "2.2.2.2"),
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_2", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "resource_record", "bbb"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "priority", "0"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "value", "c.test.com"),
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_3", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_3", "resource_record", "ccc"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_3", "type", "A"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_3", "priority", "0"),
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
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		pvtzService := PvtzService{client}

		instance, err := pvtzService.DescribeZoneRecord(recordId, zoneId)

		if err != nil {
			return err
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
		recordIdStr, zoneId, _ := splitRecordIdAndZoneId(rs.Primary.ID)
		recordId, err := strconv.Atoi(recordIdStr)
		if err != nil {
			return err
		}

		zoneRecord, err := pvtzService.DescribeZoneRecord(recordId, zoneId)

		if err != nil && !NotFoundError(err) {
			return err
		}

		if zoneRecord.Rr != "" {
			return fmt.Errorf("recordId %s still exist", zoneRecord.Rr)
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
func testAccPvtzZoneRecordConfigUpdateResourceRecord(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "@"
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
		type = "A"
		value = "2.2.2.3"
		ttl = "60"
	}
	`, rand)
}

func testAccPvtzZoneRecordConfigPriority(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "www"
		type = "MX"
		value = "aaa.test.com"
		priority = "10"
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
		type = "MX"
		value = "aaa.test.com"
		priority = "20"
		ttl = "60"
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
		type = "A"
		value = "2.2.2.2"
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
		resource_record = "@"
		type = "MX"
		value = "bbb.test.com"
		priority = "10"
		ttl = "30"
	}
	`, rand)
}

func testAccPvtzZoneRecordConfigMulti(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "zone" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "bar_1" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "aaa"
		type = "A"
		value = "2.2.2.2"
		priority = "10"
	}
	resource "alicloud_pvtz_zone_record" "bar_2" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "bbb"
		type = "CNAME"
		value = "c.test.com"
		priority = "5"
	}
	resource "alicloud_pvtz_zone_record" "bar_3" {
		zone_id = "${alicloud_pvtz_zone.zone.id}"
		resource_record = "ccc"
		type = "A"
		value = "3.3.3.3"
		priority = "3"
	}
	`, rand)
}
