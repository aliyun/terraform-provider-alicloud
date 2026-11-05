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

// Test ThreatDetection AssetSelectionConfig. >>> Resource test cases, automatically generated.
// Case AssetSelectionConfig 9090
// lintignore: AT001
func TestAccAliCloudThreatDetectionAssetSelectionConfig_basic9090(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_asset_selection_config.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionAssetSelectionConfigMap9090)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionAssetSelectionConfig")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionAssetSelectionConfigBasicDependence9090)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
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
