package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudResourceManagerSharedResourcesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"ids": `[format("%s:%s",alicloud_resource_manager_shared_resource.default.resource_id,alicloud_resource_manager_shared_resource.default.resource_type)]`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"ids": `["fake"]`,
		}),
	}
	resourceShareIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.default.resource_share_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_resource_share.fake.id}"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"status": `"Associated"`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"status": `"Associating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"ids":               `[format("%s:%s",alicloud_resource_manager_shared_resource.default.resource_id,alicloud_resource_manager_shared_resource.default.resource_type)]`,
			"resource_share_id": `"${alicloud_resource_manager_shared_resource.default.resource_share_id}"`,
			"status":            `"Associated"`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand, map[string]string{
			"ids":               `["fake"]`,
			"resource_share_id": `"${alicloud_resource_manager_resource_share.fake.id}"`,
			"status":            `"Associating"`,
		}),
	}
	var existResourceManagerSharedResourcesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"resources.#":                   "1",
			"resources.0.id":                CHECKSET,
			"resources.0.resource_id":       CHECKSET,
			"resources.0.resource_type":     "VSwitch",
			"resources.0.resource_share_id": CHECKSET,
			"resources.0.status":            "Associated",
		}
	}
	var fakeResourceManagerSharedResourcessMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"resources.#": "0",
		}
	}
	var resourceManagerSharedResourcesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_shared_resources.default",
		existMapFunc: existResourceManagerSharedResourcesMapFunc,
		fakeMapFunc:  fakeResourceManagerSharedResourcessMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	resourceManagerSharedResourcesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceShareIdConf, statusConf, allConf)
}

func testAccCheckAliCloudResourceManagerSharedResourcesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testaccResourceManagerSharedResources%d"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_vpcs" "default" {
  		name_regex = "^default-NODELETING$"
	}

	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
  		zone_id = data.alicloud_zones.default.ids.0
	}

	resource "alicloud_vswitch" "vswitch" {
  		count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  		vpc_id       = data.alicloud_vpcs.default.ids.0
  		cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  		zone_id      = data.alicloud_zones.default.ids.0
  		vswitch_name = var.name
	}

	locals {
  		vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}

	resource "alicloud_resource_manager_resource_share" "default" {
  		resource_share_name = var.name
	}

	resource "alicloud_resource_manager_resource_share" "fake" {
  		resource_share_name = var.name
	}

	resource "alicloud_resource_manager_shared_resource" "default" {
  		resource_id       = local.vswitch_id
  		resource_share_id = alicloud_resource_manager_resource_share.default.id
  		resource_type     = "VSwitch"
	}

	data "alicloud_resource_manager_shared_resources" "default" {
		%s
	}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
