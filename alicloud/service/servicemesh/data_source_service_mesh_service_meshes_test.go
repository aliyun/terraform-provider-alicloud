package servicemesh_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudServiceMeshServiceMeshesDataSource(t *testing.T) {
	rand := tfacctest.RandIntRange(1, 1000)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_service_mesh_service_mesh.default.id}"]`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_service_mesh_service_mesh.default.id}_fakeid"]`,
			"enable_details": "true",
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"name_regex":     `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}_fake"`,
			"enable_details": "true",
		}),
	}

	statusRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_service_mesh_service_mesh.default.id}"]`,
			"status":         `"running"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_service_mesh_service_mesh.default.id}"]`,
			"status":         `"initial"`,
			"enable_details": "true",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_service_mesh_service_mesh.default.id}"]`,
			"name_regex":     `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}"`,
			"status":         `"running"`,
			"enable_details": "true",
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":            `["${alicloud_service_mesh_service_mesh.default.id}_fake"]`,
			"name_regex":     `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}_fake"`,
			"status":         `"initial"`,
			"enable_details": "true",
		}),
	}

	var existDataAlicloudServiceMeshServiceMeshesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"meshes.#":                   "1",
			"meshes.0.status":            "running",
			"meshes.0.sidecar_version":   acctest.CHECKSET,
			"meshes.0.kube_config":       acctest.CHECKSET,
			"meshes.0.service_mesh_name": fmt.Sprintf("tf-testaccservicemeshservicemesh-%d", rand),
		}
	}
	var fakeDataAlicloudServiceMeshServiceMeshesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"meshes.#": "0",
		}
	}
	var alicloudServiceMeshServiceMeshCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_service_mesh_service_meshes.default",
		existMapFunc: existDataAlicloudServiceMeshServiceMeshesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudServiceMeshServiceMeshesSourceNameMapFunc,
	}

	alicloudServiceMeshServiceMeshCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, statusRegexConf, allConf)
}
func testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {	
	default = "tf-testaccservicemeshservicemesh-%d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "^default-NODELETING$"
}

resource "alicloud_vpc" "default" {
	count = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
}

data "alicloud_vswitches" "default" {
	vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}

resource "alicloud_vswitch" "default" {
	count      = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
	vpc_id     = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
	cidr_block = cidrsubnet(length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.vpcs[0].cidr_block : alicloud_vpc.default[0].cidr_block, 8, 2)
	zone_id    = data.alicloud_zones.default.zones.0.id
}

data "alicloud_service_mesh_versions" "default" {
	edition = "Default"
}

resource "alicloud_service_mesh_service_mesh" "default" {
	service_mesh_name = var.name
	edition = "Default"
	version = reverse(data.alicloud_service_mesh_versions.default.versions).0.version
	network {
		vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
		vswitche_list = [length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id]
	}
}

data "alicloud_service_mesh_service_meshes" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
