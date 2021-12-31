package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCrEENamespacesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	namespaceName := fmt.Sprintf("tf-testacc-cr-ee-ns-%d", rand)
	resourceId := "data.alicloud_cr_ee_namespaces.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, namespaceName,
		dataSourceCrEENamespacesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}-fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"ids":         []string{"test-id-fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
			"name_regex":  "${alicloud_cr_ee_namespace.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${local.instance_id}",
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

	resource "alicloud_cr_ee_instance" "default" {
	  count = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? 0 : 1
	  period = 1
	  renew_period = 0
	  payment_type = "Subscription"
	  instance_type = "Basic"
	  renewal_status = "ManualRenewal"
	  instance_name = "tf-testacc-basic"
	}
	
	data "alicloud_cr_ee_instances" "default"{
	  name_regex = "tf-testacc"
	}

	locals {
	  instance_id=length(data.alicloud_cr_ee_instances.default.ids)>0? data.alicloud_cr_ee_instances.default.ids[0] : concat(alicloud_cr_ee_instance.default.*.id, [""])[0]
	}
	
	resource "alicloud_cr_ee_namespace" "default" {
		instance_id = local.instance_id
		name = "${var.name}"
		auto_create	= true
		default_visibility = "PRIVATE"
	}
	`, name)
}
