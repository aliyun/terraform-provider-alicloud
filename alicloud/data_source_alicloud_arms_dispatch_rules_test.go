package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudARMSDispatchRulesDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_arms_dispatch_rules.default"
	name := fmt.Sprintf("tf-testacc-armsDispatchRules%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsDispatchRulesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"dispatch_rule_name": "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"enable_details":     "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"dispatch_rule_name": "fake_tf-testacc*",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alicloud_arms_dispatch_rule.default.id}"},
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_dispatch_rule.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":         "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"dispatch_rule_name": "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"ids":                []string{"${alicloud_arms_dispatch_rule.default.id}"},
			"enable_details":     "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":         "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}",
			"dispatch_rule_name": "${alicloud_arms_dispatch_rule.default.dispatch_rule_name}_fake",
			"ids":                []string{"${alicloud_arms_dispatch_rule.default.id}"},
		}),
	}
	var existArmsDispatchRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"rules.#":                               "1",
			"rules.0.id":                            CHECKSET,
			"rules.0.dispatch_rule_id":              CHECKSET,
			"rules.0.group_rules.#":                 "1",
			"rules.0.status":                        "enable",
			"rules.0.label_match_expression_grid.#": "1",
			"rules.0.notify_rules.#":                "1",
			"rules.0.dispatch_rule_name":            name,
		}
	}

	var fakeArmsDispatchRulesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"rules.#": "0",
			"names.#": "0",
			"ids.#":   "0",
		}
	}

	var ArmsDispatchRulesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existArmsDispatchRulesMapFunc,
		fakeMapFunc:  fakeArmsDispatchRulesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsDispatchRulesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, nameConf, idsConf, allConf)
}

func dataSourceArmsDispatchRulesConfigDependence(name string) string {
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

resource "alicloud_arms_dispatch_rule" "default" {
  dispatch_rule_name = var.name
  dispatch_type      = "CREATE_ALERT"
  group_rules {
    group_wait_time = 5
    group_interval  = 15
    repeat_interval = 100
    grouping_fields = [
      "alertname"]
  }
  label_match_expression_grid {
    label_match_expression_groups {
      label_match_expressions {
        key      = "_aliyun_arms_involvedObject_kind"
        value    = "app"
        operator = "eq"
      }
    }
  }

  notify_rules {
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact.default.id
      notify_type      = "ARMS_CONTACT"
      name             = var.name
    }
    notify_objects {
      notify_object_id = alicloud_arms_alert_contact_group.default.id
      notify_type      = "ARMS_CONTACT_GROUP"
      name             = var.name
    }
    notify_channels = ["dingTalk", "wechat"]
  }
}
`, name)
}
