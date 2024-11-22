package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudARMSAlertRobotsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_arms_alert_robots.default"
	name := fmt.Sprintf("tf-testacc-ArmsAlertRobots%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsAlertRobotsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_alert_robot.default.alert_robot_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_robot.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_robot.default.id}_fake"},
		}),
	}
	robotNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"alert_robot_name": "${alicloud_arms_alert_robot.default.alert_robot_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"alert_robot_name": "${alicloud_arms_alert_robot.default.alert_robot_name}-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"robot_type":       "${alicloud_arms_alert_robot.default.robot_type}",
			"name_regex":       "${alicloud_arms_alert_robot.default.alert_robot_name}",
			"ids":              []string{"${alicloud_arms_alert_robot.default.id}"},
			"alert_robot_name": "${alicloud_arms_alert_robot.default.alert_robot_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"robot_type":       "${alicloud_arms_alert_robot.default.robot_type}",
			"name_regex":       "${alicloud_arms_alert_robot.default.alert_robot_name}",
			"ids":              []string{"${alicloud_arms_alert_robot.default.id}"},
			"alert_robot_name": "${alicloud_arms_alert_robot.default.alert_robot_name}-fake",
		}),
	}
	existArmsAlertRobotsMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"names.#":                 "1",
			"robots.#":                "1",
			"robots.0.id":             CHECKSET,
			"robots.0.robot_name":     name,
			"robots.0.robot_type":     "wechat",
			"robots.0.robot_id":       CHECKSET,
			"robots.0.robot_addr":     "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=1c704e23",
			"robots.0.daily_noc":      "true",
			"robots.0.daily_noc_time": "09:30,17:00",
			"robots.0.create_time":    CHECKSET,
		}
	}

	fakeArmsAlertRobotsMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"robots.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	ArmsAlertRobotsCheckInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existArmsAlertRobotsMapFunc,
		fakeMapFunc:  fakeArmsAlertRobotsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsAlertRobotsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, robotNameConf, allConf)
}

func dataSourceArmsAlertRobotsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}
		resource "alicloud_arms_alert_robot" "default" {
		  alert_robot_name = var.name
		  robot_type       = "wechat"
		  robot_addr       = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=1c704e23"
		  daily_noc        = "true"
		  daily_noc_time   = "09:30,17:00"
		}
		`, name)
}
