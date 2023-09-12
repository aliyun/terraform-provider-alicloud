package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudResourceManagerSharedTargetsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_shared_target.default.target_id}"]`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_resource_manager_shared_target.default.target_id}_fake"]`,
		}),
	}
	resourceShareIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_target.default.resource_share_id}"`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_resource_share.fake.id}"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"status": `"Associated"`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"status": `"Associating"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_shared_target.default.target_id}"]`,
			"resource_share_id": `"${alicloud_resource_manager_shared_target.default.resource_share_id}"`,
			"status":            `"Associated"`,
		}),
		fakeConfig: testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_shared_target.default.target_id}_fake"]`,
			"resource_share_id": `"${alicloud_resource_manager_resource_share.fake.id}"`,
			"status":            `"Associating"`,
		}),
	}
	var existResourceManagerSharedTargetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"targets.#":                   "1",
			"targets.0.id":                CHECKSET,
			"targets.0.target_id":         CHECKSET,
			"targets.0.resource_share_id": CHECKSET,
			"targets.0.status":            "Associated",
		}
	}
	var fakeResourceManagerSharedTargetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"targets.#": "0",
		}
	}

	var resourceManagerShareTargetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_shared_targets.default",
		existMapFunc: existResourceManagerSharedTargetsMapFunc,
		fakeMapFunc:  fakeResourceManagerSharedTargetsMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}
	resourceManagerShareTargetsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, resourceShareIdConf, statusConf, allConf)
}

func testAccCheckAliCloudResourceManagerSharedTargetsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
	variable "name" {
  		default = "tf-testaccResourceManagerSharedTargets%d"
	}

	data "alicloud_resource_manager_accounts" "default" {
	}

	resource "alicloud_resource_manager_resource_share" "default" {
  		resource_share_name = var.name
	}

	resource "alicloud_resource_manager_resource_share" "fake" {
  		resource_share_name = var.name
	}

	resource "alicloud_resource_manager_shared_target" "default" {
  		resource_share_id = alicloud_resource_manager_resource_share.default.id
  		target_id         = data.alicloud_resource_manager_accounts.default.ids.0
	}

	data "alicloud_resource_manager_shared_targets" "default" {
		%s
	}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
