package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudResourceManagerSharedTargetsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_shared_target.example.target_id}"]`,
			"resource_share_id": `"${alicloud_resource_manager_shared_target.example.resource_share_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"ids":               `["${alicloud_resource_manager_shared_target.example.target_id}_fake"]`,
			"resource_share_id": `"${alicloud_resource_manager_shared_target.example.resource_share_id}"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_target.example.resource_share_id}"`,
			"status":            `"Associated"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_target.example.resource_share_id}"`,
			"status":            `"Associating"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_target.example.resource_share_id}"`,
			"ids":               `["${alicloud_resource_manager_shared_target.example.target_id}"]`,
			"status":            `"Associated"`,
		}),
		fakeConfig: testAccCheckAlicloudResourceManagerSharedTargetsSourceConfig(rand, map[string]string{
			"resource_share_id": `"${alicloud_resource_manager_shared_target.example.resource_share_id}"`,
			"ids":               `["${alicloud_resource_manager_shared_target.example.target_id}"]`,
			"status":            `"Associating"`,
		}),
	}

	var existResourceManagerSharedTargetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"targets.#":                   "1",
			"ids.#":                       "1",
			"targets.0.id":                CHECKSET,
			"targets.0.target_id":         CHECKSET,
			"targets.0.resource_share_id": CHECKSET,
			"targets.0.status":            "Associated",
		}
	}

	var fakeResourceManagerSharedTargetsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"targets.#": "0",
			"ids.#":     "0",
		}
	}

	var resourceManagerShareTargetsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_resource_manager_shared_targets.example",
		existMapFunc: existResourceManagerSharedTargetsMapFunc,
		fakeMapFunc:  fakeResourceManagerSharedTargetsMapFunc,
	}
	preCheck := func() {
		testAccPreCheckEnterpriseAccountEnabled(t)
	}
	resourceManagerShareTargetsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, statusConf, allConf)

}

func testAccCheckAlicloudResourceManagerSharedTargetsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testaccResourceManagerSharedTargets%d"
}

data "alicloud_resource_manager_accounts" "example" {}

resource "alicloud_resource_manager_resource_share" "example" {
  resource_share_name = var.name
}

resource "alicloud_resource_manager_shared_target" "example" {
  resource_share_id = alicloud_resource_manager_resource_share.example.id
  target_id         = data.alicloud_resource_manager_accounts.example.ids.0
}

data "alicloud_resource_manager_shared_targets" "example" {
%s
}

`, rand, strings.Join(pairs, "\n   "))
	return config
}
