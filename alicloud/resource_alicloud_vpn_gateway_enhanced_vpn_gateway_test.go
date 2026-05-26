// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test VpnGateway EnhancedVpnGateway. >>> Resource test cases, automatically generated.
// Case EnhancedVpnGateway_副本1776853707919 12764
func TestAccAliCloudVpnGatewayEnhancedVpnGateway_basic12764(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_vpn_gateway_enhanced_vpn_gateway.default"
	ra := resourceAttrInit(resourceId, AlicloudVpnGatewayEnhancedVpnGatewayMap12764)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VPNGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeVpnGatewayEnhancedVpnGateway")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccvpngateway%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudVpnGatewayEnhancedVpnGatewayBasicDependence12764)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"ap-southeast-3"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_type":                     "Normal",
					"description":                  "default",
					"disaster_recovery_vswitch_id": "${alicloud_vswitch.default23kGFr.id}",
					"vpc_id":                       "${alicloud_vpc.defaulttYTx5F.id}",
					"vpn_gateway_name":             "default",
					"network_type":                 "public",
					"vswitch_id":                   "${alicloud_vswitch.defaultTRk7k3.id}",
					"gateway_type":                 "Enhanced.SiteToSite",
					"auto_propagate":               "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_type":                     "Normal",
						"description":                  "default",
						"disaster_recovery_vswitch_id": CHECKSET,
						"vpc_id":                       CHECKSET,
						"vpn_gateway_name":             "default",
						"network_type":                 "public",
						"vswitch_id":                   CHECKSET,
						"gateway_type":                 "Enhanced.SiteToSite",
						"auto_propagate":               "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description":      "new",
					"vpn_gateway_name": "new",
					"auto_propagate":   "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":      "new",
						"vpn_gateway_name": "new",
						"auto_propagate":   "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
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

var AlicloudVpnGatewayEnhancedVpnGatewayMap12764 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudVpnGatewayEnhancedVpnGatewayBasicDependence12764(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region" {
  default = "ap-southeast-3"
}

variable "zone2" {
  default = "ap-southeast-3a"
}

variable "zone1" {
  default = "ap-southeast-3b"
}

resource "alicloud_vpc" "defaulttYTx5F" {
  cidr_block = "192.168.0.0/16"
  is_default = false
}

resource "alicloud_vswitch" "defaultTRk7k3" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone1
  cidr_block = "192.168.10.0/24"
}

resource "alicloud_vswitch" "default23kGFr" {
  vpc_id     = alicloud_vpc.defaulttYTx5F.id
  zone_id    = var.zone2
  cidr_block = "192.168.20.0/24"
}


`, name)
}

// Test VpnGateway EnhancedVpnGateway. <<< Resource test cases, automatically generated.
