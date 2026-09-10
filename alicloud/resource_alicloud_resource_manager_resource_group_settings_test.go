package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// TestAccAliCloudResourceManagerResourceGroupSettings_basic0 covers the account-level
// ResourceGroupSettings singleton. It exercises both attributes (admin setting +
// notification setting) through create and every update branch (Enable/Disable
// notification, UpdateResourceGroupAdminSetting), then re-imports the state.
func TestAccAliCloudResourceManagerResourceGroupSettings_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_resource_manager_resource_group_settings.default"
	ra := resourceAttrInit(resourceId, AlicloudResourceManagerResourceGroupSettingsMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ResourceManagerServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeResourceManagerResourceGroupSettings")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccresourcemanager%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudResourceManagerResourceGroupSettingsBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				// Create: both settings enabled. Covers the Required admin field and
				// the Optional notification field (100% attribute coverage).
				Config: testAccConfig(map[string]interface{}{
					"resource_group_admin_setting_status":        true,
					"resource_group_notification_setting_status": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_admin_setting_status":        "true",
						"resource_group_notification_setting_status": "true",
					}),
				),
			},
			{
				// Update: disable notification while keeping admin enabled. Exercises
				// the DisableResourceGroupNotification update branch.
				Config: testAccConfig(map[string]interface{}{
					"resource_group_notification_setting_status": false,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_admin_setting_status":        "true",
						"resource_group_notification_setting_status": "false",
					}),
				),
			},
			{
				// Update: switch admin to false and re-enable notification. Exercises
				// the UpdateResourceGroupAdminSetting update branch and the
				// EnableResourceGroupNotification update branch.
				Config: testAccConfig(map[string]interface{}{
					"resource_group_admin_setting_status":        false,
					"resource_group_notification_setting_status": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_admin_setting_status":        "false",
						"resource_group_notification_setting_status": "true",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudResourceManagerResourceGroupSettingsMap0 = map[string]string{
	"resource_group_admin_setting_status":        CHECKSET,
	"resource_group_notification_setting_status": CHECKSET,
}

func AlicloudResourceManagerResourceGroupSettingsBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
`, name)
}
