package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms NotifyStrategy. >>> Resource test cases.

func TestAccAliCloudCmsNotifyStrategy_basic0(t *testing.T) {
	supportedRegions := []connectivity.Region{connectivity.Hangzhou}
	var v map[string]interface{}
	resourceId := "alicloud_cms_notify_strategy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsNotifyStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsNotifyStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := testAccCmsUniqueResourceName()
	notifyTarget := testAccCmsCreateContact(t, rand, supportedRegions)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return AliCloudCmsNotifyStrategyBasicDependence0(name, notifyTarget)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, supportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_strategy_name":         name,
					"workspace":                    "${var.cms_workspace}",
					"description":                  "tf-acc cms notify strategy test",
					"ignore_restored_notification": "false",
					"grouping_setting": []map[string]interface{}{
						{
							"grouping_keys": []string{"severity"},
							"period_min":    "5",
							"silence_sec":   "300",
							"times":         "1",
						},
					},
					"routes": []map[string]interface{}{
						{
							"channels": []map[string]interface{}{
								{
									"channel_type":         "CONTACT",
									"receivers":            []string{"${var.cms_contact_id}"},
									"enabled_sub_channels": []string{"EMAIL"},
								},
							},
							"effect_time_range": []map[string]interface{}{
								{
									"day_in_week":          []interface{}{1, 2, 3, 4, 5},
									"start_time_in_minute": "0",
									"end_time_in_minute":   "1439",
									"time_zone":            "Asia/Shanghai",
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "severity",
											"op":    "IN",
											"value": "CRITICAL,ERROR,WARN",
										},
									},
									"expression": "1 OR 2",
									"relation":   "OR",
								},
							},
							"severities": []string{"CRITICAL"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_strategy_name":         name,
						"description":                  "tf-acc cms notify strategy test",
						"ignore_restored_notification": "false",
						"grouping_setting.#":           "1",
						"routes.#":                     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_strategy_name":         name + "_update",
					"description":                  "tf-acc cms notify strategy test updated",
					"ignore_restored_notification": "true",
					"grouping_setting": []map[string]interface{}{
						{
							"grouping_keys": []string{"severity", "alertname"},
							"period_min":    "10",
							"silence_sec":   "600",
							"times":         "2",
						},
					},
					"routes": []map[string]interface{}{
						{
							"channels": []map[string]interface{}{
								{
									"channel_type":         "CONTACT",
									"receivers":            []string{"${var.cms_contact_id}"},
									"enabled_sub_channels": []string{"EMAIL"},
								},
							},
							"effect_time_range": []map[string]interface{}{
								{
									"day_in_week":          []interface{}{1, 3, 5},
									"start_time_in_minute": "60",
									"end_time_in_minute":   "1380",
									"time_zone":            "UTC",
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "IN",
											"value": "ERROR,WARN",
										},
										{
											"field": "source",
											"op":    "EQ",
											"value": "cms",
										},
									},
									"expression": "1 AND 2",
									"relation":   "AND",
								},
							},
							"severities": []string{"CRITICAL", "ERROR"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_strategy_name":         name + "_update",
						"description":                  "tf-acc cms notify strategy test updated",
						"ignore_restored_notification": "true",
						"grouping_setting.#":           "1",
						"routes.#":                     "1",
					}),
				),
			},
			// Custom notification templates are external objects and this provider has
			// no resource that can create one. Keep schema and diff coverage without
			// sending a fabricated template UUID to the service.
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_template_entries": []map[string]interface{}{
						{
							"target_type":   "MAIL",
							"template_uuid": name + "-tpl-mail",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"routes": []map[string]interface{}{
						{
							"channels": []map[string]interface{}{
								{
									"channel_type":         "GROUP",
									"receivers":            []string{name + "-group"},
									"enabled_sub_channels": []string{"EMAIL", "SMS"},
								},
							},
							"severities": []string{"CRITICAL", "ERROR"},
						},
					},
					"custom_template_entries": []map[string]interface{}{
						{
							"target_type":   "ONCALL",
							"template_uuid": name + "-tpl-oncall",
						},
					},
				}),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAliCloudCmsNotifyStrategy_basic1(t *testing.T) {
	supportedRegions := []connectivity.Region{connectivity.Hangzhou}
	var v map[string]interface{}
	resourceId := "alicloud_cms_notify_strategy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsNotifyStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsNotifyStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := testAccCmsUniqueResourceName()
	notifyTarget := testAccCmsCreateContact(t, rand, supportedRegions)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, func(name string) string {
		return AliCloudCmsNotifyStrategyBasicDependence1(name, notifyTarget)
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, supportedRegions)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_strategy_name": name,
					"workspace":            "${var.cms_workspace}",
					"description":          "tf-acc cms notify strategy workspace test",
					"grouping_setting": []map[string]interface{}{
						{
							"grouping_keys": []string{"severity"},
							"period_min":    "5",
							"silence_sec":   "300",
							"times":         "1",
						},
					},
					"routes": []map[string]interface{}{
						{
							"channels": []map[string]interface{}{
								{
									"channel_type":         "CONTACT",
									"receivers":            []string{"${var.cms_contact_id}"},
									"enabled_sub_channels": []string{"EMAIL"},
								},
							},
							"severities": []string{"CRITICAL"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notify_strategy_name": name,
						"workspace":            CHECKSET,
						"grouping_setting.#":   "1",
						"routes.#":             "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-acc cms notify strategy workspace test updated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-acc cms notify strategy workspace test updated",
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

var AliCloudCmsNotifyStrategyMap0 = map[string]string{
	"notify_strategy_id": CHECKSET,
	"enable":             CHECKSET,
	"create_time":        CHECKSET,
	"update_time":        CHECKSET,
}

func AliCloudCmsNotifyStrategyBasicDependence0(name string, notifyTarget *testAccCmsContact) string {
	return fmt.Sprintf(`
variable "cms_contact_id" {
  default = %q
}

variable "cms_workspace" {
  default = %q
}
`, notifyTarget.Identifier, notifyTarget.Workspace)
}

func AliCloudCmsNotifyStrategyBasicDependence1(name string, notifyTarget *testAccCmsContact) string {
	return AliCloudCmsNotifyStrategyBasicDependence0(name, notifyTarget)
}

// Test Cms NotifyStrategy. <<< Resource test cases.
