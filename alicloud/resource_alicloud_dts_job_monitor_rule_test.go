package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// this resource depends on dts job which tf not support yet, so cannot test by automation.
func SkipTestAccAlicloudDTSJobMonitorRule_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_dts_job_monitor_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDTSJobMonitorRuleMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDtsJobMonitorRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sdtsjobmonitorrule%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDTSJobMonitorRuleBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"dts_job_id": "pf4q95n1v3gebwz",
					"type":       "delay",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dts_job_id": "pf4q95n1v3gebwz",
						"type":       "delay",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"state": "Y",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"state": "Y",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"phone": "12345678987",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"phone": "12345678987",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"delay_rule_time": "233",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"delay_rule_time": "233",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"type":            "delay",
					"state":           "N",
					"phone":           "12345678989",
					"delay_rule_time": "234",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"type":            "delay",
						"state":           "N",
						"phone":           "12345678989",
						"delay_rule_time": "234",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudDTSJobMonitorRuleMap0 = map[string]string{
	"dts_job_id": CHECKSET,
	"state":      CHECKSET,
	"phone":      CHECKSET,
}

func AlicloudDTSJobMonitorRuleBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
`, name)
}
