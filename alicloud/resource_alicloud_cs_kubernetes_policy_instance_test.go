// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			testAccPreCheckWithRegions(t, true, connectivity.ManagedKubernetesSupportedRegions)
			testAccPreCheck(t)
		},
		IDRefreshName:     resourceID,
		ProviderFactories: testAccProviderFactory,
		CheckDestroy:      rac.checkResourceDestroy(),
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
						"cluster_id":             CHECKSET,
						"policy_name":            CHECKSET,
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

func TestUnitCSKubernetesPolicyInstanceParametersForState(t *testing.T) {
	input := map[string]interface{}{
		"string": "value",
		"bool":   false,
		"int":    30,
		"object": map[string]interface{}{
			"z": 3,
			"a": true,
		},
		"list": []interface{}{3, true, "value"},
	}
	expected := map[string]interface{}{
		"string": "value",
		"bool":   "false",
		"int":    "30",
		"object": `{"a":true,"z":3}`,
		"list":   `[3,true,"value"]`,
	}

	actual, err := policyParametersForState(input)
	if err != nil {
		t.Fatalf("policyParametersForState returned an error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("policyParametersForState() = %#v, want %#v", actual, expected)
	}

	resourceData := schema.TestResourceDataRaw(t, resourceAliCloudCSKubernetesPolicyInstance().Schema, nil)
	if err := resourceData.Set("parameters", actual); err != nil {
		t.Fatalf("setting normalized policy parameters in state returned an error: %v", err)
	}

	roundTrip := NormalizeMap(actual)
	roundTripJSON, err := json.Marshal(roundTrip)
	if err != nil {
		t.Fatalf("marshalling round-trip policy parameters returned an error: %v", err)
	}
	inputJSON, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("marshalling input policy parameters returned an error: %v", err)
	}
	if string(roundTripJSON) != string(inputJSON) {
		t.Fatalf("NormalizeMap(policyParametersForState()) = %s, want %s", roundTripJSON, inputJSON)
	}
}

func TestUnitCSKubernetesPolicyInstanceParametersForStateUnsupportedType(t *testing.T) {
	_, err := policyParametersForState(map[string]interface{}{"invalid": make(chan int)})
	if err == nil {
		t.Fatal("policyParametersForState should reject unsupported value types")
	}
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

// Test Ack PolicyInstance. <<< Resource test cases, automatically generated.
