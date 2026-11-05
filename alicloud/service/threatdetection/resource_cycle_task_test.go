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

// NOTE: Test depends on data source or hardcoded are not stable and may fail at any time

// Test ThreatDetection CycleTask. >>> Resource test cases, automatically generated.
// Case CycleTask 10574
func TestAccAliCloudThreatDetectionCycleTask_basic10574(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_threat_detection_cycle_task.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudThreatDetectionCycleTaskMap10574)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return threatdetection.NewServiceV2(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeThreatDetectionCycleTask")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tFacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccthreatdetection%d", rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudThreatDetectionCycleTaskBasicDependence10574)
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
					"target_end_time":   "6",
					"task_name":         "VIRUS_VUL_SCHEDULE_SCAN",
					"task_type":         "VIRUS_VUL_SCHEDULE_SCAN",
					"param":             "{\\\"targetInfo\\\":[{\\\"name\\\":\\\"未分组\\\",\\\"type\\\":\\\"groupId\\\",\\\"target\\\":10358625},{\\\"name\\\":\\\"test\\\",\\\"type\\\":\\\"groupId\\\",\\\"target\\\":12206415}]}",
					"first_date_str":    "1650556800000",
					"interval_period":   "7",
					"enable":            "1",
					"target_start_time": "0",
					"source":            "console_batch",
					"period_unit":       "day",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_end_time":   acctest.CHECKSET,
						"task_name":         "VIRUS_VUL_SCHEDULE_SCAN",
						"task_type":         "VIRUS_VUL_SCHEDULE_SCAN",
						"param":             acctest.CHECKSET,
						"first_date_str":    acctest.CHECKSET,
						"interval_period":   acctest.CHECKSET,
						"enable":            acctest.CHECKSET,
						"target_start_time": acctest.CHECKSET,
						"source":            "console_batch",
						"period_unit":       "day",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"target_end_time":   "12",
					"param":             "{\\\"targetInfo\\\": [{\\\"name\\\":\\\"未分组\\\",\\\"type\\\":\\\"groupId\\\",\\\"target\\\":10358625}]}",
					"first_date_str":    "1664380800000",
					"interval_period":   "14",
					"enable":            "0",
					"target_start_time": "6",
					"period_unit":       "hour",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"target_end_time":   "12",
						"param":             acctest.CHECKSET,
						"first_date_str":    acctest.CHECKSET,
						"interval_period":   acctest.CHECKSET,
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
}`, name)
}

// Test ThreatDetection CycleTask. <<< Resource test cases, automatically generated.
