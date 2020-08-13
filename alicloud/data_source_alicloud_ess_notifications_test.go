package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEssNotificationsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssNotificationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_notification.default.scaling_group_id}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEssNotificationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_notification.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_notification.default.notification_arn}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEssNotificationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${alicloud_ess_notification.default.scaling_group_id}"`,
			"ids":              `["${alicloud_ess_notification.default.notification_arn}_fake"]`,
		}),
	}

	var existEssnotificationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"notifications.#":                  "1",
			"notifications.0.notification_arn": CHECKSET,
			"notifications.0.scaling_group_id": CHECKSET,
		}
	}

	var fakeEssnotificationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"notifications.#": "0",
			"ids.#":           "0",
		}
	}

	var essNotificationsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_ess_notifications.default",
		existMapFunc: existEssnotificationsMapFunc,
		fakeMapFunc:  fakeEssnotificationsMapFunc,
	}

	essNotificationsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, allConf)
}

func testAccCheckAlicloudEssNotificationsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssNotifications-%d"
}
data "alicloud_regions" "default" {
    current = true
}

data "alicloud_account" "default" {
}

resource "alicloud_ess_scaling_group" "default" {
    min_size = 1
    max_size = 1
    scaling_group_name = "${var.name}"
    removal_policies = ["OldestInstance", "NewestInstance"]
    vswitch_ids = ["${alicloud_vswitch.default.id}"]
}

resource "alicloud_mns_queue" "default"{
    name="${var.name}"
}

resource "alicloud_ess_notification" "default" {
    scaling_group_id = "${alicloud_ess_scaling_group.default.id}"
    notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS"]
    notification_arn = "acs:ess:${data.alicloud_regions.default.regions.0.id}:${data.alicloud_account.default.id}:queue/${alicloud_mns_queue.default.name}"
}

data "alicloud_ess_notifications" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
