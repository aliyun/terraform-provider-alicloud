package cs_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cs"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Ack PolicyInstance. >>> Resource test cases, automatically generated.
func TestAccAliCloudCSKubernetesPolicyInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceID := "alicloud_cs_kubernetes_policy_instance.default"
	ra := acctest.ResourceAttrInit(resourceID, AliCloudCSKubernetesPolicyInstanceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceID, &v, func() interface{} {
		return cs.NewAckServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAckPolicyInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, name, AliCloudCSKubernetesPolicyInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"action":     "deny",
					"namespaces": []string{},
					"parameters": map[string]interface{}{
						"hostNetwork": "false",
						"min":         30,
						"max":         300,
					},
					"cluster_id":  "${alicloud_cs_managed_kubernetes.CreateCluster.id}",
					"policy_name": "${var.policy_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":                 "deny",
						"namespaces.#":           "0",
						"parameters.%":           "3",
						"parameters.hostNetwork": "false",
						"parameters.min":         "30",
						"parameters.max":         "300",
						"cluster_id":             acctest.CHECKSET,
						"policy_name":            acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action":     "warn",
					"namespaces": []string{"test", "test1", "test2"},
					"parameters": map[string]interface{}{
						"hostNetwork": "true",
						"min":         50,
						"max":         500,
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":                 "warn",
						"namespaces.#":           "3",
						"parameters.%":           "3",
						"parameters.hostNetwork": "true",
						"parameters.min":         "50",
						"parameters.max":         "500",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{"test"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AliCloudCSKubernetesPolicyInstanceMap = map[string]string{
	"instance_name": acctest.CHECKSET,
}

func AliCloudCSKubernetesPolicyInstanceBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "vpc_cidr" {
  default = "10.0.0.0/8"
}

variable "vswitch_cidrs" {
  type    = list(string)
  default = ["10.1.0.0/16", "10.2.0.0/16"]
}

variable "pod_cidr" {
  default = "172.16.0.0/16"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

variable "policy_name" {
  default = "ACKPSPHostNetworkingPorts"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {}

resource "alicloud_vpc" "CreateVPC" {
  cidr_block = var.vpc_cidr
}

# According to the vswitch cidr blocks to launch several vswitches
resource "alicloud_vswitch" "CreateVSwitch" {
  count      = length(var.vswitch_cidrs)
  vpc_id     = alicloud_vpc.CreateVPC.id
  cidr_block = element(var.vswitch_cidrs, count.index)
  zone_id    = data.alicloud_enhanced_nat_available_zones.enhanced.zones[count.index].zone_id
}

resource "alicloud_cs_managed_kubernetes" "CreateCluster" {
  name                         = var.name
  cluster_spec                 = "ack.standard"
  profile                      = "Default"
  vswitch_ids                  = split(",", join(",", alicloud_vswitch.CreateVSwitch.*.id))
  pod_cidr                     = var.pod_cidr
  service_cidr                 = var.service_cidr
  is_enterprise_security_group = true
  ip_stack                     = "ipv4"
  proxy_mode                   = "ipvs"
  deletion_protection          = false

  addons {
    name = "gatekeeper"
  }
  addons {
    name = "loongcollector"
  }
  addons {
    name = "policy-template-controller"
  }

  operation_policy {
    cluster_auto_upgrade {
      enabled = false
    }
  }
  maintenance_window {
    enable = false
  }
}

resource "alicloud_cs_kubernetes_policy_instance" "test" {
  cluster_id  = alicloud_cs_managed_kubernetes.CreateCluster.id
  policy_name = var.policy_name
  action      = "deny"
  namespaces = [
    "test"
  ]
  parameters = {
    hostNetwork = true
    min         = 20
    max         = 200
  }
}
`, name)
}

func TestAccAliCloudCSKubernetesPolicyInstance_collectionOrder(t *testing.T) {
	var v map[string]interface{}
	resourceID := "alicloud_cs_kubernetes_policy_instance.default"
	ra := acctest.ResourceAttrInit(resourceID, AliCloudCSKubernetesPolicyInstanceMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceID, &v, func() interface{} {
		return cs.NewAckServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeAckPolicyInstance")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceID, name, AliCloudCSKubernetesPolicyInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"action":     "deny",
					"namespaces": []string{"test", "test1", "test2"},
					"parameters": map[string]interface{}{
						"hostNetwork": "false",
						"min":         30,
						"max":         300,
					},
					"cluster_id":  "${alicloud_cs_managed_kubernetes.CreateCluster.id}",
					"policy_name": "${var.policy_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"action":       "deny",
						"namespaces.#": "3",
						"cluster_id":   acctest.CHECKSET,
						"policy_name":  acctest.CHECKSET,
					}),
				),
			},
			{
				// namespaces: verify reorder produces a non-empty plan
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{"test2", "test1", "test"},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// Apply the reordered config and verify convergence (empty plan)
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
		},
	})
}

// Test Ack PolicyInstance. <<< Resource test cases, automatically generated.
