package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAlicloudCdnIpInfo(t *testing.T) {
	dataSourceCdnIpInfoResourceId := "data.alicloud_cdn_ip_info.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCdnIpInfosDataSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID(dataSourceCdnIpInfoResourceId),
					resource.TestCheckResourceAttrSet(dataSourceCdnIpInfoResourceId, "id"),
					resource.TestCheckResourceAttrSet(dataSourceCdnIpInfoResourceId, "cdn_ip"),
					resource.TestCheckResourceAttr(dataSourceCdnIpInfoResourceId, "cdn_ip", "False"),
				),
			},
		},
	})
}

func testAccCheckAlicloudCdnIpInfosDataSource() string {
	return fmt.Sprintf(`
data "alicloud_cdn_ip_info" "test" {
  ip = "114.114.114.114"
}
`)
}
