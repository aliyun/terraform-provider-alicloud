package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudCmsAlertRobotsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	resourceId := "alicloud_cms_alert_robot.default"
	dataSourceId := "data.alicloud_cms_alert_robots.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlertRobotBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_robot_name": "${var.name}",
					"type":             "DING",
					"url":              "https://oapi.dingtalk.com/robot/send?access_token=testacc",
				}) + testAccCmsAlertRobotsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "robots.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceId, "robots.0.alert_robot_id", resourceId, "id"),
					resource.TestCheckResourceAttrPair(dataSourceId, "robots.0.alert_robot_name", resourceId, "alert_robot_name"),
					resource.TestCheckResourceAttrPair(dataSourceId, "robots.0.type", resourceId, "type"),
				),
			},
		},
	})
}

func TestAccAliCloudCmsAlertRobotsDataSource_nameRegex(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	resourceId := "alicloud_cms_alert_robot.default"
	dataSourceId := "data.alicloud_cms_alert_robots.default"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsAlertRobotBasicDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_robot_name": "${var.name}",
					"type":             "DING",
					"url":              "https://oapi.dingtalk.com/robot/send?access_token=testacc",
				}) + testAccCmsAlertRobotsDataSourceConfigWithNameRegex(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceId, "robots.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceId, "robots.0.alert_robot_name", resourceId, "alert_robot_name"),
				),
			},
		},
	})
}

func testAccCmsAlertRobotsDataSourceConfig() string {
	return `
data "alicloud_cms_alert_robots" "default" {
  ids = [alicloud_cms_alert_robot.default.id]
}
`
}

func testAccCmsAlertRobotsDataSourceConfigWithNameRegex() string {
	return `
data "alicloud_cms_alert_robots" "default" {
  name_regex = "^${alicloud_cms_alert_robot.default.alert_robot_name}$"
}
`
}
