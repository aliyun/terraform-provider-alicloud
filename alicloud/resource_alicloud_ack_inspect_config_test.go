// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ack InspectConfig. >>> Resource test cases, automatically generated.
// Case 集群巡检配置测试用例 11778
func TestAccAliCloudAckInspectConfig_basic11778(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ack_inspect_config.default"
	ra := resourceAttrInit(resourceId, AlicloudAckInspectConfigMap11778)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &AckServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeAckInspectConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccack%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudAckInspectConfigBasicDependence11778)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"recurrence": "FREQ=DAILY;BYHOUR=10;BYMINUTE=15",
					"disabled_check_items": []string{
						"APIServerCLBListenerAbnormal"},
					"enabled":           "true",
					"inspect_config_id": "${alicloud_cs_managed_kubernetes.创建Cluster.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recurrence":             "FREQ=DAILY;BYHOUR=10;BYMINUTE=15",
						"disabled_check_items.#": "1",
						"enabled":                "true",
						"inspect_config_id":      CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"recurrence": "FREQ=DAILY;BYHOUR=10;BYMINUTE=16",
					"disabled_check_items": []string{
						"APIServerCLBListenerAbnormal", "NodeLocalCacheButNoInjection", "APIServerCLBBackendAbnormal"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"recurrence":             "FREQ=DAILY;BYHOUR=10;BYMINUTE=16",
						"disabled_check_items.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled_check_items": []string{
						"APIServerCLBBackendAbnormal", "NodeLocalCacheButNoInjection"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled_check_items.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled_check_items": []string{
						"NodeLocalCacheButNoInjection"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled_check_items.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"disabled_check_items": []string{},
					"enabled":              "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"disabled_check_items.#": "0",
						"enabled":                "false",
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

var AlicloudAckInspectConfigMap11778 = map[string]string{}

func AlicloudAckInspectConfigBasicDependence11778(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "service_cidr" {
  default = "192.168.0.0/16"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_cs_managed_kubernetes" "创建Cluster" {
  addons {
    name   = "terway-eniip"
    config = "{\"IPVlan\":\"false\",\"NetworkPolicy\":\"false\",\"ENITrunking\":\"true\"}"
  }
  addons {
    name   = "terway-controlplane"
    config = "{\"ENITrunking\":\"true\"}"
  }
  addons {
    name = "csi-plugin"
  }
  addons {
    name = "managed-csiprovisioner"
  }
  addons {
    name = "nginx-ingress-controller"
  }
  addons {
    name = "metrics-server"
  }
  addons {
    name = "coredns"
  }
  ip_stack                     = "ipv4"
  is_enterprise_security_group = true
  service_cidr                 = var.service_cidr
  proxy_mode                   = "ipvs"
  deletion_protection          = false
  operation_policy {
    cluster_auto_upgrade {
      enabled = false
    }
  }
  maintenance_window {
    enable           = true
    maintenance_time = "2025-11-03T00:00:00.000+08:00"
    duration         = "3h"
    weekly_period    = "Monday"
  }
  zone_ids = [data.alicloud_zones.default.zones.0.id]
}


`, name)
}

// Test Ack InspectConfig. <<< Resource test cases, automatically generated.
