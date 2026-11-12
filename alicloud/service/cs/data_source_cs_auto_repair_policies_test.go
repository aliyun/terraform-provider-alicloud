package cs_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCSAutoRepairPoliciesDataSource(t *testing.T) {
	resourceId := "data.alicloud_cs_auto_repair_policies.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoRepairPolicies-%d", rand)
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, resourceCSAutoRepairPoliciesDataSourceDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
					"name_regex": "${alicloud_cs_auto_repair_policy.default.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                   acctest.CHECKSET,
						"ids.#":                        "1",
						"names.#":                      "1",
						"names.0":                      name,
						"policies.#":                   "1",
						"policies.0.name":              name,
						"policies.0.type":              "custom",
						"policies.0.resource_type":     "nodepool",
						"policies.0.resource_sub_type": "ess",
						"policies.0.resource_ids.#":    "0",
						"policies.0.rules":             acctest.REGEXMATCH + `NTPProblem`,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":        "${alicloud_cs_managed_kubernetes.default.id}",
					"type":              "builtin",
					"resource_sub_type": "ess",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                   acctest.CHECKSET,
						"ids.#":                        acctest.REGEXMATCH + "^[1-9][0-9]*$",
						"names.#":                      acctest.REMOVEKEY,
						"names.0":                      acctest.REMOVEKEY,
						"policies.#":                   acctest.REMOVEKEY,
						"policies.0.name":              acctest.REMOVEKEY,
						"policies.0.type":              "builtin",
						"policies.0.resource_type":     acctest.REMOVEKEY,
						"policies.0.resource_sub_type": acctest.REGEXMATCH + "ess",
						"policies.0.resource_ids.#":    acctest.REMOVEKEY,
						"policies.0.rules":             acctest.REMOVEKEY,
					}),
				),
			},
			{
				// ids and nodepool_id narrow the result set further; filters that
				// match nothing return an empty set instead of every policy
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":  "${alicloud_cs_managed_kubernetes.default.id}",
					"name_regex":  "^tf-testAccCSAutoRepairPolicies-not-exist$",
					"ids":         []string{"${alicloud_cs_managed_kubernetes.default.id}:r-not-exist"},
					"nodepool_id": "np-not-exist",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ids.#":                        "0",
						"names.#":                      "0",
						"policies.#":                   "0",
						"policies.0.type":              acctest.REMOVEKEY,
						"policies.0.resource_sub_type": acctest.REMOVEKEY,
					}),
				),
			},
			{
				// changing the filter values re-runs the read
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":  "${alicloud_cs_managed_kubernetes.default.id}",
					"name_regex":  "^tf-testAccCSAutoRepairPolicies-still-not-exist$",
					"ids":         []string{"${alicloud_cs_managed_kubernetes.default.id}:r-still-not-exist"},
					"nodepool_id": "np-still-not-exist",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ids.#":      "0",
						"names.#":    "0",
						"policies.#": "0",
					}),
				),
			},
			{
				// cluster_id is required and every read is scoped to it, so an
				// unknown cluster must surface the service error
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "c-not-exist",
				}),
				ExpectError: regexp.MustCompile(`(?i)not found`),
			},
		},
	})
}

func resourceCSAutoRepairPoliciesDataSourceDependence(name string) string {
	return resourceCSAutoRepairPolicyConfigDependence(name) + `
resource "alicloud_cs_auto_repair_policy" "default" {
  cluster_id        = alicloud_cs_managed_kubernetes.default.id
  name              = var.name
  resource_sub_type = "ess"

  rules {
    incidents {
      name = "NTPProblem"
    }
    repair_procedure {
      name = "QuarantineNode"
    }
    repair_procedure {
      name = "Evict"
    }
    repair_procedure {
      name = "Repair"
    }
    repair_procedure {
      name = "UnquarantineNode"
    }
  }
}
`
}
