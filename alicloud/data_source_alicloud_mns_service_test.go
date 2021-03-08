package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMnsServiceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudMnsServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_mns_service.current"),
					resource.TestCheckResourceAttrSet("data.alicloud_mns_service.current", "id"),
					resource.TestCheckResourceAttr("data.alicloud_mns_service.current", "status", "Opened"),
				),
			},
		},
	})
}

const testAccCheckAlicloudMnsServiceDataSource = `
data "alicloud_mns_service" "current" {
	enable = "On"
}
`
