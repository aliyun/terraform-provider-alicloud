package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerSharedResourcesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `[format("%s:%s",alicloud_resource_manager_shared_resource.example.resource_id,alicloud_resource_manager_shared_resource.example.resource_type)]`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `["fake"]`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `[format("%s:%s",alicloud_resource_manager_shared_resource.example.resource_id,alicloud_resource_manager_shared_resource.example.resource_type)]`,
			"status":            `"Associated"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.example.resource_share_id}"`,
			"ids":               `["fake"]`,
			"status":            `"Associating"`,
		}),
	}

	var existResourceManagerSharedResourcesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#":                   "1",
			"ids.#":                         "1",
			"resources.0.id":                CHECKSET,
			"resources.0.resource_type":     "VSwitch",
			"resources.0.resource_share_id": CHECKSET,
			"resources.0.resource_id":       CHECKSET,
			"resources.0.status":            "Associated",
		}
	}

	var fakeResourceManagerSharedResourcessMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"resources.#": "0",
			"ids.#":       "0",
		}
	}

	var resourceManagerSharedResourcesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_shared_resources.example",
		existMapFunc: existResourceManagerSharedResourcesMapFunc,
		fakeMapFunc:  fakeResourceManagerSharedResourcessMapFunc,
	}
	preCheck := func() {
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	resourceManagerSharedResourcesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf)

}

func testAccCheckAlicloudResourceManagerSharedResourcesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccResourceManagerSharedResources%d"
}
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}
data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.example.ids.0
}

resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.example.ids.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_resource_manager_resource_share" "example" {
	resource_share_name = var.name
}
resource "alicloud_resource_manager_shared_resource" "example" {
  resource_id       = local.vswitch_id
  resource_share_id = alicloud_resource_manager_resource_share.example.id
  resource_type     = "VSwitch"
}
data "alicloud_resource_manager_shared_resources" "example"{
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
