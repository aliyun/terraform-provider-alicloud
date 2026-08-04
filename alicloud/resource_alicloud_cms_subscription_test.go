package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms Subscription. >>> Resource test cases.

func TestAccAliCloudCmsSubscription_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_subscription.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsSubscriptionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsSubscription")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsSubscriptionBasicDependence0)
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
					"subscription_name":  name,
					"description":        "tf-acc cms subscription test",
					"notify_strategy_id": "${alicloud_cms_notify_strategy.default.id}",
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
					"pushing_setting": []map[string]interface{}{
						{
							"template_uuid":      name + "-tpl-1",
							"alert_action_ids":   []string{name + "-alert-1"},
							"restore_action_ids": []string{name + "-restore-1"},
							"response_plan_id":   name + "-resp-1",
						},
					},
					"agent_config": []map[string]interface{}{
						{
							"agent_uuid": name + "-agent-1",
							"routes": []map[string]interface{}{
								{
									"channels": []map[string]interface{}{
										{
											"channel_type":         "CONTACT",
											"receivers":            []string{"agent-receiver-1"},
											"enabled_sub_channels": []string{"SMS", "EMAIL"},
										},
									},
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"day_in_week":          []interface{}{1, 2, 3, 4, 5},
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subscription_name":  name,
						"description":        "tf-acc cms subscription test",
						"notify_strategy_id": CHECKSET,
						"filter_setting.#":   "1",
						"pushing_setting.#":  "1",
						"agent_config.#":     "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"subscription_name":  name + "_update",
					"description":        "tf-acc cms subscription test updated",
					"notify_strategy_id": "${alicloud_cms_notify_strategy.second.id}",
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
					"pushing_setting": []map[string]interface{}{
						{
							"template_uuid":      name + "-tpl-2",
							"alert_action_ids":   []string{name + "-alert-2"},
							"restore_action_ids": []string{name + "-restore-2"},
							"response_plan_id":   name + "-resp-2",
						},
					},
					"agent_config": []map[string]interface{}{
						{
							"agent_uuid": name + "-agent-2",
							"routes": []map[string]interface{}{
								{
									"channels": []map[string]interface{}{
										{
											"channel_type":         "CONTACT",
											"receivers":            []string{"agent-receiver-2", "agent-receiver-3"},
											"enabled_sub_channels": []string{"EMAIL", "CALL", "SMS"},
										},
										{
											"channel_type":         "WEBHOOK",
											"receivers":            []string{"agent-receiver-7"},
											"enabled_sub_channels": []string{},
										},
									},
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "UTC",
											"day_in_week":          []interface{}{0, 1, 2, 3, 4, 5},
											"start_time_in_minute": "60",
											"end_time_in_minute":   "1380",
										},
									},
								},
								{
									"channels": []map[string]interface{}{
										{
											"channel_type":         "CONTACT",
											"receivers":            []string{"agent-receiver-4"},
											"enabled_sub_channels": []string{"SMS"},
										},
									},
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"day_in_week":          []interface{}{1, 3, 5},
											"start_time_in_minute": "120",
											"end_time_in_minute":   "1320",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subscription_name":  name + "_update",
						"description":        "tf-acc cms subscription test updated",
						"notify_strategy_id": CHECKSET,
						"filter_setting.#":   "1",
						"pushing_setting.#":  "1",
						"agent_config.#":     "1",
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

func TestAccAliCloudCmsSubscription_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_subscription.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsSubscriptionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsSubscription")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsSubscriptionBasicDependence1)
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
					"subscription_name":  name,
					"workspace":          "${alicloud_cms_workspace.default.id}",
					"notify_strategy_id": "${alicloud_cms_notify_strategy.default.id}",
					"filter_setting": []map[string]interface{}{
						{
							"conditions": []map[string]interface{}{
								{
									"field": "severity",
									"op":    "EQ",
									"value": "CRITICAL",
								},
							},
							"relation": "AND",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"subscription_name":  name,
						"workspace":          CHECKSET,
						"notify_strategy_id": CHECKSET,
						"filter_setting.#":   "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "tf-acc cms subscription workspace test",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "tf-acc cms subscription workspace test",
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

var AliCloudCmsSubscriptionMap0 = map[string]string{
	"subscription_id": CHECKSET,
	"enable":          CHECKSET,
	"create_time":     CHECKSET,
	"update_time":     CHECKSET,
}

func AliCloudCmsSubscriptionBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_cms_notify_strategy" "default" {
  notify_strategy_name = var.name
  grouping_setting {
    grouping_keys = ["severity"]
    period_min    = 5
    silence_sec   = 300
    times         = 1
  }
  routes {
    channels {
      channel_type = "DING"
      receivers    = ["CONTACT"]
    }
    severities = ["CRITICAL"]
  }
}

resource "alicloud_cms_notify_strategy" "second" {
  notify_strategy_name = "${var.name}-2"
  grouping_setting {
    grouping_keys = ["severity"]
    period_min    = 5
    silence_sec   = 300
    times         = 1
  }
  routes {
    channels {
      channel_type = "DING"
      receivers    = ["CONTACT"]
    }
    severities = ["CRITICAL"]
  }
}
`, name)
}

func AliCloudCmsSubscriptionBasicDependence1(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
}

resource "alicloud_cms_notify_strategy" "default" {
  notify_strategy_name = var.name
  grouping_setting {
    grouping_keys = ["severity"]
    period_min    = 5
    silence_sec   = 300
    times         = 1
  }
  routes {
    channels {
      channel_type = "DING"
      receivers    = ["CONTACT"]
    }
    severities = ["CRITICAL"]
  }
}
`, name)
}

// Test Cms Subscription. <<< Resource test cases.
