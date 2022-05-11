package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSAEServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudSaeServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_sae_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_sae_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_sae_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudSaeServiceDataSource = `
data "alicloud_sae_service" "current" {
	enable = "On"
}
`
