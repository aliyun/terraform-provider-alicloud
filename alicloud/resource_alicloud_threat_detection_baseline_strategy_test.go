package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case 1
func TestAccAlicloudThreatDetectionBaselineStrategy_basic1862(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_baseline_strategy.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionBaselineStrategyMap1862)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionBaselineStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionBaselineStrategy%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionBaselineStrategyBasicDependence1862)
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
					"custom_type":            "custom",
					"end_time":               "08:00:00",
					"baseline_strategy_name": "${var.name}",
					"cycle_days":             "3",
					"target_type":            "groupId",
					"start_time":             "05:00:00",
					"risk_sub_type_name":     "hc_exploit_redis",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_type":            "custom",
						"end_time":               "08:00:00",
						"baseline_strategy_name": CHECKSET,
						"cycle_days":             "3",
						"target_type":            "groupId",
						"start_time":             "05:00:00",
					}),
				),
			}, {
				Config: testAccConfig(map[string]interface{}{
					"custom_type":        "custom",
					"end_time":           "10:00:00",
					"cycle_days":         "1",
					"target_type":        "groupId",
					"start_time":         "06:00:00",
					"risk_sub_type_name": "hc_exploit_redis",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_type": "custom",
						"end_time":    "10:00:00",
						"cycle_days":  "1",
						"target_type": "groupId",
						"start_time":  "06:00:00",
					}),
				),
			}, {
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"risk_sub_type_name"},
			},
		},
	})
}

var AlicloudThreatDetectionBaselineStrategyMap1862 = map[string]string{}

func AlicloudThreatDetectionBaselineStrategyBasicDependence1862(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}
