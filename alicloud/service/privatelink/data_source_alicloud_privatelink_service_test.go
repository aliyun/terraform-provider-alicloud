package privatelink_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudPrivateLinkServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_privatelink_service.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAliCloudPrivateLinkServiceDataSourceNil,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "PrivateLinkServiceHasNotBeenOpened",
						"status": "",
						"enable": "Off",
					}),
				),
			},
			{
				Config: testAccCheckAliCloudPrivateLinkServiceDataSourceWithOff,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     "PrivateLinkServiceHasNotBeenOpened",
						"status": "",
						"enable": "Off",
					}),
				),
			},
			{
				Config: testAccCheckAliCloudPrivateLinkServiceDataSourceWithOn,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     acctest.CHECKSET,
						"status": "Opened",
						"enable": "On",
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
