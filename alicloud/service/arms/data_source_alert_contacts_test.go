package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudARMSAlertContactsDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_arms_alert_contacts.default"
	name := fmt.Sprintf("tf-testacc-armsAlertContacts%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsAlertContactsConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_alert_contact.default.alert_contact_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"alert_contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"alert_contact_name": "fake_tf-testacc*",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact.default.id}_fake"},
		}),
	}

	emailConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"email": "${alicloud_arms_alert_contact.default.email}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"email": "${alicloud_arms_alert_contact.default.email}_fake",
		}),
	}
	phoneNumConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"phone_num": "${alicloud_arms_alert_contact.default.phone_num}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"phone_num": "${alicloud_arms_alert_contact.default.phone_num}00",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex":         "${alicloud_arms_alert_contact.default.alert_contact_name}",
			"alert_contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}",
			"ids":                []string{"${alicloud_arms_alert_contact.default.id}"},
			"email":              "${alicloud_arms_alert_contact.default.email}",
			"phone_num":          "${alicloud_arms_alert_contact.default.phone_num}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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
			"contacts.0.id":                     acctest.CHECKSET,
			"contacts.0.alert_contact_id":       acctest.CHECKSET,
			"contacts.0.alert_contact_name":     name,
			"contacts.0.create_time":            acctest.CHECKSET,
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

	var ArmsAlertContactsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existArmsAlertContactsMapFunc,
		FakeMapFunc:  fakeArmsAlertContactsMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsAlertContactsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, nameConf, idsConf, emailConf, phoneNumConf, allConf)
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
