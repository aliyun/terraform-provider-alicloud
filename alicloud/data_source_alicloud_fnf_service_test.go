package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudFnfServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudFnfServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_fnf_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_fnf_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_fnf_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudFnfServiceDataSource = `
data "alicloud_fnf_service" "current" {
	enable = "On"
}
`
