package cdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceAlicloudCdnIpInfo(t *testing.T) {
	dataSourceCdnIpInfoResourceId := "data.alicloud_cdn_ip_info.test"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCdnIpInfosDataSource(),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckAlicloudDataSourceID(dataSourceCdnIpInfoResourceId),
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
