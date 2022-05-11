package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudNASServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudNasServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_nas_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_nas_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_nas_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudNasServiceDataSource = `
data "alicloud_nas_service" "current" {
	enable = "On"
}
`
