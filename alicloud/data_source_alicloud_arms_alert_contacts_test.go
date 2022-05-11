package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudARMSAlertContactsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_arms_alert_contacts.default"
	name := fmt.Sprintf("tf-testacc-armsAlertContacts%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsAlertContactsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_alert_contact.default.alert_contact_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"alert_contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"alert_contact_name": "fake_tf-testacc*",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact.default.id}_fake"},
		}),
	}

	emailConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"email": "${alicloud_arms_alert_contact.default.email}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"email": "${alicloud_arms_alert_contact.default.email}_fake",
		}),
	}
	phoneNumConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"phone_num": "${alicloud_arms_alert_contact.default.phone_num}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"phone_num": "${alicloud_arms_alert_contact.default.phone_num}00",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":         "${alicloud_arms_alert_contact.default.alert_contact_name}",
			"alert_contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}",
			"ids":                []string{"${alicloud_arms_alert_contact.default.id}"},
			"email":              "${alicloud_arms_alert_contact.default.email}",
			"phone_num":          "${alicloud_arms_alert_contact.default.phone_num}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":         "${alicloud_arms_alert_contact.default.alert_contact_name}",
			"alert_contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}_fake",
			"ids":                []string{"${alicloud_arms_alert_contact.default.id}"},
			"email":              "${alicloud_arms_alert_contact.default.email}",
			"phone_num":          "${alicloud_arms_alert_contact.default.phone_num}00",
		}),
	}
	var existArmsAlertContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"contacts.#":                        "1",
			"contacts.0.id":                     CHECKSET,
			"contacts.0.alert_contact_id":       CHECKSET,
			"contacts.0.alert_contact_name":     name,
			"contacts.0.create_time":            CHECKSET,
			"contacts.0.ding_robot_webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=91f2f7",
			"contacts.0.email":                  "hello.uuuu@aaa.com",
			"contacts.0.phone_num":              "12345678901",
			"contacts.0.system_noc":             "false",
			"contacts.0.webhook":                "",
		}
	}

	var fakeArmsAlertContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"contacts.#": "0",
			"names.#":    "0",
			"ids.#":      "0",
		}
	}

	var ArmsAlertContactsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existArmsAlertContactsMapFunc,
		fakeMapFunc:  fakeArmsAlertContactsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsAlertContactsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, nameConf, idsConf, emailConf, phoneNumConf, allConf)
}

func dataSourceArmsAlertContactsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		resource "alicloud_arms_alert_contact" "default" {
		  alert_contact_name = var.name
		  ding_robot_webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=91f2f7"
          email = "hello.uuuu@aaa.com"
          phone_num = "12345678901"
          system_noc = "false"
		}
		`, name)
}
