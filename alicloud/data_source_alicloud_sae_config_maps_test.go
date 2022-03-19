package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSAEConfigmapsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 200)
	checkoutSupportedRegions(t, true, connectivity.SaeSupportRegions)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeConfigmapsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_config_map.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudSaeConfigmapsDataSourceName(rand, map[string]string{
			"ids": `["${alicloud_sae_config_map.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeConfigmapsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_config_map.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSaeConfigmapsDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_sae_config_map.default.name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSaeConfigmapsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_sae_config_map.default.id}"]`,
			"name_regex": `"${alicloud_sae_config_map.default.name}"`,
		}),
		fakeConfig: testAccCheckAlicloudSaeConfigmapsDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_sae_config_map.default.id}_fake"]`,
			"name_regex": `"${alicloud_sae_config_map.default.name}_fake"`,
		}),
	}
	var existAlicloudSaeNamespaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                "1",
			"names.#":              "1",
			"maps.#":               "1",
			"maps.0.name":          fmt.Sprintf("tf-testaccsaenames-%d", rand),
			"maps.0.description":   fmt.Sprintf("tf-testaccsaenames-%d", rand),
			"maps.0.namespace_id":  fmt.Sprintf("%s:configtest%d", defaultRegionToTest, rand),
			"maps.0.data":          "{\"env.home\":\"/root\",\"env.shell\":\"/bin/sh\"}",
			"maps.0.create_time":   CHECKSET,
			"maps.0.config_map_id": CHECKSET,
		}
	}
	var fakeAlicloudSaeNamespaceDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alicloudSaeNamespaceCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sae_config_maps.default",
		existMapFunc: existAlicloudSaeNamespaceDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlicloudSaeNamespaceDataSourceNameMapFunc,
	}
	alicloudSaeNamespaceCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlicloudSaeConfigmapsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testaccsaenames-%d"
}

variable "namespace_id" {	
	default = "%s:configtest%d"
}

resource "alicloud_sae_namespace" "default" {
	namespace_id = var.namespace_id
	namespace_name = var.name
	namespace_description = var.name
}

resource "alicloud_sae_config_map" "default" {
	namespace_id = alicloud_sae_namespace.default.namespace_id
	name = var.name
	description = var.name
	data = jsonencode({"env.home": "/root", "env.shell": "/bin/sh"})
}

data "alicloud_sae_config_maps" "default" {	
	namespace_id = "%s:configtest%d"
	%s
}
`, rand, defaultRegionToTest, rand, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
