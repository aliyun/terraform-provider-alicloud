package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ddoscoo"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudDdoscooSchedulerRule_basic(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	var v ddoscoo.SchedulerRule
	resourceId := "alicloud_ddoscoo_scheduler_rule.default"
	ra := resourceAttrInit(resourceId, DdosCooSchedulerRuleMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DdoscooService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDdoscooSchedulerRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccDdoscooSchedulerRule%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, DdosCooSchedulerRuleBasicdependence)
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
					"rule_name": name,
					"rule_type": "6",
					"rules": []map[string]string{
						{
							"priority":   "50",
							"region_id":  "cn-hangzhou",
							"type":       "A",
							"value_type": "6",
							"value":      "1.1.1.1",
						},
						{
							"priority":   "100",
							"region_id":  "cn-hangzhou",
							"type":       "A",
							"value":      "203.107.54.136",
							"value_type": "1",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_name": name,
						"rule_type": "6",
						"rules.#":   "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"param", "resource_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resource_group_id": "${data.alicloud_resource_manager_resource_groups.defaults.groups[0].id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resource_group_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"rule_type": "6",
					"rules": []map[string]string{
						{
							"priority":   "100",
							"region_id":  "cn-hangzhou",
							"type":       "A",
							"value":      "170.33.2.125",
							"value_type": "1",
						},
						{
							"priority":   "50",
							"region_id":  "cn-hangzhou",
							"type":       "A",
							"value":      "170.33.14.193",
							"value_type": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"rule_type": "6",
						"rules.#":   "2",
					}),
				),
			},
		},
	})
}

var DdosCooSchedulerRuleMap = map[string]string{}

func DdosCooSchedulerRuleBasicdependence(name string) string {
	return `
data "alicloud_resource_manager_resource_groups" "defaults" {
  name_regex = "default"
}
`
}
