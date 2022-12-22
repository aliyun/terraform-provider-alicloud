package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudEcsElasticityAssuranceDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_elasticity_assurance.default.id}_fake"]`,
		}),
	}
	privatePoolOptionsIdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"private_pool_options_ids": `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"private_pool_options_ids": `["${alicloud_ecs_elasticity_assurance.default.id}_fake"]`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
			"tags": `{ 
						"Created" = "tfTestAcc0"
    					"For"     = "Tftestacc 0" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"ids": `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
			"tags": `{ 
						"Created" = "tfTestAcc0-fake"
    					"For"     = "Tftestacc 0-fake" 
					}`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"resource_group_id": `"${alicloud_ecs_elasticity_assurance.default.resource_group_id}"`,
			"ids":               `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"resource_group_id": `"${alicloud_ecs_elasticity_assurance.default.resource_group_id}_fake"`,
			"ids":               `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"status": `"${alicloud_ecs_elasticity_assurance.default.status}"`,
			"ids":    `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"status": `"Released"`,
			"ids":    `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"ids":                      `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
			"private_pool_options_ids": `["${alicloud_ecs_elasticity_assurance.default.id}"]`,
			"resource_group_id":        `"${alicloud_ecs_elasticity_assurance.default.resource_group_id}"`,
			"status":                   `"${alicloud_ecs_elasticity_assurance.default.status}"`,
			"tags": `{ 
						"Created" = "tfTestAcc0"
    					"For"     = "Tftestacc 0" 
					}`,
		}),
		fakeConfig: testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand, map[string]string{
			"ids":                      `["${alicloud_ecs_elasticity_assurance.default.id}_fake"]`,
			"private_pool_options_ids": `["${alicloud_ecs_elasticity_assurance.default.id}_fake"]`,
			"resource_group_id":        `"${alicloud_ecs_elasticity_assurance.default.resource_group_id}_fake"`,
			"status":                   `"Released"`,
			"tags": `{ 
						"Created" = "tfTestAcc0-fake"
    					"For"     = "Tftestacc 0-fake" 
					}`,
		}),
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithTime(t, []int{1})
	}

	EcsElasticityAssuranceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, privatePoolOptionsIdsConf, tagsConf, resourceGroupIdConf, statusConf, allConf)
}

var existEcsElasticityAssuranceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                              "1",
		"assurances.#":                       "1",
		"assurances.0.id":                    CHECKSET,
		"assurances.0.allocated_resources.#": "1",
		"assurances.0.allocated_resources.0.instance_type": CHECKSET,
		"assurances.0.allocated_resources.0.total_amount":  "1",
		"assurances.0.allocated_resources.0.used_amount":   CHECKSET,
		"assurances.0.allocated_resources.0.zone_id":       CHECKSET,
		"assurances.0.description":                         "before",
		"assurances.0.elasticity_assurance_id":             CHECKSET,
		"assurances.0.end_time":                            CHECKSET,
		"assurances.0.instance_charge_type":                CHECKSET,
		"assurances.0.private_pool_options_id":             CHECKSET,
		"assurances.0.private_pool_options_match_criteria": "Open",
		"assurances.0.private_pool_options_name":           "test_before",
		"assurances.0.resource_group_id":                   CHECKSET,
		"assurances.0.start_time":                          CHECKSET,
		"assurances.0.start_time_type":                     CHECKSET,
		"assurances.0.status":                              CHECKSET,
		"assurances.0.tags.%":                              "2",
		"assurances.0.tags.Created":                        "tfTestAcc0",
		"assurances.0.tags.For":                            "Tftestacc 0",
		"assurances.0.total_assurance_times":               CHECKSET,
		"assurances.0.used_assurance_times":                CHECKSET,
	}
}

var fakeEcsElasticityAssuranceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":        "0",
		"assurances.#": "0",
	}
}

var EcsElasticityAssuranceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ecs_elasticity_assurances.default",
	existMapFunc: existEcsElasticityAssuranceMapFunc,
	fakeMapFunc:  fakeEcsElasticityAssuranceMapFunc,
}

func testAccCheckAlicloudEcsElasticityAssuranceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccEcsElasticityAssurance%d"
}

data "alicloud_resource_manager_resource_groups" "default"{
	status = "OK"
}
data "alicloud_zones" "default" {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
	instance_type_family = "ecs.c6"
}

resource "alicloud_ecs_elasticity_assurance" "default" {
  instance_amount = 1
  description     = "before"
  zone_ids = [data.alicloud_zones.default.zones[0].id]
  private_pool_options_name           = "test_before"
  period                              = 1
  private_pool_options_match_criteria = "Open"
  instance_type = [data.alicloud_instance_types.default.instance_types.0.id]
  period_unit     = "Month"
  assurance_times = "Unlimited"
  tags = {
		Created =  "tfTestAcc0"
		For =      "Tftestacc 0"
	}
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

data "alicloud_ecs_elasticity_assurances" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
