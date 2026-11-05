package cdn_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlicloudCDNServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCdnServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID("data.alicloud_cdn_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_cdn_service.current", "status", "Opened"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "internet_charge_type"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "opening_time"),
					//resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "changing_charge_type"),
					//resource.TestCheckResourceAttrSet("data.alicloud_cdn_service.current", "changing_affect_time"),
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
