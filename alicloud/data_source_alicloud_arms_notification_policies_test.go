package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudARMSNotificationPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_arms_notification_policies.default"
	name := fmt.Sprintf("tf-testacc-armsNotificationPolicies%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsNotificationPoliciesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_arms_notification_policy.default.name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name":           "${alicloud_arms_notification_policy.default.name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name": "fake_tf-testacc*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_arms_notification_policy.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_notification_policy.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_arms_notification_policy.default.name}",
			"name":           "${alicloud_arms_notification_policy.default.name}",
			"ids":            []string{"${alicloud_arms_notification_policy.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_notification_policy.default.name}",
			"name":       "${alicloud_arms_notification_policy.default.name}_fake",
			"ids":        []string{"${alicloud_arms_notification_policy.default.id}"},
		}),
	}

	var existArmsNotificationPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    "1",
			"names.#":                  "1",
			"policies.#":              "1",
			"policies.0.id":           CHECKSET,
			"policies.0.name":         name,
			"policies.0.notify_rule.#": "1",
			"policies.0.group_rule.#":  "1",
		}
	}

	var fakeArmsNotificationPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"policies.#": "0",
			"names.#":    "0",
			"ids.#":      "0",
		}
	}

	var ArmsNotificationPoliciesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existArmsNotificationPoliciesMapFunc,
		fakeMapFunc:  fakeArmsNotificationPoliciesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsNotificationPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, nameConf, idsConf, allConf)
}

func dataSourceArmsNotificationPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%v"
}

resource "alicloud_arms_alert_contact" "default" {
  alert_contact_name = var.name
  email              = "${var.name}@aaa.com"
}

resource "alicloud_arms_alert_contact_group" "default" {
  alert_contact_group_name = var.name
  contact_ids              = [alicloud_arms_alert_contact.default.id]
}

resource "alicloud_arms_notification_policy" "default" {
  name = var.name

  group_rule {
    group_wait      = 5
    group_interval  = 30
    grouping_fields = ["alertname"]
  }

  matching_rules {
    matching_conditions {
      key      = "_aliyun_arms_involvedObject_kind"
      value    = "app"
      operator = "eq"
    }
  }

  notify_rule {
    notify_start_time = "00:00"
    notify_end_time   = "23:59"
    notify_channels   = ["dingTalk", "email"]
    notify_objects {
      notify_object_type = "ARMS_CONTACT"
      notify_object_id   = alicloud_arms_alert_contact.default.id
      notify_object_name = var.name
    }
  }

  repeat          = true
  repeat_interval = 600
}
`, name)
}
