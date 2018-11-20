package alicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudPvtzZoneRecord_Basic(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_updateRr(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigUpdateResourceRecord,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "@"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_updateType(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigUpdateType,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "TXT"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}

func TestAccAlicloudPvtzZoneRecord_updateValue(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigUpdateValue,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.3"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}
func TestAccAlicloudPvtzZoneRecord_updatePriority(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigUpdatePriority,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "10"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
		},
	})

}
func TestAccAlicloudPvtzZoneRecord_updateTTL(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigUpdateTTl,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "30"),
				),
			},
		},
	})

}
func TestAccAlicloudPvtzZoneRecord_updateAll(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "www"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "2.2.2.2"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "A"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "5"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "60"),
				),
			},
			resource.TestStep{
				Config: testAccPvtzZoneRecordConfigUpdateAll,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &record),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "resource_record", "@"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "value", "bbb.test.com"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "type", "CNAME"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "priority", "10"),
					resource.TestCheckResourceAttr("alicloud_pvtz_zone_record.foo", "ttl", "30"),
				),
			},
		},
	})

}
func TestAccAlicloudPvtzZoneRecord_multi(t *testing.T) {
	if !isRegionSupports(PrivateZone) {
		logTestSkippedBecauseOfUnsupportedRegionalFeatures(t.Name(), PrivateZone)
		return
	}

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
						"alicloud_pvtz_zone_record.bar_1", "resource_record", "aaa"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "type", "A"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "priority", "10"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_1", "value", "2.2.2.2"),
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_2", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "resource_record", "bbb"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "type", "CNAME"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "priority", "5"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_2", "value", "c.test.com"),
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.bar_3", &record),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_3", "resource_record", "ccc"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_3", "type", "A"),
					resource.TestCheckResourceAttr(
						"alicloud_pvtz_zone_record.bar_3", "priority", "3"),
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

const testAccPvtzZoneRecordConfig = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
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
const testAccPvtzZoneRecordConfigUpdateResourceRecord = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "@"
	type = "A"
	value = "2.2.2.2"
	priority = "5"
	ttl = "60"
}
`
const testAccPvtzZoneRecordConfigUpdateType = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "www"
	type = "TXT"
	value = "2.2.2.2"
	priority = "5"
	ttl = "60"
}
`

const testAccPvtzZoneRecordConfigUpdateValue = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "www"
	type = "A"
	value = "2.2.2.3"
	priority = "5"
	ttl = "60"
}
`

const testAccPvtzZoneRecordConfigUpdatePriority = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "www"
	type = "A"
	value = "2.2.2.2"
	priority = "10"
	ttl = "60"
}
`

const testAccPvtzZoneRecordConfigUpdateTTl = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "www"
	type = "A"
	value = "2.2.2.2"
	priority = "5"
	ttl = "30"
}
`

const testAccPvtzZoneRecordConfigUpdateAll = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "@"
	type = "CNAME"
	value = "bbb.test.com"
	priority = "10"
	ttl = "30"
}
`

const testAccPvtzZoneRecordConfigMulti = `
resource "alicloud_pvtz_zone" "zone" {
	name = "tf-testacc.test.com"
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
`
