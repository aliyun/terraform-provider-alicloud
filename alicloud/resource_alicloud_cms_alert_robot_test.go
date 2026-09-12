package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsAlertRobot_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_robot.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsAlertRobotMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRobot")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlertRobotBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_robot_name":      "${var.name}",
					"type":                  "DING",
					"url":                   "https://oapi.dingtalk.com/robot/send?access_token=testacc",
					"lang":                  "zh_CN",
					"workspace":             "test-workspace",
					"digital_employee_name": "test-employee",
					"robot_sign_key":        "test-sign-key",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_robot_name":      name,
						"type":                  "DING",
						"url":                   "https://oapi.dingtalk.com/robot/send?access_token=testacc",
						"lang":                  "zh_CN",
						"workspace":             "test-workspace",
						"digital_employee_name": "test-employee",
						"robot_sign_key":        "test-sign-key",
						"alert_robot_id":        CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_robot_name":      "${var.name}_update",
					"type":                  "DING",
					"url":                   "https://oapi.dingtalk.com/robot/send?access_token=updated",
					"lang":                  "en_US",
					"workspace":             "test-workspace",
					"digital_employee_name": "updated-employee",
					"robot_sign_key":        "updated-sign-key",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_robot_name":      name + "_update",
						"url":                   "https://oapi.dingtalk.com/robot/send?access_token=updated",
						"lang":                  "en_US",
						"digital_employee_name": "updated-employee",
						"robot_sign_key":        "updated-sign-key",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"url"},
			},
		},
	})
}

var AliCloudCmsAlertRobotMap = map[string]string{
	"alert_robot_id": CHECKSET,
}

func AliCloudCmsAlertRobotBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
`, name)
}
