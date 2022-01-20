package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func Test_dataSourceAlicloudCdnIpInfo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCdnIpInfosDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_cdn_ip_info.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_ip_info.current", "id"),
					resource.TestCheckResourceAttrSet("data.alicloud_cdn_ip_info.current", "cdn_ip"),
					resource.TestCheckResourceAttr("data.alicloud_cdn_ip_info.current", "cdn_ip", "False"),
				),
			},
		},
	})
}

const testAccCheckAlicloudCdnIpInfosDataSource = `
data "alicloud_cdn_ip_info" "current" {
	ip = "114.114.114.114"
}
`
