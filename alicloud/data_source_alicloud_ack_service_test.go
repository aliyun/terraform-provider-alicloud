package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudAckServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_ack_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudAckServiceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
						"type":   "propayasgo",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudAckServiceDataSource = `
data "alicloud_ack_service" "current" {
	enable = "On"
	type   = "propayasgo"
}
`
