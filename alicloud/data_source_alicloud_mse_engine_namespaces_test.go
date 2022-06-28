package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMseEngineNamespacesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	checkoutSupportedRegions(t, true, connectivity.MSESupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMseEngineNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_engine_namespace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudMseEngineNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_mse_engine_namespace.default.id}_fake"]`,
		}),
	}
	var existAlicloudMseEngineNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"namespaces.#":                     "1",
			"namespaces.0.namespace_desc":      "",
			"namespaces.0.namespace_show_name": fmt.Sprintf("tf-testAccEngineNamespace-%d", rand),
			"namespaces.0.namespace_id":        fmt.Sprintf("tf-testAccEngineNamespace-%d", rand),
			"namespaces.0.service_count":       CHECKSET,
			"namespaces.0.quota":               CHECKSET,
			"namespaces.0.type":                CHECKSET,
			"namespaces.0.config_count":        CHECKSET,
			"namespaces.0.id":                  CHECKSET,
		}
	}
	var fakeAlicloudMseEngineNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#": "0",
		}
	}
	var alicloudMseEngineNamespacesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_mse_engine_namespaces.default",
		existMapFunc: existAlicloudMseEngineNamespacesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudMseEngineNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alicloudMseEngineNamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}
func testAccCheckAlicloudMseEngineNamespacesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccEngineNamespace-%d"
}

data "alicloud_mse_clusters" "default" {
	name_regex = "default-NODELETING"
}
resource "alicloud_mse_engine_namespace" "default" {
	cluster_id = data.alicloud_mse_clusters.default.clusters.0.cluster_id
	namespace_show_name = var.name
	namespace_id = var.name
}

data "alicloud_mse_engine_namespaces" "default" {	
	cluster_id = data.alicloud_mse_clusters.default.clusters.0.cluster_id
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
