package alicloud

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudCRNamespacesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cr_namespaces.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-cr-ns-%d", rand),
		dataSourceCRNamespacesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cr_namespace.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cr_namespace.default.name}-fake",
		}),
	}

	var existCRNamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                         "1",
			"names.0":                         fmt.Sprintf("tf-testacc-cr-ns-%d", rand),
			"namespaces.#":                    "1",
			"namespaces.0.name":               fmt.Sprintf("tf-testacc-cr-ns-%d", rand),
			"namespaces.0.default_visibility": "PUBLIC",
			"namespaces.0.auto_create":        "false",
		}
	}

	var fakeCRNamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":      "0",
			"namespaces.#": "0",
		}
	}

	var crNamespacesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCRNamespacesMapFunc,
		fakeMapFunc:  fakeCRNamespacesMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
	}
	crNamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
}

func dataSourceCRNamespacesConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	resource "alicloud_cr_namespace" "default" {
		name = "${var.name}"
		auto_create	= false
		default_visibility = "PUBLIC"
	}
	`, name)
}
