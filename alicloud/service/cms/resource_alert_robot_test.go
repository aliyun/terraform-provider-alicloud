package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccAliCloudCmsAlertRobot_basic exercises the full Framework CRUD path
// through the muxed provider: create with defaults, in-place update of every
// updatable field, and import verification (url is write-only at the API
// level and ignored on import).
func TestAccAliCloudCmsAlertRobot_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlertRobotConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_robot.default", "id"),
					resource.TestCheckResourceAttrSet("alicloud_cms_alert_robot.default", "alert_robot_id"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "alert_robot_name", fmt.Sprintf("tf-acc-cms-%d", r)),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "type", "DING"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "url", "https://oapi.dingtalk.com/robot/send?access_token=tf-acc"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "lang", "zh_CN"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "workspace", "tf-acc-workspace"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "digital_employee_name", "tf-acc-employee"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "robot_sign_key", "tf-acc-sign-key"),
				),
			},
			{
				Config: testAccCmsAlertRobotConfigUpdate(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "alert_robot_name", fmt.Sprintf("tf-acc-cms-%d-upd", r)),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "url", "https://oapi.dingtalk.com/robot/send?access_token=tf-acc-upd"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "lang", "en_US"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "digital_employee_name", "tf-acc-employee-upd"),
					resource.TestCheckResourceAttr("alicloud_cms_alert_robot.default", "robot_sign_key", "tf-acc-sign-key-upd"),
				),
			},
			{
				ResourceName:            "alicloud_cms_alert_robot.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"url"},
			},
		},
	})
}

func testAccCmsAlertRobotConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_robot" "default" {
  alert_robot_name      = "tf-acc-cms-%[1]d"
  type                  = "DING"
  url                   = "https://oapi.dingtalk.com/robot/send?access_token=tf-acc"
  workspace             = "tf-acc-workspace"
  digital_employee_name = "tf-acc-employee"
  robot_sign_key        = "tf-acc-sign-key"
}
`, r)
}

func testAccCmsAlertRobotConfigUpdate(r int) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_robot" "default" {
  alert_robot_name      = "tf-acc-cms-%[1]d-upd"
  type                  = "DING"
  url                   = "https://oapi.dingtalk.com/robot/send?access_token=tf-acc-upd"
  lang                  = "en_US"
  workspace             = "tf-acc-workspace"
  digital_employee_name = "tf-acc-employee-upd"
  robot_sign_key        = "tf-acc-sign-key-upd"
}
`, r)
}
