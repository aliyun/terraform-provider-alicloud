package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudECSKeyPairsDataSource(t *testing.T) {
	resourceId := "data.alicloud_ecs_key_pairs.default"
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsKeyPairsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsKeyPairsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ecs_key_pair.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"ids":        []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"ids":        []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Create": "TF",
				"For":    "Ecs Key Pairs",
			},
			"ids": []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"tags": map[string]string{
				"Create": "TF-fake",
				"For":    "Ecs Key Pairs Fake",
			},
			"ids": []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
			"ids":               []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
			"ids":               []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": name,
			"tags": map[string]string{
				"Create": "TF",
				"For":    "Ecs Key Pairs",
			},
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.0}",
			"ids":               []string{"${alicloud_ecs_key_pair.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": name + "fake",
			"tags": map[string]string{
				"Create": "TF-fake",
				"For":    "Ecs Key Pairs Fake",
			},
			"resource_group_id": "${data.alicloud_resource_manager_resource_groups.default.ids.1}",
			"ids":               []string{"${alicloud_ecs_key_pair.default.id}-fake"},
		}),
	}
	var existEcsKeyPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         CHECKSET,
			"names.#":                       "1",
			"names.0":                       name,
			"pairs.#":                       "1",
			"pairs.0.id":                    CHECKSET,
			"pairs.0.finger_print":          CHECKSET,
			"pairs.0.key_name":              name,
			"pairs.0.resource_group_id":     CHECKSET,
			"pairs.0.tags.%":                "2",
			"pairs.0.tags.Create":           "TF",
			"pairs.0.tags.For":              "Ecs Key Pairs",
			"key_pairs.#":                   "1",
			"key_pairs.0.id":                CHECKSET,
			"key_pairs.0.finger_print":      CHECKSET,
			"key_pairs.0.key_name":          name,
			"key_pairs.0.resource_group_id": CHECKSET,
			"key_pairs.0.tags.%":            "2",
			"key_pairs.0.tags.Create":       "TF",
			"key_pairs.0.tags.For":          "Ecs Key Pairs",
		}
	}

	var fakeEcsKeyPairsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"pairs.#":     "0",
			"key_pairs.#": "0",
		}
	}

	var EcsKeyPairsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsKeyPairsMapFunc,
		fakeMapFunc:  fakeEcsKeyPairsMapFunc,
	}

	EcsKeyPairsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, tagsConf, resourceGroupIdConf, allConf)
}

func dataSourceEcsKeyPairsDependence(name string) string {
	return fmt.Sprintf(`
	data "alicloud_resource_manager_resource_groups" "default"{
		status = "OK"
	}
	resource "alicloud_ecs_key_pair" "default" {
		key_name              = "%s"
		resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
		tags = {
			Create = "TF"
    		For = "Ecs Key Pairs",
  			}
	}`, name)
}
