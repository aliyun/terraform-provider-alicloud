package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudARMSAlertContactGroupsDataSource(t *testing.T) {

	rand := acctest.RandInt()
	resourceId := "data.alicloud_arms_alert_contact_groups.default"
	name := fmt.Sprintf("tf-testacc-ArmsAlertContactGroups%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsAlertContactGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"alert_contact_group_name": "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"alert_contact_group_name": "fake_tf-testacc*",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact_group.default.id}_fake"},
		}),
	}

	contactIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_arms_alert_contact_group.default.id}"},
			"contact_id": "${alicloud_arms_alert_contact.default.id}",
		}),
	}
	contactNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_arms_alert_contact_group.default.id}"},
			"contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_arms_alert_contact_group.default.id}"},
			"contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":               "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}",
			"alert_contact_group_name": "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}",
			"ids":                      []string{"${alicloud_arms_alert_contact_group.default.id}"},
			"contact_id":               "${alicloud_arms_alert_contact.default.id}",
			"contact_name":             "${alicloud_arms_alert_contact.default.alert_contact_name}",
		}),
		// There is an API error when fetching one resource with multi conditions
		//fakeConfig: testAccConfig(map[string]interface{}{
		//	"name_regex":               "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}",
		//	"alert_contact_group_name": "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}-fake",
		//	"ids":                      []string{"${alicloud_arms_alert_contact_group.default.id}"},
		//	"contact_id":               "${alicloud_arms_alert_contact.default.id}",
		//	"contact_name":             "${alicloud_arms_alert_contact.default.alert_contact_name}-fake",
		//}),
	}
	var existArmsAlertContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"groups.#":                          "1",
			"groups.0.id":                       CHECKSET,
			"groups.0.alert_contact_group_id":   CHECKSET,
			"groups.0.alert_contact_group_name": name,
			"groups.0.contact_ids.#":            "1",
			"groups.0.create_time":              CHECKSET,
		}
	}

	var fakeArmsAlertContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var ArmsAlertContactGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existArmsAlertContactGroupsMapFunc,
		fakeMapFunc:  fakeArmsAlertContactGroupsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsAlertContactGroupsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, nameConf, idsConf, contactIdConf, contactNameConf, allConf)
}

func dataSourceArmsAlertContactGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
		variable "name" {
		 default = "%v"
		}

		resource "alicloud_arms_alert_contact" "default" {
		  alert_contact_name = var.name
          email = "${var.name}@aaa.com"
		}
		resource "alicloud_arms_alert_contact_group" "default" {
		  alert_contact_group_name = var.name
          contact_ids = [alicloud_arms_alert_contact.default.id]
		}
		`, name)
}
