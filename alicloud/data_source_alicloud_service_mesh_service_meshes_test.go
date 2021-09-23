package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudServiceMeshServiceMeshesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1, 1000)

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
			"name_regex": `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"name_regex": `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}_fake"`,
		}),
	}

	statusRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_service_mesh_service_mesh.default.id}"]`,
			"status": `"running"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":    `["${alicloud_service_mesh_service_mesh.default.id}"]`,
			"status": `"initial"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_service_mesh_service_mesh.default.id}"]`,
			"name_regex": `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}"`,
			"status":     `"running"`,
		}),
		fakeConfig: testAccCheckAlicloudServiceMeshServiceMeshDataSourceName(rand, map[string]string{
			"ids":        `["${alicloud_service_mesh_service_mesh.default.id}_fake"]`,
			"name_regex": `"${alicloud_service_mesh_service_mesh.default.service_mesh_name}_fake"`,
			"status":     `"initial"`,
		}),
	}

	var existDataAlicloudServiceMeshServiceMeshesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"meshes.#":                   "1",
			"meshes.0.status":            "running",
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
 name_regex = "default-NODELETING"
}
resource "alicloud_vpc" "default" {
    count = length(data.alicloud_vpcs.default.ids) > 0 ? 0 : 1
	vpc_name = var.name
}
data "alicloud_vswitches" "default" {
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}
resource "alicloud_vswitch" "default" {
  count         = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id        = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
  cidr_block    = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 2)
  zone_id     	= data.alicloud_zones.default.zones.0.id
  vswitch_name  = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : alicloud_vswitch.default[0].id
  vpc_id = length(data.alicloud_vpcs.default.ids) > 0 ? data.alicloud_vpcs.default.ids[0] : alicloud_vpc.default[0].id
}

resource "alicloud_service_mesh_service_mesh" "default" {
	service_mesh_name = var.name
	network {
		vpc_id = local.vpc_id
		vswitche_list = [local.vswitch_id]
	}
}

data "alicloud_service_mesh_service_meshes" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
