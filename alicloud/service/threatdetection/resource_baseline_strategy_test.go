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

// Case 1
func TestAccAliCloudThreatDetectionBaselineStrategy_basic1862(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_baseline_strategy.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionBaselineStrategyMap1862)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewSasService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionBaselineStrategy")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sThreatDetectionBaselineStrategy%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionBaselineStrategyBasicDependence1862)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.CheckResourceDestroy(),
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
						"baseline_strategy_name": acctest.CHECKSET,
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
}`, name)
}
