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
	var v map[string]interface{}
	resourceId := "alicloud_cms_notify_strategy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsNotifyStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsNotifyStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsNotifyStrategyBasicDependence0)
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
					"notify_strategy_name":         name,
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
									"receivers":            []string{"DING", "WEIXIN"},
									"enabled_sub_channels": []string{"SMS", "EMAIL"},
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
					"custom_template_entries": []map[string]interface{}{
						{
							"target_type":   "MAIL",
							"template_uuid": name + "-tpl-mail",
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
						"custom_template_entries.#":    "1",
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
									"channel_type":         "FEISHU",
									"receivers":            []string{"GROUP"},
									"enabled_sub_channels": []string{},
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
					"custom_template_entries": []map[string]interface{}{
						{
							"target_type":   "ONCALL",
							"template_uuid": name + "-tpl-oncall",
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
						"custom_template_entries.#":    "1",
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

func TestAccAliCloudCmsNotifyStrategy_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_notify_strategy.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsNotifyStrategyMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsNotifyStrategy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsNotifyStrategyBasicDependence1)
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
					"notify_strategy_name": name,
					"workspace":            "${alicloud_cms_workspace.default.id}",
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
									"channel_type": "DING",
									"receivers":    []string{"CONTACT"},
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

func AliCloudCmsNotifyStrategyBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
`, name)
}

func AliCloudCmsNotifyStrategyBasicDependence1(name string) string {
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
`, name)
}

// Test Cms NotifyStrategy. <<< Resource test cases.
