package cr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCRNamespacesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_cr_namespaces.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-cr-ns-%d", rand),
		dataSourceCRNamespacesConfigDependence)

	nameRegexConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_cr_namespace.default.name}",
		}),
		FakeConfig: testAccConfig(map[string]interface{}{
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

	var crNamespacesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existCRNamespacesMapFunc,
		FakeMapFunc:  fakeCRNamespacesMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, false, connectivity.CRNoSupportedRegions)
	}
	crNamespacesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf)
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
