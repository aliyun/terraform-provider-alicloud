package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCREENamespacesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	namespaceName := fmt.Sprintf("tf-testacc-cr-ee-ns-%d", rand)
	resourceId := "data.alicloud_cr_ee_namespaces.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, namespaceName,
		dataSourceCrEENamespacesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}-fake",
		}),
	}

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"ids":         []string{"${alicloud_cr_ee_namespace.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"ids":         []string{"test-id-fake"},
		}),
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}",
			"ids":         []string{"${alicloud_cr_ee_namespace.default.id}"},
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${alicloud_cr_ee_namespace.default.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}-fake",
			"ids":         []string{"test-id-fake"},
		}),
	}

	var existCrEENamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           "1",
			"ids.0":                           acctest.CHECKSET,
			"names.#":                         "1",
			"names.0":                         namespaceName,
			"namespaces.#":                    "1",
			"namespaces.0.name":               namespaceName,
			"namespaces.0.default_visibility": "PRIVATE",
			"namespaces.0.auto_create":        "true",
			"namespaces.0.instance_id":        acctest.CHECKSET,
			"namespaces.0.id":                 acctest.CHECKSET,
			"namespaces.0.namespace_name":     namespaceName,
			"namespaces.0.namespace_id":       acctest.CHECKSET,
		}
	}

	var fakeCrEENamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"namespaces.#": "0",
		}
	}

	var crEENamespacesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCrEENamespacesMapFunc,
		FakeMapFunc:  fakeCrEENamespacesMapFunc,
	}
	crEENamespacesCheckInfo.DataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceCrEENamespacesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	# Provision a dedicated EE instance instead of relying on
	# data.alicloud_cr_ee_instances {} — that data source returns an empty list
	# when no instance exists at test time (frequent in cn-beijing under
	# concurrent ACC runs). Creating one inline makes this datasource test
	# self-contained.
	resource "alicloud_cr_ee_instance" "default" {
		payment_type   = "Subscription"
		period         = 1
		renew_period   = 1
		renewal_status = "AutoRenewal"
		instance_type  = "Advanced"
		instance_name  = var.name
	}

	resource "alicloud_cr_ee_namespace" "default" {
		instance_id        = alicloud_cr_ee_instance.default.id
		name               = "${var.name}"
		auto_create        = true
		default_visibility = "PRIVATE"
	}
	`, name)
}
