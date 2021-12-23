package alicloud

import (
	"testing"
)

func TestAccAlicloudCrEEInstancesDataSource(t *testing.T) {
	resourceId := "data.alicloud_cr_ee_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "",
		dataSourceCrEEInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "^tf-testacc",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "test-fake.*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"test-id-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":     "^tf-testacc",
			"enable_details": "true",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"test-id-fake"},
			"name_regex": "test-fake.*",
		}),
	}

	var existCrEEInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                         CHECKSET,
			"names.0":                         CHECKSET,
			"instances.#":                     CHECKSET,
			"instances.0.id":                  CHECKSET,
			"instances.0.name":                CHECKSET,
			"instances.0.region":              CHECKSET,
			"instances.0.specification":       CHECKSET,
			"instances.0.namespace_quota":     CHECKSET,
			"instances.0.namespace_usage":     CHECKSET,
			"instances.0.repo_quota":          CHECKSET,
			"instances.0.repo_usage":          CHECKSET,
			"instances.0.vpc_endpoints.#":     CHECKSET,
			"instances.0.public_endpoints.#":  CHECKSET,
			"instances.0.authorization_token": CHECKSET,
			"instances.0.temp_username":       CHECKSET,
		}
	}

	var fakeCrEEInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var crEEInstancesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCrEEInstancesMapFunc,
		fakeMapFunc:  fakeCrEEInstancesMapFunc,
	}

	crEEInstancesCheckInfo.dataSourceTestCheck(t, 0, nameRegexConf, idsConf, allConf)
}

func dataSourceCrEEInstancesConfigDependence(name string) string {
	return ""
}
