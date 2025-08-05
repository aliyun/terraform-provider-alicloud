// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

// Test ThreatDetection AttackPathSensitiveAssetConfig. >>> Resource test cases, automatically generated.
// Case AttackPathSensitiveAssetConfig_250219_03 10276
func TestAccAliCloudThreatDetectionAttackPathSensitiveAssetConfig_basic10276(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_attack_path_sensitive_asset_config.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAttackPathSensitiveAssetConfigMap10276)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAttackPathSensitiveAssetConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAttackPathSensitiveAssetConfigBasicDependence10276)
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
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "${data.alicloud_ddoscoo_instances.default.instances.0.id}",
							"vendor":         "0",
							"asset_type":     "16",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attack_path_asset_list.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "${data.alicloud_slb_load_balancers.default.balancers.0.id}",
							"vendor":         "0",
							"asset_type":     "1",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
						{
							"instance_id":    "${data.alicloud_slb_load_balancers.default.balancers.1.id}",
							"vendor":         "0",
							"asset_type":     "1",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
						{
							"instance_id":    "${data.alicloud_ssl_certificates_service_certificates.default.certificates.0.id}",
							"vendor":         "0",
							"asset_type":     "13",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attack_path_asset_list.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "${data.alicloud_ddoscoo_instances.default.instances.0.id}",
							"vendor":         "0",
							"asset_type":     "16",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
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

func TestAccAliCloudThreatDetectionAttackPathSensitiveAssetConfig_basic10276_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_attack_path_sensitive_asset_config.default"
	ra := resourceAttrInit(resourceId, AliCloudThreatDetectionAttackPathSensitiveAssetConfigMap10276)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAttackPathSensitiveAssetConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudThreatDetectionAttackPathSensitiveAssetConfigBasicDependence10276)
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
					"attack_path_asset_list": []map[string]interface{}{
						{
							"instance_id":    "${data.alicloud_slb_load_balancers.default.balancers.0.id}",
							"vendor":         "0",
							"asset_type":     "1",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
						{
							"instance_id":    "${data.alicloud_slb_load_balancers.default.balancers.1.id}",
							"vendor":         "0",
							"asset_type":     "1",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
						{
							"instance_id":    "${data.alicloud_ssl_certificates_service_certificates.default.certificates.0.id}",
							"vendor":         "0",
							"asset_type":     "13",
							"asset_sub_type": "0",
							"region_id":      "${data.alicloud_regions.default.regions.0.id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"attack_path_asset_list.#": "3",
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

var AliCloudThreatDetectionAttackPathSensitiveAssetConfigMap10276 = map[string]string{}

func AliCloudThreatDetectionAttackPathSensitiveAssetConfigBasicDependence10276(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_regions" "default" {
  current = true
}

data "alicloud_ddoscoo_instances" "default" {
}

data "alicloud_slb_load_balancers" "default" {
}

data "alicloud_ssl_certificates_service_certificates" "default" {
}

`, name)
}

// Test ThreatDetection AttackPathSensitiveAssetConfig. <<< Resource test cases, automatically generated.
