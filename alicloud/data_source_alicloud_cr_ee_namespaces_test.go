package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCREENamespacesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	namespaceName := fmt.Sprintf("tf-testacc-cr-ee-ns-%d", rand)
	resourceId := "data.alicloud_cr_ee_namespaces.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, namespaceName,
		dataSourceCrEENamespacesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"ids":         []string{"${alicloud_cr_ee_namespace.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"ids":         []string{"test-id-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}",
			"ids":         []string{"${alicloud_cr_ee_namespace.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}-fake",
			"ids":         []string{"test-id-fake"},
		}),
	}

	var existCrEENamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           CHECKSET,
			"names.#":                         "1",
			"names.0":                         namespaceName,
			"namespaces.#":                    "1",
			"namespaces.0.name":               namespaceName,
			"namespaces.0.default_visibility": "PRIVATE",
			"namespaces.0.auto_create":        "true",
			"namespaces.0.instance_id":        CHECKSET,
			"namespaces.0.id":                 CHECKSET,
			"namespaces.0.namespace_name":     namespaceName,
			"namespaces.0.namespace_id":       CHECKSET,
		}
	}

	var fakeCrEENamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"namespaces.#": "0",
		}
	}

	var crEENamespacesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCrEENamespacesMapFunc,
		fakeMapFunc:  fakeCrEENamespacesMapFunc,
	}
	crEENamespacesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceCrEENamespacesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	data "alicloud_cr_ee_instances" "default"{}
	
	resource "alicloud_cr_ee_namespace" "default" {
		instance_id = data.alicloud_cr_ee_instances.default.ids.0
		name = "${var.name}"
		auto_create	= true
		default_visibility = "PRIVATE"
	}
	`, name)
}
