package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudExpressConnectPhysicalConnectionServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_express_connect_physical_connection_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudPhysicalConnectionServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudPhysicalConnectionServiceDataSource = `
data "alicloud_express_connect_physical_connection_service" "current" {
	enable = "On"
}
`
