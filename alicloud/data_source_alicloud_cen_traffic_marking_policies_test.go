package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCenTrafficMarkingPoliciesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_traffic_marking_policy.default.traffic_marking_policy_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"name_regex": `"${alicloud_cen_traffic_marking_policy.default.traffic_marking_policy_name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_traffic_marking_policy.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_cen_traffic_marking_policy.default.id}_fake"]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cen_traffic_marking_policy.default.id}"]`,
			"description": `"${alicloud_cen_traffic_marking_policy.default.description}"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cen_traffic_marking_policy.default.id}"]`,
			"description": `"${alicloud_cen_traffic_marking_policy.default.description}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cen_traffic_marking_policy.default.id}"]`,
			"status": `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids":    `["${alicloud_cen_traffic_marking_policy.default.id}"]`,
			"status": `"Creating"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cen_traffic_marking_policy.default.id}"]`,
			"name_regex":  `"${alicloud_cen_traffic_marking_policy.default.traffic_marking_policy_name}"`,
			"description": `"${alicloud_cen_traffic_marking_policy.default.description}"`,
			"status":      `"Active"`,
		}),
		fakeConfig: testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand, map[string]string{
			"ids":         `["${alicloud_cen_traffic_marking_policy.default.id}_fake"]`,
			"name_regex":  `"${alicloud_cen_traffic_marking_policy.default.traffic_marking_policy_name}_fake"`,
			"description": `"${alicloud_cen_traffic_marking_policy.default.description}_fake"`,
			"status":      `"Creating"`,
		}),
	}

	var existCenTrafficMarkingPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"policies.#":                             "1",
			"policies.0.priority":                    "1",
			"policies.0.description":                 fmt.Sprintf("tf-testAccCenTrafficMarkingPolicies%d", rand),
			"policies.0.traffic_marking_policy_name": fmt.Sprintf("tf-testAccCenTrafficMarkingPolicies%d", rand),
			"policies.0.traffic_marking_policy_id":   CHECKSET,
			"policies.0.marking_dscp":                "1",
			"policies.0.transit_router_id":           CHECKSET,
			"policies.0.status":                      "Active",
			"policies.0.id":                          CHECKSET,
		}
	}

	var fakeCenTrafficMarkingPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}

	var cenTrafficMarkingPoliciesCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cen_traffic_marking_policies.default",
		existMapFunc: existCenTrafficMarkingPoliciesMapFunc,
		fakeMapFunc:  fakeCenTrafficMarkingPoliciesMapFunc,
	}
	preCheck := func() {
		testAccPreCheck(t)
	}

	cenTrafficMarkingPoliciesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, descriptionConf, allConf)

}

func testAccCheckAlicloudCenTrafficMarkingPoliciesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCenTrafficMarkingPolicies%d"
}

resource "alicloud_cen_instance" "default" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "default" {
  cen_id              = alicloud_cen_instance.default.id
  transit_router_name = var.name
}

resource "alicloud_cen_traffic_marking_policy" "default" {
  marking_dscp                = 1
  priority                    = 1
  traffic_marking_policy_name = var.name
  description                 = var.name
  transit_router_id           = alicloud_cen_transit_router.default.transit_router_id
}

data "alicloud_cen_traffic_marking_policies" "default" {
	transit_router_id = alicloud_cen_transit_router.default.transit_router_id
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
