package alicloud

import (
	"testing"

	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudPvtzZoneRecordsDataSource_basic(t *testing.T) {
	var pvtzZoneRecord pvtz.Record
	num := acctest.RandIntRange(10000, 999999)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzZoneRecordsDataSource_basic(num),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &pvtzZoneRecord),
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zone_records.keyword"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_pvtz_zone_records.keyword", "records.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.resource_record", "www"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.type", "A"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.ttl", "60"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.value", "2.2.2.2"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.status", "ENABLE"),
				),
			},
		},
	})
}

func TestAccAlicloudPvtzZoneRecordsDataSource_keyword(t *testing.T) {

	num := acctest.RandIntRange(10000, 999999)
	var pvtzZoneRecord pvtz.Record

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzZoneRecordsDataSource_keyword(num),
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &pvtzZoneRecord),
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zone_records.keyword"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_pvtz_zone_records.keyword", "records.0.id"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.resource_record", "www"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.type", "A"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.ttl", "60"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.priority", "0"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.value", "2.2.2.2"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.status", "ENABLE"),
				),
			},
			{
				Config: testAccCheckAlicloudPvtzZoneRecordsDataSource_keyword_empty(num),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zone_records.keyword"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.#", "0"),
				),
			},
		},
	})
}

func testAccCheckAlicloudPvtzZoneRecordsDataSource_basic(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "basic" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.basic.id}"
		resource_record = "www"
		type = "A"
		value = "2.2.2.2"
		ttl = "60"
	}

	data "alicloud_pvtz_zone_records" "keyword" {
		zone_id = "${alicloud_pvtz_zone_record.foo.zone_id}"
	}
	`, rand)
}

func testAccCheckAlicloudPvtzZoneRecordsDataSource_keyword(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "basic" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.basic.id}"
		resource_record = "www"
		type = "A"
		value = "2.2.2.2"
		ttl = "60"
	}

	data "alicloud_pvtz_zone_records" "keyword" {
		zone_id = "${alicloud_pvtz_zone.basic.id}"
		keyword = "${alicloud_pvtz_zone_record.foo.value}"
	}
	`, rand)
}

func testAccCheckAlicloudPvtzZoneRecordsDataSource_keyword_empty(rand int) string {
	return fmt.Sprintf(`
	resource "alicloud_pvtz_zone" "basic" {
		name = "tf-testacc%d.test.com"
	}

	resource "alicloud_pvtz_zone_record" "foo" {
		zone_id = "${alicloud_pvtz_zone.basic.id}"
		resource_record = "www"
		type = "A"
		value = "2.2.2.2"
		ttl = "60"
	}

	data "alicloud_pvtz_zone_records" "keyword" {
		zone_id = "${alicloud_pvtz_zone.basic.id}"
		keyword = "${alicloud_pvtz_zone_record.foo.value}-fake"
	}
	`, rand)
}
