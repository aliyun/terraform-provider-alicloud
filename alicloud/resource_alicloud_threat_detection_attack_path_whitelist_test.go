// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection AttackPathWhitelist. >>> Resource test cases, automatically generated.
func TestAccAliCloudThreatDetectionAttackPathWhitelist_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_attack_path_whitelist.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAttackPathWhitelistMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAttackPathWhitelist")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAttackPathWhitelistBasicDependence)
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
					"path_type":      "role_escalation",
					"whitelist_type": "ALL_ASSET",
					"whitelist_name": name,
					"path_name":      "ecs_get_credential_by_create_login_profile",
					"remark":         name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path_type":                "role_escalation",
						"whitelist_type":           "ALL_ASSET",
						"whitelist_name":           name,
						"path_name":                "ecs_get_credential_by_create_login_profile",
						"remark":                   name,
						"attack_path_asset_list.#": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"path_type":      "role_escalation",
					"whitelist_type": "PART_ASSET",
					"whitelist_name": name + "_update",
					"path_name":      "ecs_get_credential_by_create_login_profile",
					"remark":         name + "_update",
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "AliyunYundunSASReadOnlyAccess::System",
							"vendor":         "0",
							"asset_type":     "15",
							"asset_sub_type": "2",
							"region_id":      "cn-hangzhou",
							"node_type":      "end",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"whitelist_type":           "PART_ASSET",
						"whitelist_name":           name + "_update",
						"remark":                   name + "_update",
						"attack_path_asset_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"path_type":      "role_escalation",
					"whitelist_type": "PART_ASSET",
					"whitelist_name": name + "_update",
					"path_name":      "ecs_get_credential_by_create_login_profile",
					"remark":         name + "_update",
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "AliyunYundunSASReadOnlyAccess::System",
							"vendor":         "0",
							"asset_type":     "15",
							"asset_sub_type": "2",
							"region_id":      "cn-hangzhou",
							"node_type":      "start",
						},
						{
							"instance_id":    "${data.alicloud_ddoscoo_instances.default.instances.0.id}",
							"vendor":         "0",
							"asset_type":     "16",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
							"node_type":      "end",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attack_path_asset_list.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"path_type":      "sensitive_asset",
					"whitelist_type": "PART_ASSET",
					"whitelist_name": name + "_update",
					"path_name":      "ram_user_has_access_to_sensitive_asset",
					"remark":         name + "_update",
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "AliyunYundunSASReadOnlyAccess::System",
							"vendor":         "0",
							"asset_type":     "15",
							"asset_sub_type": "2",
							"region_id":      "cn-hangzhou",
							"node_type":      "start",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path_type":                "sensitive_asset",
						"path_name":                "ram_user_has_access_to_sensitive_asset",
						"attack_path_asset_list.#": "1",
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

func TestAccAliCloudThreatDetectionAttackPathWhitelist_basic_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_attack_path_whitelist.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAttackPathWhitelistMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAttackPathWhitelist")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAttackPathWhitelistBasicDependence)
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
					"path_type":      "role_escalation",
					"whitelist_type": "PART_ASSET",
					"whitelist_name": name,
					"path_name":      "ecs_get_credential_by_create_login_profile",
					"remark":         name,
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "AliyunYundunSASReadOnlyAccess::System",
							"vendor":         "0",
							"asset_type":     "15",
							"asset_sub_type": "2",
							"region_id":      "cn-hangzhou",
							"node_type":      "end",
						},
						{
							"instance_id":    "AliyunYundunSASReadOnlyAccess::System",
							"vendor":         "0",
							"asset_type":     "15",
							"asset_sub_type": "2",
							"region_id":      "cn-hangzhou",
							"node_type":      "start",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path_type":                "role_escalation",
						"whitelist_type":           "PART_ASSET",
						"whitelist_name":           name,
						"path_name":                "ecs_get_credential_by_create_login_profile",
						"remark":                   name,
						"attack_path_asset_list.#": "2",
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

func TestAccAliCloudThreatDetectionAttackPathWhitelist_allasset(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_attack_path_whitelist.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAttackPathWhitelistMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAttackPathWhitelist")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAttackPathWhitelistBasicDependence)
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
					"path_type":      "role_escalation",
					"whitelist_type": "ALL_ASSET",
					"whitelist_name": name,
					"path_name":      "ecs_get_credential_by_create_login_profile",
					"remark":         name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"path_type":                "role_escalation",
						"whitelist_type":           "ALL_ASSET",
						"whitelist_name":           name,
						"path_name":                "ecs_get_credential_by_create_login_profile",
						"remark":                   name,
						"attack_path_asset_list.#": "0",
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

var AliCloudThreatDetectionAttackPathWhitelistMap = map[string]string{}

func AliCloudThreatDetectionAttackPathWhitelistBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_ddoscoo_instances" "default" {
}

`, name)
}

// Test ThreatDetection AttackPathWhitelist. <<< Resource test cases, automatically generated.
