package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection AssetSelectionConfig. >>> Resource test cases, automatically generated.
// Case AssetSelectionConfig 9090
func TestAccAliCloudThreatDetectionAssetSelectionConfig_basic9090(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_asset_selection_config.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionAssetSelectionConfigMap9090)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionAssetSelectionConfig")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionAssetSelectionConfigBasicDependence9090)
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
					"business_type": "agentlesss_vul_white_1",
					"target_type":   "instance",
					"platform":      "all",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"business_type": "agentlesss_vul_white_1",
						"target_type":   "instance",
						"platform":      "all",
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

var AlicloudThreatDetectionAssetSelectionConfigMap9090 = map[string]string{}

func AlicloudThreatDetectionAssetSelectionConfigBasicDependence9090(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection AssetSelectionConfig. <<< Resource test cases, automatically generated.
