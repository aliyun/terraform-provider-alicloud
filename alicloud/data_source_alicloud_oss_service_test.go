package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudOssServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOssServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_oss_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_oss_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudOssServiceDataSource = `
data "alicloud_oss_service" "current" {
	enable = "On"
}
`
