package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCloudCloudMonitorServiceEnterpriseServiceDataSource(t *testing.T) {
	resourceId := "data.alicloud_cloud_monitor_service_enterprise_service.current"
	testAccCheck := resourceAttrInit(resourceId, map[string]string{}).resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudCloudCloudMonitorServiceEnterpriseServiceDataOnSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":          CHECKSET,
						"create_time": CHECKSET,
						"status":      "Opened",
					}),
				),
			},
			{
				Config: testAccCheckAlicloudCloudCloudMonitorServiceEnterpriseServiceOffDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"id":          CHECKSET,
						"create_time": "",
						"status":      "",
					}),
				),
			},
		},
	})
}

const testAccCheckAlicloudCloudCloudMonitorServiceEnterpriseServiceDataOnSource = `
data "alicloud_cloud_monitor_service_enterprise_service" "current" {
	enable = "On"
}
`

const testAccCheckAlicloudCloudCloudMonitorServiceEnterpriseServiceOffDataSource = `
data "alicloud_cloud_monitor_service_enterprise_service" "current" {
	enable = "Off"
}
`
