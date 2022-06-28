package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCDNServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCdnServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cdn_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_cdn_service.current", "status", "Opened"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "internet_charge_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "opening_time"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "changing_charge_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "changing_affect_time"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCdnServiceDataSource = `
data "alicloud_cdn_service" "current" {
	enable = "On"
	internet_charge_type = "PayByTraffic"
}
`
