package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudCmsMetricRuleBlackList_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_black_list.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsMetricRuleBlackListMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMetricRuleBlackList")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsMetricRuleBlackList%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsMetricRuleBlackListBasicDependence)
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
					"instances": []string{
						"{\\\"instanceId\\\":\\\"instance_id_to_test\\\"}"},
					"metrics": []map[string]interface{}{
						{
							"metric_name": "disk_utilization",
						},
					},
					"category":                    "ecs",
					"enable_end_time":             "1640608200000",
					"namespace":                   "acs_ecs_dashboard",
					"enable_start_time":           "1640237400000",
					"metric_rule_black_list_name": "henghai1342432432432",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instances.#":                 "1",
						"metrics.#":                   "1",
						"category":                    "ecs",
						"enable_end_time":             "1640608200000",
						"namespace":                   "acs_ecs_dashboard",
						"enable_start_time":           "1640237400000",
						"metric_rule_black_list_name": "henghai1342432432432",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"metric_rule_black_list_name": "henghai1342432432432update",
					"enable_end_time":             "1640609200000",
					"enable_start_time":           "1640238400000",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_rule_black_list_name": "henghai1342432432432update",
						"enable_end_time":             "1640609200000",
						"enable_start_time":           "1640238400000",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			}},
	})
}

var AlicloudCmsMetricRuleBlackListMap = map[string]string{}

func AlicloudCmsMetricRuleBlackListBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}
