// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Eflo ErRouteMap. >>> Resource test cases, automatically generated.
// Case er_route_map_vpd2vpd_permit 12479
func TestAccAliCloudEfloErRouteMap_basic12479(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_er_route_map.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloErRouteMapMap12479)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloErRouteMap")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloErRouteMapBasicDependence12479)
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
					"transmission_instance_type":  "VPD",
					"action":                      "permit",
					"reception_instance_type":     "VPD",
					"description":                 "permit-route-map-vpd-to-vpd",
					"reception_instance_id":       "${alicloud_eflo_vpd.reception_vpd.id}",
					"er_id":                       "${alicloud_eflo_er.ER.id}",
					"reception_instance_owner":    "${data.alicloud_account.default.id}",
					"transmission_instance_owner": "${data.alicloud_account.default.id}",
					"transmission_instance_id":    "${alicloud_eflo_vpd.transmission_vpd.id}",
					"er_route_map_num":            "1001",
					"destination_cidr_block":      "0.0.0.0/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transmission_instance_type":  "VPD",
						"action":                      "permit",
						"reception_instance_type":     "VPD",
						"description":                 "permit-route-map-vpd-to-vpd",
						"reception_instance_id":       CHECKSET,
						"er_id":                       CHECKSET,
						"reception_instance_owner":    CHECKSET,
						"transmission_instance_owner": CHECKSET,
						"transmission_instance_id":    CHECKSET,
						"er_route_map_num":            "1001",
						"destination_cidr_block":      "0.0.0.0/0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "permit-route-map-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "permit-route-map-updated",
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

var AlicloudEfloErRouteMapMap12479 = map[string]string{
	"status":          CHECKSET,
	"create_time":     CHECKSET,
	"er_route_map_id": CHECKSET,
	"region_id":       CHECKSET,
}

func AlicloudEfloErRouteMapBasicDependence12479(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "zone_id" {
  default = "cn-hangzhou-a"
}

data "alicloud_account" "default" {}

resource "alicloud_eflo_vpd" "transmission_vpd" {
  cidr     = "10.0.0.0/8"
  vpd_name = "tf-transmission-${var.name}"
}

resource "alicloud_eflo_vpd" "reception_vpd" {
  cidr     = "192.168.0.0/16"
  vpd_name = "tf-reception-${var.name}"
}

resource "alicloud_eflo_er" "ER" {
  er_name        = "tf-er-routemap-${var.name}"
  master_zone_id = var.zone_id
}

`, name)
}

// Case er_route_map_vpd2vpd_deny 12484
func TestAccAliCloudEfloErRouteMap_basic12484(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eflo_er_route_map.default"
	ra := resourceAttrInit(resourceId, AlicloudEfloErRouteMapMap12484)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EfloServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEfloErRouteMap")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacceflo%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudEfloErRouteMapBasicDependence12484)
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
					"transmission_instance_type":  "VPD",
					"action":                      "deny",
					"reception_instance_type":     "VPD",
					"description":                 "deny-route-map-vpd-to-vpd",
					"reception_instance_id":       "${alicloud_eflo_vpd.reception_vpd.id}",
					"er_id":                       "${alicloud_eflo_er.ER.id}",
					"reception_instance_owner":    "${data.alicloud_account.default.id}",
					"transmission_instance_owner": "${data.alicloud_account.default.id}",
					"transmission_instance_id":    "${alicloud_eflo_vpd.transmission_vpd.id}",
					"er_route_map_num":            "1002",
					"destination_cidr_block":      "0.0.0.0/0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"transmission_instance_type":  "VPD",
						"action":                      "deny",
						"reception_instance_type":     "VPD",
						"description":                 "deny-route-map-vpd-to-vpd",
						"reception_instance_id":       CHECKSET,
						"er_id":                       CHECKSET,
						"reception_instance_owner":    CHECKSET,
						"transmission_instance_owner": CHECKSET,
						"transmission_instance_id":    CHECKSET,
						"er_route_map_num":            "1002",
						"destination_cidr_block":      "0.0.0.0/0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "deny-route-map-updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "deny-route-map-updated",
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

var AlicloudEfloErRouteMapMap12484 = map[string]string{
	"status":          CHECKSET,
	"create_time":     CHECKSET,
	"er_route_map_id": CHECKSET,
	"region_id":       CHECKSET,
}

func AlicloudEfloErRouteMapBasicDependence12484(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "region_id" {
  default = "cn-hangzhou"
}

variable "zone_id" {
  default = "cn-hangzhou-a"
}

data "alicloud_account" "default" {}

resource "alicloud_eflo_vpd" "transmission_vpd" {
  cidr     = "10.0.0.0/8"
  vpd_name = "tf-transmission-${var.name}"
}

resource "alicloud_eflo_vpd" "reception_vpd" {
  cidr     = "192.168.0.0/16"
  vpd_name = "tf-reception-${var.name}"
}

resource "alicloud_eflo_er" "ER" {
  er_name        = "tf-er-routemap-${var.name}"
  master_zone_id = var.zone_id
}

`, name)
}

// Test Eflo ErRouteMap. <<< Resource test cases, automatically generated.
