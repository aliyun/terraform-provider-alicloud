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
	policyName := name + "-policy"
	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, name, resourceAutoRepairPoliciesDataSourceDependence)

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
						"names.0":                      policyName,
						"policies.#":                   "1",
						"policies.0.name":              policyName,
						"policies.0.type":              "custom",
						"policies.0.resource_type":     "nodepool",
						"policies.0.resource_sub_type": "ess",
						"policies.0.rules":             acctest.REGEXMATCH + `NTPProblem`,
					}),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", "alicloud_cs_auto_repair_policy.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "policies.0.policy_id", "alicloud_cs_auto_repair_policy.default", "policy_id"),
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
						"policies.0.resource_ids.#":    "0",
						"policies.0.rules":             acctest.REMOVEKEY,
					}),
				),
			},
			{
				// ids selects one policy by its <cluster_id>:<policy_id> identity
				Config: testAccConfig(map[string]interface{}{
					"cluster_id": "${alicloud_cs_managed_kubernetes.default.id}",
					"ids":        []string{"${alicloud_cs_auto_repair_policy.default.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ids.#":                        "1",
						"names.#":                      "1",
						"names.0":                      policyName,
						"policies.#":                   "1",
						"policies.0.name":              policyName,
						"policies.0.type":              "custom",
						"policies.0.resource_sub_type": "ess",
						"policies.0.resource_ids.#":    "1",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "ids.0", "alicloud_cs_auto_repair_policy.default", "id"),
					resource.TestCheckResourceAttrPair(resourceId, "policies.0.policy_id", "alicloud_cs_auto_repair_policy.default", "policy_id"),
					resource.TestCheckResourceAttrPair(resourceId, "policies.0.resource_ids.0", "alicloud_cs_kubernetes_node_pool.default", "node_pool_id"),
				),
			},
			{
				// nodepool_id selects the policies bound to that node pool
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":  "${alicloud_cs_managed_kubernetes.default.id}",
					"nodepool_id": "${alicloud_cs_kubernetes_node_pool.default.node_pool_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ids.#":                     "1",
						"names.0":                   policyName,
						"policies.#":                "1",
						"policies.0.name":           policyName,
						"policies.0.resource_ids.#": "1",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "policies.0.resource_ids.0", "alicloud_cs_kubernetes_node_pool.default", "node_pool_id"),
				),
			},
			{
				// the filters intersect, so a query that matches no policy returns an
				// empty set instead of every policy
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
						"names.0":                      acctest.REMOVEKEY,
						"policies.0.name":              acctest.REMOVEKEY,
						"policies.0.type":              acctest.REMOVEKEY,
						"policies.0.resource_sub_type": acctest.REMOVEKEY,
						"policies.0.resource_ids.#":    acctest.REMOVEKEY,
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

func resourceAutoRepairPoliciesDataSourceDependence(name string) string {
	return resourceAutoRepairPolicyNodePoolConfigDependence(name) + `
# The policy is bound to this node pool so that the nodepool_id filter has a
# positive match to select.
resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "${var.name}-np"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  key_name             = alicloud_key_pair.default.key_name
  system_disk_category = "cloud_essd"
  system_disk_size     = 40
  desired_size         = 1

  management {
    auto_repair = true

    auto_repair_policy {
      auto_repair_policy_id = alicloud_cs_auto_repair_policy.default.policy_id
    }
  }

  depends_on = [alicloud_cs_kubernetes_addon.nlc, alicloud_cs_kubernetes_addon.npd]
}
`
}
