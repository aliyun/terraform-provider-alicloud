package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudQuotasQuotaAlarm_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_quotas_quota_alarm.default"
	ra := resourceAttrInit(resourceId, AlicloudQuotasQuotaAlarmMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &QuotasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeQuotasQuotaAlarm")
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccQuotasQuotaAlarmTest%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudQuotasQuotaAlarmBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name":  name,
					"product_code":      "ecs",
					"quota_action_code": "q_prepaid-instance-count-per-once-purchase",
					"threshold":         "100",
					"quota_dimensions": []map[string]interface{}{
						{
							"key":   "regionId",
							"value": "cn-hangzhou",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name":   name,
						"product_code":       "ecs",
						"quota_action_code":  "q_prepaid-instance-count-per-once-purchase",
						"threshold":          "100",
						"quota_dimensions.#": "1",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name": name + "Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name": name + "Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"threshold": "200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"threshold": "200",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"quota_alarm_name": name,
					"threshold":        "100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"quota_alarm_name": name,
						"threshold":        "100",
					}),
				),
			},
		},
	})
}

var AlicloudQuotasQuotaAlarmMap = map[string]string{}

func AlicloudQuotasQuotaAlarmBasicDependence(name string) string {
	return ""
}
