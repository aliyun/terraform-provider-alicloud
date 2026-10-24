// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Apig Plugin. >>> Resource test cases, automatically generated.
// Case plugin_basic_test 12943
func TestAccAliCloudApigPlugin_basic12943(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPluginMap12943)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccapig%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPluginBasicDependence12943)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"plugin_class_id": "pls-crpqb35lhtgo800k2m86",
					"gateway_id":      "${alicloud_apig_gateway.plugin_gateway_pre.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_class_id": "pls-crpqb35lhtgo800k2m86",
						"gateway_id":      CHECKSET,
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

var AlicloudApigPluginMap12943 = map[string]string{
	"plugin_class_name": CHECKSET,
	"gateway_name":      CHECKSET,
}

func AlicloudApigPluginBasicDependence12943(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_apig_gateway" "plugin_gateway_pre" {
  network_access_config {
    type = "Internet"
  }
  vswitch {
    vswitch_id = data.alicloud_vswitches.default.ids.0
  }
  zone_config {
    select_option = "Auto"
  }
  vpc {
    vpc_id = data.alicloud_vpcs.default.ids.0
  }
  gateway_type = "API"
  payment_type = "PayAsYouGo"
  gateway_name = var.name
  spec         = "apigw.small.x1"
  log_config {
    sls {
      enable = true
    }
  }
}


`, name)
}

// Test Apig Plugin. <<< Resource test cases, automatically generated.

// TestAccAliCloudApigPlugin_aiGatewayReadBack verifies that an AI gateway plugin
// can be read back after Create. ListPlugins (GET /v1/plugins) does not return AI
// gateway plugins unless gatewayType is specified, so before the DescribeApigPlugin
// fix the Read following Create reported Plugin NotFound and apply failed. This
// case also covers plan/refresh drift detection and ImportStatePassthrough for AI
// gateway plugins.
func TestAccAliCloudApigPlugin_aiGatewayReadBack(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_apig_plugin.default"
	ra := resourceAttrInit(resourceId, AlicloudApigPluginMap12943)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ApigServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeApigPlugin")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc-apig-plugin-ai-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudApigPluginAIGatewayDependence)
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
					"plugin_class_id": "pls-crpqb35lhtgo800k2m8g",
					"gateway_id":      "${alicloud_apig_gateway.ai_gateway.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"plugin_class_id": "pls-crpqb35lhtgo800k2m8g",
						"gateway_id":      CHECKSET,
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

func AlicloudApigPluginAIGatewayDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "k" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "cn-hangzhou-k"
  cidr_block = "172.16.1.0/24"
}

resource "alicloud_vswitch" "j" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = "cn-hangzhou-j"
  cidr_block = "172.16.2.0/24"
}

resource "alicloud_apig_gateway" "ai_gateway" {
  gateway_name = var.name
  gateway_type = "AI"
  payment_type = "PayAsYouGo"
  spec         = "aigw.small.x1"
  vpc {
    vpc_id = alicloud_vpc.default.id
  }
  vswitch {
    vswitch_id = alicloud_vswitch.j.id
  }
  zone_config {
    select_option = "Manual"
  }
  zones {
    zone_id    = alicloud_vswitch.j.zone_id
    vswitch_id = alicloud_vswitch.j.id
  }
  zones {
    zone_id    = alicloud_vswitch.k.zone_id
    vswitch_id = alicloud_vswitch.k.id
  }
  network_access_config {
    type = "Intranet"
  }
  log_config {
    sls {
      enable = false
    }
  }
}

`, name)
}
