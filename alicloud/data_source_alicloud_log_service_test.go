package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudLogServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_log_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_log_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_log_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudLogServiceDataSource = `
data "alicloud_log_service" "current" {
	enable = "On"
}
`
