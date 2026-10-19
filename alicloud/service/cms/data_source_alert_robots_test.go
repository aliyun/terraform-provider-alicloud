package cms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCmsAlertRobotsDataSource_basic(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlertRobotsDataSourceConfig(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_alert_robots.default", "robots.#", "1"),
					resource.TestCheckResourceAttrPair("data.alicloud_cms_alert_robots.default", "robots.0.alert_robot_id", "alicloud_cms_alert_robot.default", "id"),
					resource.TestCheckResourceAttrPair("data.alicloud_cms_alert_robots.default", "robots.0.alert_robot_name", "alicloud_cms_alert_robot.default", "alert_robot_name"),
					resource.TestCheckResourceAttrPair("data.alicloud_cms_alert_robots.default", "robots.0.type", "alicloud_cms_alert_robot.default", "type"),
				),
			},
		},
	})
}

func TestAccAliCloudCmsAlertRobotsDataSource_nameRegex(t *testing.T) {
	r := tfacctest.RandIntRange(10000, 99999)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccCmsAlertRobotsDataSourceConfigNameRegex(r),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.alicloud_cms_alert_robots.default", "robots.#", "1"),
					resource.TestCheckResourceAttrPair("data.alicloud_cms_alert_robots.default", "robots.0.alert_robot_name", "alicloud_cms_alert_robot.default", "alert_robot_name"),
				),
			},
		},
	})
}

func testAccCmsAlertRobotsDataSourceConfig(r int) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_robot" "default" {
  alert_robot_name = "tf-acc-cmsds-%[1]d"
  type             = "DING"
  url              = "https://oapi.dingtalk.com/robot/send?access_token=tf-acc"
}

data "alicloud_cms_alert_robots" "default" {
  ids = [alicloud_cms_alert_robot.default.id]
}
`, r)
}

func testAccCmsAlertRobotsDataSourceConfigNameRegex(r int) string {
	return fmt.Sprintf(`
resource "alicloud_cms_alert_robot" "default" {
  alert_robot_name = "tf-acc-cmsds-%[1]d"
  type             = "DING"
  url              = "https://oapi.dingtalk.com/robot/send?access_token=tf-acc"
}

data "alicloud_cms_alert_robots" "default" {
  name_regex = "^${alicloud_cms_alert_robot.default.alert_robot_name}$"
}
`, r)
}
