package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDataWorksServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudDataWorksServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_dataworks_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_dataworks_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_dataworks_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudDataWorksServiceDataSource = `
data "alicloud_dataworks_service" "current" {
	enable = "On"
}
`
