package alicloud

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudSSOServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cloud_sso_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	checkoutSupportedRegions(t, true, connectivity.CloudSsoSupportRegions)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckEnterpriseAccountEnabled(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCloudSsoServiceDataOnSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "Opened",
					}),
				),
			},
			{
				Config: testAccCheckAlicloudCloudSsoServiceOffDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":     CHECKSET,
						"status": "",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudCloudSsoServiceDataOnSource = `
data "alicloud_cloud_sso_service" "current" {
	enable = "On"
}
`

const testAccCheckAlicloudCloudSsoServiceOffDataSource = `
data "alicloud_cloud_sso_service" "current" {
	enable = "Off"
}
`
