package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudSaeApplicationGroupsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_sae_application_groups.default"
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSaeApplicationGroupsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"app_id": "${data.alicloud_sae_application_groups.image.app_id}",
			"ids":    []string{"${data.alicloud_sae_application_groups.image.groups.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"app_id": "${data.alicloud_sae_application_groups.image.app_id}",
			"ids":    []string{"${data.alicloud_sae_application_groups.image.groups.0.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.image.app_id}",
			"name_regex": "${data.alicloud_sae_application_groups.image.groups.0.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.image.app_id}",
			"name_regex": "${data.alicloud_sae_application_groups.image.groups.0.group_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.image.app_id}",
			"ids":        []string{"${data.alicloud_sae_application_groups.image.groups.0.id}"},
			"name_regex": "${data.alicloud_sae_application_groups.image.groups.0.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.image.app_id}",
			"ids":        []string{"${data.alicloud_sae_application_groups.image.groups.0.id}_fake"},
			"name_regex": "${data.alicloud_sae_application_groups.image.groups.0.group_name}_fake",
		}),
	}

	var existAliCloudSaeApplicationGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"groups.#":                        "1",
			"groups.0.id":                     CHECKSET,
			"groups.0.group_id":               CHECKSET,
			"groups.0.group_name":             CHECKSET,
			"groups.0.group_type":             CHECKSET,
			"groups.0.replicas":               CHECKSET,
			"groups.0.package_type":           CHECKSET,
			"groups.0.package_version":        CHECKSET,
			"groups.0.package_version_id":     CHECKSET,
			"groups.0.package_url":            "",
			"groups.0.image_url":              CHECKSET,
			"groups.0.web_container":          "",
			"groups.0.jdk":                    "",
			"groups.0.running_instances":      CHECKSET,
			"groups.0.edas_container_version": CHECKSET,
		}
	}

	var fakeAliCloudSaeApplicationGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var aliCloudSaeApplicationGroupsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_application_groups.default",
		existMapFunc: existAliCloudSaeApplicationGroupsMapFunc,
		fakeMapFunc:  fakeAliCloudSaeApplicationGroupsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudSaeApplicationGroupsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func TestAccAliCloudSaeApplicationGroupsDataSource_basic1(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_sae_application_groups.default"
	name := fmt.Sprintf("tfacc%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceSaeApplicationGroupsConfig)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"app_id": "${data.alicloud_sae_application_groups.war.app_id}",
			"ids":    []string{"${data.alicloud_sae_application_groups.war.groups.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"app_id": "${data.alicloud_sae_application_groups.war.app_id}",
			"ids":    []string{"${data.alicloud_sae_application_groups.war.groups.0.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.war.app_id}",
			"name_regex": "${data.alicloud_sae_application_groups.war.groups.0.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.war.app_id}",
			"name_regex": "${data.alicloud_sae_application_groups.war.groups.0.group_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.war.app_id}",
			"ids":        []string{"${data.alicloud_sae_application_groups.war.groups.0.id}"},
			"name_regex": "${data.alicloud_sae_application_groups.war.groups.0.group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"app_id":     "${data.alicloud_sae_application_groups.war.app_id}",
			"ids":        []string{"${data.alicloud_sae_application_groups.war.groups.0.id}_fake"},
			"name_regex": "${data.alicloud_sae_application_groups.war.groups.0.group_name}_fake",
		}),
	}

	var existAliCloudSaeApplicationGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"names.#":                         "1",
			"groups.#":                        "1",
			"groups.0.id":                     CHECKSET,
			"groups.0.group_id":               CHECKSET,
			"groups.0.group_name":             CHECKSET,
			"groups.0.group_type":             CHECKSET,
			"groups.0.replicas":               CHECKSET,
			"groups.0.package_type":           CHECKSET,
			"groups.0.package_version":        CHECKSET,
			"groups.0.package_version_id":     CHECKSET,
			"groups.0.package_url":            CHECKSET,
			"groups.0.image_url":              "",
			"groups.0.web_container":          CHECKSET,
			"groups.0.jdk":                    CHECKSET,
			"groups.0.running_instances":      CHECKSET,
			"groups.0.edas_container_version": "",
		}
	}

	var fakeAliCloudSaeApplicationGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var aliCloudSaeApplicationGroupsInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_application_groups.default",
		existMapFunc: existAliCloudSaeApplicationGroupsMapFunc,
		fakeMapFunc:  fakeAliCloudSaeApplicationGroupsMapFunc,
	}

	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.TestSalveRegions)
	}

	aliCloudSaeApplicationGroupsInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}

func dataSourceSaeApplicationGroupsConfig(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_sae_applications" "image" {
  app_name = "default-NODELETING-image"
}

data "alicloud_sae_applications" "war" {
  app_name = "default-NODELETING-war"
}

data "alicloud_sae_application_groups" "image" {
  app_id = data.alicloud_sae_applications.image.applications.0.id
}

data "alicloud_sae_application_groups" "war" {
  app_id = data.alicloud_sae_applications.war.applications.0.id
}
`, name)
}
