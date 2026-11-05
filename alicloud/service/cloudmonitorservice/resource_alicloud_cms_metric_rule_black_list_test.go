package cloudmonitorservice_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/cloudmonitorservice"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudCmsMetricRuleBlackList_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_black_list.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsMetricRuleBlackListMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsMetricRuleBlackList")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsMetricRuleBlackList%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsMetricRuleBlackListBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instances": []string{
						"{\\\"instanceId\\\":\\\"instance_id_to_test\\\"}"},
					"metrics": []map[string]interface{}{
						{
							"metric_name": "disk_utilization",
							"resource":    "host1",
						},
					},
					"category":       "ecs",
					"effective_time": "2022-01-01 00:00:00",
					"is_enable":      true,
					"scope_type":     "USER",
					"scope_value": []string{
						"scope-a",
					},
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
						"effective_time":              "2022-01-01 00:00:00",
						"scope_type":                  "USER",
						"scope_value.#":               "1",
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
					"category":                    "rds",
					"effective_time":              "2022-02-01 00:00:00",
					"is_enable":                   false,
					"instances": []string{
						"{\\\"instanceId\\\":\\\"instance_id_to_test_update\\\"}",
					},
					"metrics": []map[string]interface{}{
						{
							"metric_name": "cpu_utilization",
							"resource":    "host2",
						},
					},
					"scope_type": "GROUP",
					"scope_value": []string{
						"scope-b",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"metric_rule_black_list_name": "henghai1342432432432update",
						"enable_end_time":             "1640609200000",
						"enable_start_time":           "1640238400000",
						"category":                    "rds",
						"effective_time":              "2022-02-01 00:00:00",
						"instances.#":                 "1",
						"metrics.#":                   "1",
						"scope_type":                  "GROUP",
						"scope_value.#":               "1",
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

// TypeList order coverage for alicloud_cms_metric_rule_black_list.instances.
func TestAccAliCloudCmsMetricRuleBlackList_instancesOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_black_list.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsMetricRuleBlackListMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsMetricRuleBlackList")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsMetricRuleBlackListOrder%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsMetricRuleBlackListBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instances": []string{
						"{\\\"instanceId\\\":\\\"i-aaa\\\"}",
						"{\\\"instanceId\\\":\\\"i-bbb\\\"}",
					},
					"metrics": []map[string]interface{}{
						{
							"metric_name": "disk_utilization",
						},
					},
					"category":                    "ecs",
					"enable_end_time":             "1640608200000",
					"namespace":                   "acs_ecs_dashboard",
					"enable_start_time":           "1640237400000",
					"metric_rule_black_list_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instances.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instances": []string{
						"{\\\"instanceId\\\":\\\"i-bbb\\\"}",
						"{\\\"instanceId\\\":\\\"i-aaa\\\"}",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instances.#": "2",
					}),
				),
			},
		},
	})
}

// TypeList order coverage for alicloud_cms_metric_rule_black_list.scope_value.
func TestAccAliCloudCmsMetricRuleBlackList_scopeValueOrder(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_metric_rule_black_list.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudCmsMetricRuleBlackListMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return cloudmonitorservice.NewCmsService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeCmsMetricRuleBlackList")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sCmsMetricRuleBlackListScope%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudCmsMetricRuleBlackListBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instances": []string{
						"{\\\"instanceId\\\":\\\"i-aaa\\\"}",
					},
					"metrics": []map[string]interface{}{
						{
							"metric_name": "disk_utilization",
						},
					},
					"category":                    "ecs",
					"enable_end_time":             "1640608200000",
					"namespace":                   "acs_ecs_dashboard",
					"enable_start_time":           "1640237400000",
					"metric_rule_black_list_name": name,
					"scope_value": []string{
						"scope-a",
						"scope-b",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope_value.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"scope_value": []string{
						"scope-b",
						"scope-a",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scope_value.#": "2",
					}),
				),
			},
		},
	})
}
