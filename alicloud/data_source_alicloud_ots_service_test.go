package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudOtsServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOtsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_ots_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_ots_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_ots_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudOtsServiceDataSource = `
data "alicloud_ots_service" "current" {
	enable = "On"
}
`
