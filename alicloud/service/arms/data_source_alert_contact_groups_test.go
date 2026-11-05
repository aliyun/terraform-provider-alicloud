package arms_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudARMSAlertContactGroupsDataSource(t *testing.T) {

	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_arms_alert_contact_groups.default"
	name := fmt.Sprintf("tf-testacc-ArmsAlertContactGroups%v", rand)

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsAlertContactGroupsConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_tf-testacc*",
		}),
	}

	nameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"alert_contact_group_name": "${alicloud_arms_alert_contact_group.default.alert_contact_group_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"alert_contact_group_name": "fake_tf-testacc*",
		}),
	}
	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact_group.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_alert_contact_group.default.id}_fake"},
		}),
	}

	contactIdConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alicloud_arms_alert_contact_group.default.id}"},
			"contact_id": "${alicloud_arms_alert_contact.default.id}",
		}),
	}
	contactNameConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_arms_alert_contact_group.default.id}"},
			"contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alicloud_arms_alert_contact_group.default.id}"},
			"contact_name": "${alicloud_arms_alert_contact.default.alert_contact_name}-fake",
		}),
	}
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
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
			"groups.0.id":                       acctest.CHECKSET,
			"groups.0.alert_contact_group_id":   acctest.CHECKSET,
			"groups.0.alert_contact_group_name": name,
			"groups.0.contact_ids.#":            "1",
			"groups.0.create_time":              acctest.CHECKSET,
		}
	}

	var fakeArmsAlertContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"groups.#": "0",
			"names.#":  "0",
			"ids.#":    "0",
		}
	}

	var ArmsAlertContactGroupsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existArmsAlertContactGroupsMapFunc,
		FakeMapFunc:  fakeArmsAlertContactGroupsMapFunc,
	}
	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, connectivity.ARMSSupportRegions)
	}
	ArmsAlertContactGroupsCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, nameConf, idsConf, contactIdConf, contactNameConf, allConf)
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
