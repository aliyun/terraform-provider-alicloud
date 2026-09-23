package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceMeshVersionsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshVersionsDataSourceName(rand, map[string]string{
			"ids": `["${data.alicloud_service_mesh_versions.pre.ids.0}"]`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshVersionsDataSourceName(rand, map[string]string{
			"ids": `["${data.alicloud_service_mesh_versions.pre.ids.0}_fakeid"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshVersionsDataSourceName(rand, map[string]string{
			"ids":     `["${data.alicloud_service_mesh_versions.pre.ids.0}"]`,
			"edition": `"Default"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshVersionsDataSourceName(rand, map[string]string{
			"ids":     `["${data.alicloud_service_mesh_versions.pre.ids.0}_fake"]`,
			"edition": `"Default"`,
		}),
	}

	var existDataAlicloudServiceMeshVersionsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":              "1",
			"versions.#":         "1",
			"versions.0.id":      CHECKSET,
			"versions.0.version": CHECKSET,
			"versions.0.edition": "Default",
		}
	}
	var fakeDataAlicloudServiceMeshVersionsSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"versions.#": "0",
		}
	}
	var alicloudServiceMeshServiceMeshCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_service_mesh_versions.default",
		existMapFunc: existDataAlicloudServiceMeshVersionsSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudServiceMeshVersionsSourceNameMapFunc,
	}

	alicloudServiceMeshServiceMeshCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}
func testAccCheckAlicloudServiceMeshVersionsDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_service_mesh_versions" "pre" {
	edition = "Default"
}

data "alicloud_service_mesh_versions" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
	return config
}
