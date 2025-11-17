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
// Case 策略管理测试-delayedTime 11657
func TestAccAliCloudAckPolicyInstance_basic11657(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_policy_instance.default"
	ra := resourceAttrInit(resourceId, AlicloudAckPolicyInstanceMap11657)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckPolicyInstance")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckPolicyInstanceBasicDependence11657)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_action": "deny",
					"namespaces":    []string{},
					"parameters": map[string]interface{}{
						"\"repos\"": "[\\\"registry-vpc.cn-hangzhou.aliyuncs.com\\\\/acs\\\\/\\\"]",
					},
					"cluster_id":  "${alicloud_cs_managed_kubernetes.创建Cluster.id}",
					"policy_name": "${var.policy_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_action": "deny",
						"namespaces.#":  "0",
						"cluster_id":    CHECKSET,
						"policy_name":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_action": "warn",
					"namespaces": []string{
						"${var.policy_scope}", "${var.policy_scope-kube-public}", "${var.policy_scope_update}"},
					"parameters": map[string]interface{}{
						"\"repos\"": "[\\\"registry-vpc.cn-hangzhou.aliyuncs.com\\\\/acs\\\\/\\\",\\\"registry.cn-hangzhou.aliyuncs.com\\\\/acs\\\\/\\\"]",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_action": "warn",
						"namespaces.#":  "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"namespaces": []string{
						"${var.policy_scope}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"namespaces.#": "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudAckPolicyInstanceMap11657 = map[string]string{
	"instance_name": CHECKSET,
}

func AlicloudAckPolicyInstanceBasicDependence11657(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

variable "policy-template-controller_version" {
  default = "v0.4.0.0-gddee19d-aliyun"
}

variable "zone_id" {
  default = "cn-shanghai-b"
}

variable "cluster_name" {
  default = "test-create-cluster"
}

variable "policy_name" {
  default = "ACKAllowedRepos"
}

variable "policy_scope" {
  default = "default"
}

variable "policy_scope-kube-public" {
  default = "kube-public"
}

variable "region_id" {
  default = "cn-shanghai"
}

variable "policy_scope_update" {
  default = "kube-system"
}

variable "cidr_block" {
  default = "172.18.0.0/21"
}

variable "gatekeeper_version" {
  default = "3.18.2-release"
}

variable "loongcollector_version" {
  default = "3.1.6"
}

resource "alicloud_vpc" "创建VPC" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "创建VSwitch" {
  vpc_id     = alicloud_vpc.创建VPC.id
  zone_id    = var.zone_id
  cidr_block = var.cidr_block
}

resource "alicloud_cs_managed_kubernetes" "创建Cluster" {
  addons {
    name     = "gatekeeper"
    disabled = false
  }
  addons {
    name     = "loongcollector"
    disabled = false
  }
  addons {
    name     = "policy-template-controller"
    disabled = false
  }
  addons {
    name     = "terway-eniip"
    config   = "{\"IPVlan\":\"false\",\"NetworkPolicy\":\"false\",\"ENITrunking\":\"true\"}"
    disabled = false
  }
  addons {
    name     = "terway-controlplane"
    config   = "{\"ENITrunking\":\"true\"}"
    disabled = false
  }
  addons {
    name     = "coredns"
    disabled = false
  }
  addons {
    name     = "metrics-server"
    disabled = false
  }
  addons {
    name     = "nginx-ingress-controller"
    disabled = false
  }
  addons {
    name     = "managed-csiprovisioner"
    disabled = false
  }
  addons {
    name     = "csi-plugin"
    disabled = false
  }
  addons {
    name     = "storage-operator"
    disabled = false
  }
  is_enterprise_security_group = true
  vswitch_ids                  = ["${alicloud_vswitch.创建VSwitch.id}"]
  service_cidr                 = var.service_cidr
  pod_vswitch_ids              = ["${alicloud_vswitch.创建VSwitch.id}"]
  ip_stack                     = "ipv4"
  proxy_mode                   = "ipvs"
  deletion_protection          = false
  operation_policy {
    cluster_auto_upgrade {
      enabled = false
    }
  }
  maintenance_window {
    enable = false
  }
  profile      = "Default"
  cluster_spec = "ack.pro.small"
}


`, name)
}

// Test Ack PolicyInstance. <<< Resource test cases, automatically generated.
