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
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zone_records.keyword", "records.0.value", "2.2.2.2"),
				),
			},
		},
	})
}

const testAccCheckAlicloudPvtzZoneRecordsDataSourceBasic = `
resource "alicloud_pvtz_zone" "basic" {
	name = "basic.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.basic.id}"
	resource_record = "www"
	type = "A"
	value = "2.2.2.2"
	priority = "5"
	ttl = "60"
}

data "alicloud_pvtz_zone_records" "keyword" {
	zone_id = "${alicloud_pvtz_zone.basic.id}"
	keyword = "${alicloud_pvtz_zone_record.foo.value}"
}
`
