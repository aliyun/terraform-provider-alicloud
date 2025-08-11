package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudArmsEnvironmentsDataSource_basic0(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	resourceId := "data.alicloud_arms_environments.default"
	name := fmt.Sprintf("tf-testacc%sarmsenvironment%d", defaultRegionToTest, rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceArmsEnvironmentsConfig0)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_environment.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_arms_environment.default.id}_fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_environment.default.environment_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_arms_environment.default.environment_name}_fake",
		}),
	}
	environmentTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${alicloud_arms_environment.default.id}"},
			"environment_type": "${alicloud_arms_environment.default.environment_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"environment_type": "Cloud",
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${alicloud_arms_environment.default.resource_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.2}",
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Environment",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Created": "TF_Fake",
				"For":     "Environment_Fake",
			},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_arms_environment.default.id}"},
			"name_regex":        "${alicloud_arms_environment.default.environment_name}",
			"environment_type":  "${alicloud_arms_environment.default.environment_type}",
			"resource_group_id": "${alicloud_arms_environment.default.resource_group_id}",
			"tags": map[string]string{
				"Created": "TF",
				"For":     "Environment",
			},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":               []string{"${alicloud_arms_environment.default.id}_fake"},
			"name_regex":        "${alicloud_arms_environment.default.environment_name}_fake",
			"environment_type":  "Cloud",
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.2}",
			"tags": map[string]string{
				"Created": "TF_Fake",
				"For":     "Environment_Fake",
			},
		}),
	}
	var existAliCloudArmsEnvironmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 "1",
			"names.#":                               "1",
			"environments.#":                        "1",
			"environments.0.id":                     CHECKSET,
			"environments.0.bind_resource_id":       CHECKSET,
			"environments.0.bind_resource_type":     CHECKSET,
			"environments.0.bind_vpc_cidr":          CHECKSET,
			"environments.0.environment_id":         CHECKSET,
			"environments.0.environment_name":       CHECKSET,
			"environments.0.environment_type":       CHECKSET,
			"environments.0.grafana_datasource_uid": CHECKSET,
			"environments.0.grafana_folder_uid":     CHECKSET,
			"environments.0.managed_type":           CHECKSET,
			"environments.0.prometheus_instance_id": CHECKSET,
			"environments.0.region_id":              CHECKSET,
			"environments.0.resource_group_id":      CHECKSET,
			"environments.0.tags.%":                 "2",
			"environments.0.tags.Created":           "TF",
			"environments.0.tags.For":               "Environment",
			"environments.0.user_id":                CHECKSET,
		}
	}
	var fakeAliCloudArmsEnvironmentsDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"names.#":        "0",
			"environments.#": "0",
		}
	}
	var alicloudArmsEnvironmentsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_arms_environments.default",
		existMapFunc: existAliCloudArmsEnvironmentsDataSourceNameMapFunc,
		fakeMapFunc:  fakeAliCloudArmsEnvironmentsDataSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudArmsEnvironmentsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, environmentTypeConf, resourceGroupIdConf, tagsConf, allConf)
}

func dataSourceArmsEnvironmentsConfig0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_resource_manager_resource_groups" "default"{
		status = "OK"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	resource "alicloud_arms_environment" "default" {
  		bind_resource_id     = data.alicloud_vpcs.default.ids.0
  		environment_sub_type = "ECS"
  		environment_type     = "ECS"
  		environment_name     = var.name
  		resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.1
  		tags = {
    		Created = "TF"
    		For     = "Environment"
  		}
	}
`, name)
}
