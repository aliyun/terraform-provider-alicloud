// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ack PolicyInstance. >>> Resource test cases, automatically generated.
func TestAccAliCloudCSKubernetesPolicyInstance_basic(t *testing.T) {
	var v map[string]interface{}
	resourceID := "alicloud_cs_kubernetes_policy_instance.default"
	ra := resourceAttrInit(resourceID, AliCloudCSKubernetesPolicyInstanceMap)
	rc := resourceCheckInitWithDescribeMethod(resourceID, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckPolicyInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceID, name, AliCloudCSKubernetesPolicyInstanceBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceID,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
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
						"action":       "deny",
						"namespaces.#": "0",
						"cluster_id":   CHECKSET,
						"policy_name":  CHECKSET,
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
						"action":       "warn",
						"namespaces.#": "3",
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
				ResourceName:            resourceID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters"},
			},
		},
	})
}

var AliCloudCSKubernetesPolicyInstanceMap = map[string]string{
	"instance_name": CHECKSET,
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

variable "cluster_name" {
  default = "example-create-cluster"
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
  cluster_spec                 = "ack.pro.small"
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

// Test Ack PolicyInstance. <<< Resource test cases, automatically generated.
