package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSaeConfigmapsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 200)
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
			"maps.0.name":          fmt.Sprintf("tf-testaccsaenames-%d",rand),
			"maps.0.description":   fmt.Sprintf("tf-testaccsaenamespacedesc-%d",rand),
			"maps.0.namespace_id":  fmt.Sprintf("%s:configtest",os.Getenv("ALICLOUD_REGION")),
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
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.SaeSupportRegions)
	}
	alicloudSaeNamespaceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand,preCheck, idsConf, nameRegexConf, allConf)
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
variable "desc" {	
	default = "tf-testaccsaenamespacedesc-%d"
}
variable "namespace_id" {	
	default = "%s:configtest"
}

resource "alicloud_sae_config_map" "default" {
	namespace_id = var.namespace_id
	name = var.name
	description = var.desc
	data = jsonencode({"env.home": "/root", "env.shell": "/bin/sh"})
}

data "alicloud_sae_config_maps" "default" {	
	namespace_id = "%s:configtest"
	%s
}
`, rand, rand,os.Getenv("ALICLOUD_REGION"),os.Getenv("ALICLOUD_REGION"), strings.Join(pairs, " \n "))
	return config
}
