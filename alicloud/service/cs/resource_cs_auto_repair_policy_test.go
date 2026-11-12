package cs_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("alicloud_cs_auto_repair_policy", &resource.Sweeper{
		Name: "alicloud_cs_auto_repair_policy",
		F:    testSweepCSAutoRepairPolicy,
	})
}

func testSweepCSAutoRepairPolicy(region string) error {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting Alicloud client: %s", err)
	}
	client := rawClient.(*connectivity.AliyunClient)

	csClient, err := client.NewRoaCsClient()
	if err != nil {
		return fmt.Errorf("error getting CS client: %s", err)
	}
	clusters, err := csClient.DescribeClusters(nil)
	if err != nil {
		return fmt.Errorf("error listing CS clusters: %s", err)
	}
	if clusters == nil || clusters.Body == nil {
		return nil
	}

	prefixes := []string{"tf-testAcc", "tf_testAcc", "tf-example", "tf_example", "terraform-example"}
	for _, cluster := range clusters.Body {
		clusterId := tea.StringValue(cluster.ClusterId)
		if clusterId == "" {
			continue
		}
		policies, err := csClient.ListAutoRepairPolicies(tea.String(clusterId))
		if err != nil {
			log.Printf("[ERROR] Failed to list auto repair policies of cluster %s: %s", clusterId, err)
			continue
		}
		if policies == nil || policies.Body == nil {
			continue
		}
		for _, policy := range policies.Body.Items {
			if tea.StringValue(policy.Type) == "builtin" {
				continue
			}
			name := tea.StringValue(policy.Name)
			if !sweep.SweepAll() {
				matched := false
				for _, prefix := range prefixes {
					if strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}
			log.Printf("[INFO] Deleting auto repair policy: %s (%s) of cluster %s", name, tea.StringValue(policy.Id), clusterId)
			if _, err := csClient.DeleteAutoRepairPolicy(tea.String(clusterId), policy.Id); err != nil {
				log.Printf("[ERROR] Failed to delete auto repair policy %s: %s", name, err)
			}
		}
	}
	return nil
}

func TestAccAliCloudCSAutoRepairPolicy_basic(t *testing.T) {
	resourceId := "alicloud_cs_auto_repair_policy.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{"policy_id": acctest.CHECKSET}).ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoRepairPolicy-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAutoRepairPolicyConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":        "${alicloud_cs_managed_kubernetes.default.id}",
					"name":              "${var.name}",
					"resource_sub_type": "ess",
					"rules":             autoRepairPolicyRules(true),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cluster_id":                        acctest.CHECKSET,
						"name":                              name,
						"resource_type":                     "nodepool",
						"resource_sub_type":                 "ess",
						"type":                              "custom",
						"resource_ids.#":                    "0",
						"rules.#":                           "1",
						"rules.0.incidents.#":               "1",
						"rules.0.incidents.0.name":          "NTPProblem",
						"rules.0.repair_procedure.#":        "4",
						"rules.0.repair_procedure.1.name":   "Evict",
						"rules.0.repair_procedure.1.config": `{"evict_timeout":300}`,
						"rules.0.repair_procedure.2.intervention.#":                   "1",
						"rules.0.repair_procedure.2.intervention.0.enable":            "true",
						"rules.0.repair_procedure.2.intervention.0.type":              "label",
						"rules.0.repair_procedure.2.intervention.0.inquiring_label.#": "1",
						"rules.0.repair_procedure.2.intervention.0.approved_label.#":  "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":  "${var.name}-updated",
					"rules": autoRepairPolicyRules(false),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                                                        name + "-updated",
						"rules.0.repair_procedure.#":                                  "4",
						"rules.0.repair_procedure.1.config":                           `{"evict_timeout":600}`,
						"rules.0.repair_procedure.2.intervention.#":                   "0",
						"rules.0.repair_procedure.2.intervention.0.enable":            acctest.REMOVEKEY,
						"rules.0.repair_procedure.2.intervention.0.type":              acctest.REMOVEKEY,
						"rules.0.repair_procedure.2.intervention.0.inquiring_label.#": acctest.REMOVEKEY,
						"rules.0.repair_procedure.2.intervention.0.approved_label.#":  acctest.REMOVEKEY,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// collection-order: bypass rules.repair_procedure -- the service validates the configured procedures as a strict subsequence of the built-in flow order, so only the built-in order is accepted and every permutation is rejected
func TestAccAliCloudCSAutoRepairPolicy_ruleMemberOrder(t *testing.T) {
	resourceId := "alicloud_cs_auto_repair_policy.default"
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoRepairPolicy-order-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAutoRepairPolicyConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"cluster_id":        "${alicloud_cs_managed_kubernetes.default.id}",
					"name":              "${var.name}",
					"resource_sub_type": "ess",
					"rules": []map[string]interface{}{
						{
							"incidents": []map[string]interface{}{
								{"name": "NTPProblem"},
								{"name": "RuntimeOffline"},
							},
							"repair_procedure": []map[string]interface{}{
								{"name": "QuarantineNode"},
								{"name": "Evict", "config": "${jsonencode({evict_timeout = 300})}"},
								{"name": "Repair"},
								{"name": "UnquarantineNode"},
							},
						},
						{
							"incidents": []map[string]interface{}{
								{"name": "Ready"},
							},
							"repair_procedure": []map[string]interface{}{
								{"name": "QuarantineNode"},
								{"name": "Evict", "config": "${jsonencode({evict_timeout = 300})}"},
								{"name": "Repair"},
								{"name": "UnquarantineNode"},
							},
						},
					},
				}),
			},
			{
				// incident order carries no meaning: reordering the incidents of the
				// first rule must converge without changing anything else
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"incidents": []map[string]interface{}{
								{"name": "RuntimeOffline"},
								{"name": "NTPProblem"},
							},
							"repair_procedure": []map[string]interface{}{
								{"name": "QuarantineNode"},
								{"name": "Evict", "config": "${jsonencode({evict_timeout = 300})}"},
								{"name": "Repair"},
								{"name": "UnquarantineNode"},
							},
						},
						{
							"incidents": []map[string]interface{}{
								{"name": "Ready"},
							},
							"repair_procedure": []map[string]interface{}{
								{"name": "QuarantineNode"},
								{"name": "Evict", "config": "${jsonencode({evict_timeout = 300})}"},
								{"name": "Repair"},
								{"name": "UnquarantineNode"},
							},
						},
					},
				}),
			},
			{
				// rules are identified by their content, so swapping the two rules
				// must converge as well
				Config: testAccConfig(map[string]interface{}{
					"rules": []map[string]interface{}{
						{
							"incidents": []map[string]interface{}{
								{"name": "Ready"},
							},
							"repair_procedure": []map[string]interface{}{
								{"name": "QuarantineNode"},
								{"name": "Evict", "config": "${jsonencode({evict_timeout = 300})}"},
								{"name": "Repair"},
								{"name": "UnquarantineNode"},
							},
						},
						{
							"incidents": []map[string]interface{}{
								{"name": "RuntimeOffline"},
								{"name": "NTPProblem"},
							},
							"repair_procedure": []map[string]interface{}{
								{"name": "QuarantineNode"},
								{"name": "Evict", "config": "${jsonencode({evict_timeout = 300})}"},
								{"name": "Repair"},
								{"name": "UnquarantineNode"},
							},
						},
					},
				}),
			},
		},
	})
}

func TestAccAliCloudCSAutoRepairPolicy_nodePoolBinding(t *testing.T) {
	resourceId := "alicloud_cs_kubernetes_node_pool.default"
	testAccCheck := acctest.ResourceAttrInit(resourceId, map[string]string{}).ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccCSAutoRepairPolicy-%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, resourceCSAutoRepairPolicyNodePoolConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(autoRepairPolicyNodePoolConfig("${alicloud_cs_auto_repair_policy.default.policy_id}", "2", "10", true)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"management.0.auto_repair":                                        "true",
						"management.0.auto_repair_policy.0.auto_repair_policy_id":         acctest.CHECKSET,
						"management.0.auto_repair_policy.0.max_parallel_repairing_nodes":  "2",
						"management.0.auto_repair_policy.0.max_unhealthy_nodes_threshold": "10",
					}),
				),
			},
			{
				// re-applying refreshes the policy, which then reports the node pool
				// it is bound to
				Config: testAccConfig(autoRepairPolicyNodePoolConfig("${alicloud_cs_auto_repair_policy.default.policy_id}", "2", "10", true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cs_auto_repair_policy.default", "resource_ids.#", "1"),
					resource.TestCheckResourceAttrPair("alicloud_cs_auto_repair_policy.default", "resource_ids.0", "alicloud_cs_kubernetes_node_pool.default", "node_pool_id"),
				),
			},
			{
				// turning auto repair off does not clear the policy id; the policy
				// merely stops taking effect
				Config: testAccConfig(autoRepairPolicyNodePoolConfig("${alicloud_cs_auto_repair_policy.default.policy_id}", "2", "10", false)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"management.0.auto_repair":                                        "false",
						"management.0.auto_repair_policy.0.max_parallel_repairing_nodes":  "2",
						"management.0.auto_repair_policy.0.max_unhealthy_nodes_threshold": "10",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "management.0.auto_repair_policy.0.auto_repair_policy_id", "alicloud_cs_auto_repair_policy.default", "policy_id"),
				),
			},
			{
				// re-applying refreshes the policy, which stops reporting the node
				// pool while auto repair is off
				Config: testAccConfig(autoRepairPolicyNodePoolConfig("${alicloud_cs_auto_repair_policy.default.policy_id}", "2", "10", false)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceId, "management.0.auto_repair_policy.0.auto_repair_policy_id", "alicloud_cs_auto_repair_policy.default", "policy_id"),
					resource.TestCheckResourceAttr("alicloud_cs_auto_repair_policy.default", "resource_ids.#", "0"),
				),
			},
			{
				// the retained policy takes effect again once auto repair is back on
				Config: testAccConfig(autoRepairPolicyNodePoolConfig("${alicloud_cs_auto_repair_policy.default.policy_id}", "2", "10", true)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"management.0.auto_repair":                                        "true",
						"management.0.auto_repair_policy.0.max_parallel_repairing_nodes":  "2",
						"management.0.auto_repair_policy.0.max_unhealthy_nodes_threshold": "10",
					}),
					resource.TestCheckResourceAttrPair(resourceId, "management.0.auto_repair_policy.0.auto_repair_policy_id", "alicloud_cs_auto_repair_policy.default", "policy_id"),
				),
			},
			{
				// the policy reports the node pool again on the next refresh
				Config: testAccConfig(autoRepairPolicyNodePoolConfig("${alicloud_cs_auto_repair_policy.default.policy_id}", "2", "10", true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alicloud_cs_auto_repair_policy.default", "resource_ids.#", "1"),
					resource.TestCheckResourceAttrPair("alicloud_cs_auto_repair_policy.default", "resource_ids.0", resourceId, "node_pool_id"),
				),
			},
		},
	})
}

func autoRepairPolicyRules(withIntervention bool) []map[string]interface{} {
	repair := map[string]interface{}{"name": "Repair"}
	if withIntervention {
		repair["intervention"] = []map[string]interface{}{
			{
				"enable": true,
				"type":   "label",
				"inquiring_label": []map[string]interface{}{
					{"key": "intervention.autorepair.node.k8s.aliyun.com/repair", "value": "inquiring"},
				},
				"approved_label": []map[string]interface{}{
					{"key": "intervention.autorepair.node.k8s.aliyun.com/repair", "value": "approved"},
				},
			},
		}
	}
	evictTimeout := "300"
	if !withIntervention {
		evictTimeout = "600"
	}
	return []map[string]interface{}{
		{
			"incidents": []map[string]interface{}{
				{"name": "NTPProblem"},
			},
			"repair_procedure": []map[string]interface{}{
				{"name": "QuarantineNode"},
				{"name": "Evict", "config": fmt.Sprintf("${jsonencode({evict_timeout = %s})}", evictTimeout)},
				repair,
				{"name": "UnquarantineNode"},
			},
		},
	}
}

func autoRepairPolicyNodePoolConfig(policyId, maxParallelRepairingNodes, maxUnhealthyNodesThreshold string, autoRepair bool) map[string]interface{} {
	management := map[string]interface{}{
		"auto_repair": autoRepair,
	}
	if policyId != "" || maxParallelRepairingNodes != "" || maxUnhealthyNodesThreshold != "" {
		management["auto_repair_policy"] = []map[string]interface{}{
			{
				"auto_repair_policy_id":         policyId,
				"max_parallel_repairing_nodes":  maxParallelRepairingNodes,
				"max_unhealthy_nodes_threshold": maxUnhealthyNodesThreshold,
			},
		}
	}
	return map[string]interface{}{
		"name":                 "${var.name}",
		"cluster_id":           "${alicloud_cs_managed_kubernetes.default.id}",
		"vswitch_ids":          []string{"${alicloud_vswitch.default.id}"},
		"instance_types":       []string{"${data.alicloud_instance_types.default.instance_types.0.id}"},
		"key_name":             "${alicloud_key_pair.default.key_name}",
		"system_disk_category": "cloud_efficiency",
		"system_disk_size":     40,
		"desired_size":         1,
		"management":           []map[string]interface{}{management},
		"depends_on":           []string{"alicloud_cs_kubernetes_addon.nlc", "alicloud_cs_kubernetes_addon.npd"},
	}
}

func resourceCSAutoRepairPolicyConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  cluster_spec                 = "ack.pro.small"
  vswitch_ids                  = [alicloud_vswitch.default.id]
  pod_vswitch_ids              = [alicloud_vswitch.default.id]
  service_cidr                 = cidrsubnet("172.16.0.0/16", 4, 3)
  new_nat_gateway              = false
  is_enterprise_security_group = true
  slb_internet_enabled         = false
  deletion_protection          = false

  addons {
    name = "terway-eniip"
  }
}
`, name)
}

func resourceCSAutoRepairPolicyNodePoolConfigDependence(name string) string {
	return resourceCSAutoRepairPolicyConfigDependence(name) + fmt.Sprintf(`
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_efficiency"
  sorted_by            = "Price"
}

resource "alicloud_key_pair" "default" {
  key_name = var.name
}

# Auto repair requires both add-ons: binding an auto repair policy requires
# node-lifecycle-controller at v2.3.13-b581b592 or later, and automatic fault
# diagnosis requires ack-node-problem-detector at 1.2.41 or later, so both are
# pinned to their minimum supported versions here.
resource "alicloud_cs_kubernetes_addon" "nlc" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  name       = "node-lifecycle-controller"
  version    = "v2.3.13-b581b592"
}

resource "alicloud_cs_kubernetes_addon" "npd" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  name       = "ack-node-problem-detector"
  version    = "1.2.41"
}

resource "alicloud_cs_auto_repair_policy" "default" {
  cluster_id        = alicloud_cs_managed_kubernetes.default.id
  name              = "${var.name}-policy"
  resource_sub_type = "ess"

  rules {
    incidents {
      name = "NTPProblem"
    }
    repair_procedure {
      name = "QuarantineNode"
    }
    repair_procedure {
      name   = "Evict"
      config = jsonencode({ evict_timeout = 300 })
    }
    repair_procedure {
      name = "Repair"
    }
    repair_procedure {
      name = "UnquarantineNode"
    }
  }
}
`)
}
