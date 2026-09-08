// Package alicloud. This file is hand-written from the Cms 2024-03-30 CloudSpec Mission definition.
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAliCloudCmsMission_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_mission.default"
	ra := resourceAttrInit(resourceId, AliCloudCmsMissionMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsMission")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms-mission%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudCmsMissionBasicDependence)
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
					"name":                  name,
					"digital_employee_name": "${var.digital_employee_name}",
					"display_name":          name,
					"description":           name,
					"enabled":               true,
					"variables": map[string]interface{}{
						"workspace": name,
						"project":   name,
					},
					"notification": []interface{}{
						map[string]interface{}{
							"dingtalk": []interface{}{"https://oapi.dingtalk.com/robot/send?access_token=token1"},
							"feishu":   []interface{}{"https://open.feishu.cn/open-apis/bot/v2/hook/hook1"},
							"slack":    []interface{}{"https://hooks.slack.com/services/slack1"},
							"wechat":   []interface{}{"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=key1"},
							"call":     []interface{}{"13800138000"},
							"sms":      []interface{}{"13800138000"},
							"email":    []interface{}{"test@example.com"},
							"webhook":  []interface{}{"https://example.com/webhook1"},
						},
					},
					"notification_policy": []interface{}{
						map[string]interface{}{
							"region":    "cn-hangzhou",
							"workspace": name,
						},
					},
					"configuration": []interface{}{
						map[string]interface{}{
							"credits": 100,
						},
					},
					"blueprint": []interface{}{
						map[string]interface{}{
							"cron": []interface{}{
								map[string]interface{}{
									"id":              fmt.Sprintf("%s-cron", name),
									"display_name":    fmt.Sprintf("%s-cron", name),
									"description":     "cron blueprint",
									"prompt":          "hello",
									"cron_expression": "0 0 * * *",
									"variables": map[string]interface{}{
										"key": "value1",
									},
									"time_zone":          "+0800",
									"delay":              0,
									"run_immediately":    false,
									"priority":           1,
									"concurrency_policy": "skip",
									"timeout_seconds":    300,
								},
							},
							"calendar": []interface{}{
								map[string]interface{}{
									"id":           fmt.Sprintf("%s-cal", name),
									"display_name": fmt.Sprintf("%s-cal", name),
									"description":  "calendar blueprint",
									"prompt":       "hello",
									"rrule":        "FREQ=DAILY",
									"variables": map[string]interface{}{
										"key": "value1",
									},
									"time_zone":          "+0800",
									"delay":              0,
									"run_immediately":    false,
									"priority":           1,
									"concurrency_policy": "skip",
									"timeout_seconds":    300,
								},
							},
							"event": []interface{}{
								map[string]interface{}{
									"id":               fmt.Sprintf("%s-evt", name),
									"display_name":     fmt.Sprintf("%s-evt", name),
									"description":      "event blueprint",
									"prompt":           "hello",
									"workspace":        name,
									"max_concurrency":  1,
									"debounce_seconds": 0,
									"variables": map[string]interface{}{
										"key": "value1",
									},
									"priority":           1,
									"concurrency_policy": "skip",
									"timeout_seconds":    300,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"digital_employee_name": "${var.digital_employee_name}",
						"display_name":          name,
						"description":           name,
						"enabled":               "true",
						"create_time":           CHECKSET,
						"update_time":           CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"digital_employee_name": "tf-test-digital-employee-v2",
					"display_name":          fmt.Sprintf("%s-updated", name),
					"description":           fmt.Sprintf("%s-updated", name),
					"enabled":               false,
					"variables": map[string]interface{}{
						"workspace": fmt.Sprintf("%s-updated", name),
						"project":   fmt.Sprintf("%s-updated", name),
					},
					"notification": []interface{}{
						map[string]interface{}{
							"dingtalk": []interface{}{"https://oapi.dingtalk.com/robot/send?access_token=token2"},
							"feishu":   []interface{}{"https://open.feishu.cn/open-apis/bot/v2/hook/hook2"},
							"slack":    []interface{}{"https://hooks.slack.com/services/slack2"},
							"wechat":   []interface{}{"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=key2"},
							"call":     []interface{}{"13900139000"},
							"sms":      []interface{}{"13900139000"},
							"email":    []interface{}{"test2@example.com"},
							"webhook":  []interface{}{"https://example.com/webhook2"},
						},
					},
					"notification_policy": []interface{}{
						map[string]interface{}{
							"region":    "cn-beijing",
							"workspace": fmt.Sprintf("%s-updated", name),
						},
					},
					"configuration": []interface{}{
						map[string]interface{}{
							"credits": 200,
						},
					},
					"blueprint": []interface{}{
						map[string]interface{}{
							"cron": []interface{}{
								map[string]interface{}{
									"id":              fmt.Sprintf("%s-cron-v2", name),
									"display_name":    fmt.Sprintf("%s-cron-v2", name),
									"description":     "cron blueprint v2",
									"prompt":          "hello v2",
									"cron_expression": "30 0 * * *",
									"variables": map[string]interface{}{
										"key": "value2",
									},
									"time_zone":          "+0900",
									"delay":              10,
									"run_immediately":    true,
									"priority":           5,
									"concurrency_policy": "queue",
									"timeout_seconds":    600,
								},
							},
							"calendar": []interface{}{
								map[string]interface{}{
									"id":           fmt.Sprintf("%s-cal-v2", name),
									"display_name": fmt.Sprintf("%s-cal-v2", name),
									"description":  "calendar blueprint v2",
									"prompt":       "hello v2",
									"rrule":        "FREQ=HOURLY",
									"variables": map[string]interface{}{
										"key": "value2",
									},
									"time_zone":          "+0900",
									"delay":              10,
									"run_immediately":    true,
									"priority":           5,
									"concurrency_policy": "queue",
									"timeout_seconds":    600,
								},
							},
							"event": []interface{}{
								map[string]interface{}{
									"id":               fmt.Sprintf("%s-evt-v2", name),
									"display_name":     fmt.Sprintf("%s-evt-v2", name),
									"description":      "event blueprint v2",
									"prompt":           "hello v2",
									"workspace":        fmt.Sprintf("%s-updated", name),
									"max_concurrency":  2,
									"debounce_seconds": 5,
									"variables": map[string]interface{}{
										"key": "value2",
									},
									"priority":           5,
									"concurrency_policy": "queue",
									"timeout_seconds":    600,
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":         name,
						"display_name": fmt.Sprintf("%s-updated", name),
						"description":  fmt.Sprintf("%s-updated", name),
						"enabled":      "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccAliCloudCmsMissionsDataSource_basic(t *testing.T) {
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms-mission-ds%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc("data.alicloud_cms_missions.default", name, AliCloudCmsMissionBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: func(*terraform.State) error {
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"digital_employee_name": "${var.digital_employee_name}",
					"name_regex":            name,
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.alicloud_cms_missions.default", "missions.#"),
				),
			},
		},
	})
}

var AliCloudCmsMissionMap = map[string]string{
	"create_time": CHECKSET,
	"update_time": CHECKSET,
}

func AliCloudCmsMissionBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "digital_employee_name" {
    default = "tf-test-digital-employee"
}
`, name)
}
