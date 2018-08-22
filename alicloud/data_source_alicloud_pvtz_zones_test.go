package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/pvtz"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudPvtzZonesDataSource_basic(t *testing.T) {
	var pvtzZone pvtz.DescribeZoneInfoResponse

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccAlicloudPvtzZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPvtzZoneDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccAlicloudPvtzZoneExists("alicloud_pvtz_zone.basic", &pvtzZone),
					testAccCheckAlicloudDataSourceID("data.alicloud_pvtz_zones.keyword"),
					resource.TestCheckResourceAttr("data.alicloud_pvtz_zones.keyword", "zones.0.name", "basic.test.com"),
				),
			},
		},
	})
}

const testAccCheckAlicloudPvtzZoneDataSourceBasic = `
resource "alicloud_pvtz_zone" "basic" {
	name = "basic.test.com"
}
data "alicloud_pvtz_zones" "keyword" {
	keyword = "${alicloud_pvtz_zone.basic.name}"
}
`
