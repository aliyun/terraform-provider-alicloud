// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package threatdetection_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/threatdetection"
	tFacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test ThreatDetection AttackPathSensitiveAssetConfig. >>> Resource test cases, automatically generated.
// Case AttackPathSensitiveAssetConfig_250219_03 10276
func TestAccAliCloudThreatDetectionAttackPathSensitiveAssetConfig_basic10276(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_attack_path_sensitive_asset_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AliCloudThreatDetectionAttackPathSensitiveAssetConfigMap10276)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionAttackPathSensitiveAssetConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudThreatDetectionAttackPathSensitiveAssetConfigBasicDependence10276)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
	ra := acctest.ResourceAttrInit(resourceId, AliCloudThreatDetectionAttackPathSensitiveAssetConfigMap10276)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionAttackPathSensitiveAssetConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudThreatDetectionAttackPathSensitiveAssetConfigBasicDependence10276)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
