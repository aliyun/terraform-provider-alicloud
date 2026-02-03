package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 检查项扫描策略配置20251209005 12031
// lintignore: AT001
func TestAccAliCloudThreatDetectionCheckConfig_basic12031(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_check_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionCheckConfigMap12031)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionCheckConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionCheckConfigBasicDependence12031)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "18",
					"enable_auto_check": "true",
					"vendors": []string{
						"ALIYUN"},
					"cycle_days": []string{
						"7", "1", "2"},
					"enable_add_check": "true",
					"start_time":       "12",
					"configure":        "not",
					"system_config":    "false",
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   370,
							"section_id": 515,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time":          "18",
						"enable_auto_check": "true",
						"vendors.#":         "1",
						"cycle_days.#":      "3",
						"enable_add_check":  "true",
						"start_time":        "12",
						"configure":         "not",
						"system_config":     "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"end_time":          "12",
					"enable_auto_check": "false",
					"enable_add_check":  "false",
					"start_time":        "6",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"end_time":          "12",
						"enable_auto_check": "false",
						"enable_add_check":  "false",
						"start_time":        "6",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_auto_check": "true",
					"enable_add_check":  "true",
					"cycle_days": []string{
						"4"},
					"configure": "not",
					"selected_checks": []map[string]interface{}{
						{
							"check_id":   23,
							"section_id": 11,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_auto_check": "true",
						"cycle_days.#":      "1",
						"enable_add_check":  "true",
						"configure":         "not",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"configure", "system_config", "vendors"},
			},
		},
	})
}

var AlicloudThreatDetectionCheckConfigMap12031 = map[string]string{}

func AlicloudThreatDetectionCheckConfigBasicDependence12031(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}

// Test ThreatDetection CheckConfig. <<< Resource test cases, automatically generated.
