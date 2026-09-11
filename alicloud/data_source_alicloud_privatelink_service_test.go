package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudPrivateLinkServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_service.default"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudPrivateLinkServiceDataSourceNil,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "PrivateLinkServiceHasNotBeenOpened",
						"status": "",
					}),
				),
			},
			{
				Config: testAccCheckAliCloudPrivateLinkServiceDataSourceWithOff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "PrivateLinkServiceHasNotBeenOpened",
						"status": "",
					}),
				),
			},
			{
				Config: testAccCheckAliCloudPrivateLinkServiceDataSourceWithOn,
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

const testAccCheckAliCloudPrivateLinkServiceDataSourceNil = `
data "alicloud_privatelink_service" "default" {
}
`

const testAccCheckAliCloudPrivateLinkServiceDataSourceWithOff = `
data "alicloud_privatelink_service" "default" {
  enable = "Off"
}
`

const testAccCheckAliCloudPrivateLinkServiceDataSourceWithOn = `
data "alicloud_privatelink_service" "default" {
  enable = "On"
}
`
