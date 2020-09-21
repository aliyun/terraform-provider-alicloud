package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudEdasServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudEdasServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_edas_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_edas_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_edas_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudEdasServiceDataSource = `
data "alicloud_edas_service" "current" {
	enable = "On"
}
`
