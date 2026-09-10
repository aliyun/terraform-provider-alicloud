// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test ThreatDetection CycleTask. >>> Resource test cases, automatically generated.
// Case CycleTask 10574
func TestAccAliCloudThreatDetectionCycleTask_basic10574(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_cycle_task.default"
	ra := resourceAttrInit(resourceId, AlicloudThreatDetectionCycleTaskMap10574)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ThreatDetectionServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeThreatDetectionCycleTask")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudThreatDetectionCycleTaskBasicDependence10574)
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
					"target_end_time":   "6",
					"task_name":         "VIRUS_VUL_SCHEDULE_SCAN",
					"task_type":         "VIRUS_VUL_SCHEDULE_SCAN",
					"param":             "{\\n      \\\"targetInfo\\\": [\\n            {\\n                  \\\"type\\\": \\\"groupId\\\",\\n                  \\\"name\\\": \\\"TI HOST\\\",\\n                  \\\"target\\\": 10597***\\n            },\\n            {\\n                  \\\"type\\\": \\\"groupId\\\",\\n                  \\\"name\\\": \\\"expense HOST\\\",\\n                  \\\"target\\\": 10597***\\n            }\\n      ]\\n}",
					"first_date_str":    "1650556800000",
					"interval_period":   "7",
					"enable":            "1",
					"target_start_time": "0",
					"source":            "console_batch",
					"period_unit":       "day",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_end_time":   CHECKSET,
						"task_name":         "VIRUS_VUL_SCHEDULE_SCAN",
						"task_type":         "VIRUS_VUL_SCHEDULE_SCAN",
						"param":             CHECKSET,
						"first_date_str":    CHECKSET,
						"interval_period":   CHECKSET,
						"enable":            CHECKSET,
						"target_start_time": CHECKSET,
						"source":            "console_batch",
						"period_unit":       "day",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_end_time":   "12",
					"param":             "{\\n      \\\"targetInfo\\\": [\\n            {\\n                  \\\"type\\\": \\\"groupId\\\",\\n                  \\\"name\\\": \\\"TI HOST\\\",\\n                  \\\"target\\\": 10597***\\n            }\\n      ]\\n}",
					"first_date_str":    "1664380800000",
					"interval_period":   "14",
					"enable":            "0",
					"target_start_time": "6",
					"period_unit":       "hour",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_end_time":   "12",
						"param":             CHECKSET,
						"first_date_str":    CHECKSET,
						"interval_period":   CHECKSET,
						"enable":            "0",
						"target_start_time": "6",
						"period_unit":       "hour",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"source"},
			},
		},
	})
}

var AlicloudThreatDetectionCycleTaskMap10574 = map[string]string{}

func AlicloudThreatDetectionCycleTaskBasicDependence10574(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test ThreatDetection CycleTask. <<< Resource test cases, automatically generated.
