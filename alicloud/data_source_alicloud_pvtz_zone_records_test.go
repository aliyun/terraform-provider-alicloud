package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudPvtzZoneRecordsDataSource_basic(t *testing.T) {
	var pvtzZoneRecord pvtz.Record

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzZoneRecordsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneRecordExists("alicloud_pvtz_zone_record.foo", &pvtzZoneRecord),
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zone_records.keyword"),
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

func TestAccAlicloudPvtzZoneRecordsDataSource_empty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzZoneRecordsDataSourceEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zone_records.keyword"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.id"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.resource_record"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.type"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.ttl"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.priority"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.value"),
					resource.TestCheckNoResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.status"),
				),
			},
		},
	})
}

const testAccCheckAlicloudPvtzZoneRecordsDataSourceBasic = `
resource "alicloud_pvtz_zone" "basic" {
	name = "tf-testacc.test.com"
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
`

const testAccCheckAlicloudPvtzZoneRecordsDataSourceEmpty = `
resource "alicloud_pvtz_zone" "basic" {
	name = "tf-testacc.test.com"
}

data "alicloud_pvtz_zone_records" "keyword" {
	zone_id = "${alicloud_pvtz_zone.basic.id}"
	keyword = "tf-testacc-fake-name"
}
`
