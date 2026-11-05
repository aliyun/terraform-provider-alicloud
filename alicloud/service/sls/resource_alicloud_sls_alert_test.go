package sls_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/sls"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Test Sls Alert. >>> Resource test cases, automatically generated.
// Case Alert_Terraform_Schedule 5844
func TestAccAliCloudSlsAlert_basic5844(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
								{
									"store_type": "meta",
									"store":      "mytest1",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "no_group",
									"fields": []string{
										"a", "b"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "qqq",
								},
								{
									"type":      "inner_join",
									"condition": "fefefe",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-shanghai-intranet.log.aliyuncs.com",
									"project":     "wkb-wangren",
									"event_store": "alert",
									"role_arn":    "acs:ram::1654218965343050:role/aliyunlogetlrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"template_id": "sls.app.ack.autoscaler.cluster_unhealthy",
									"type":        "sys",
									"version":     "1.0",
									"lang":        "cn",
									"tokens": map[string]interface{}{
										"\"a\"": "b",
									},
									"annotations": map[string]interface{}{
										"\"x\"": "y",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 3",
									"count_condition": "__count__ < 3",
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic",
									"action_policy_id": "wkb-action",
									"repeat_interval":  "1m",
								},
							},
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
							"tags": []string{
								"wkb", "wangren", "sls"},
						},
					},
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"display_name": "openapi-terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "create alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "create alert",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":      "default",
							"version":   "2.0",
							"dashboard": "internal-alert-analysis",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select count(*)",
									"time_span_type": "Custom",
									"start":          "-5m",
									"end":            "-2m",
									"store_type":     "metric",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "aler",
									"region":         "cn-shanghai",
									"power_sql_mode": "auto",
									"chart_title":    "wkb-chart-1",
									"dashboard_id":   "wkb-dashboard-1",
									"ui":             "{\\\"app\\\": \\\"trace\\\"}",
								},
								{
									"store_type": "meta",
									"store":      " user.es.authtest",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "custom",
									"fields": []string{
										"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "cross_join",
									"condition": "$0.id == $1.id",
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a1",
									"value": "b1",
								},
								{
									"key":   "a11",
									"value": "b11",
								},
								{
									"key":   "a111",
									"value": "b111",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x1",
									"value": "y1",
								},
								{
									"key":   "x11",
									"value": "y11",
								},
								{
									"key":   "x111",
									"value": "y111",
								},
							},
							"auto_annotation": "false",
							"send_resolved":   "true",
							"threshold":       "2",
							"no_data_fire":    "true",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-beijing-intranet.log.aliyuncs.com",
									"project":     "${alicloud_log_project.defaultINsMgl.name}",
									"event_store": "test",
									"role_arn":    " acs:ram::1654218965343050:role/aliyunlogetlrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "10",
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt < 120",
										},
									},
								},
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt < 200",
										},
									},
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic)",
									"action_policy_id": "wkb-test",
									"repeat_interval":  "1h",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"type":        "user",
									"version":     "2.0",
									"lang":        "en",
									"template_id": "sls.app.ack.autoscaler.instance_expired",
									"tokens": map[string]interface{}{
										"\"a1\"": "b1",
									},
									"annotations": map[string]interface{}{
										"\"x1\"": "y1",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 5",
									"count_condition": "__count__ > 10",
								},
							},
							"mute_until":       "1",
							"no_data_severity": "4",
							"tags": []string{
								"wkb", "wangren"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"delay":           "20",
							"time_zone":       "+0700",
							"cron_expression": "0 0/1 * * *",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "openapi-terraform-step1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "openapi-terraform-step1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":      "default",
							"version":   "2.0",
							"dashboard": "internal-alert-analysis",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "wkb-wangren",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"role_arn":       "acs:ram::1654218965343050:role/alert-monitor-role",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "custom",
									"fields": []string{},
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "2",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 6",
											"count_condition": "cnt < 30",
										},
									},
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"tags": []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "false",
							"cron_expression": "0 0/1 * * *",
							"delay":           "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "openapi-terraform-step1-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "openapi-terraform-step1-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "this is alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "this is alert",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
								{
									"store_type": "meta",
									"store":      "mytest1",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "no_group",
									"fields": []string{
										"a", "b"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "qqq",
								},
								{
									"type":      "inner_join",
									"condition": "fefefe",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-shanghai-intranet.log.aliyuncs.com",
									"project":     "wkb-wangren",
									"event_store": "alert",
									"role_arn":    "acs:ram::1654218965343050:role/aliyunlogetlrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"template_id": "sls.app.ack.autoscaler.cluster_unhealthy",
									"type":        "sys",
									"version":     "1.0",
									"lang":        "cn",
									"tokens": map[string]interface{}{
										"\"a\"": "b",
									},
									"annotations": map[string]interface{}{
										"\"x\"": "y",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 3",
									"count_condition": "__count__ < 3",
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic",
									"action_policy_id": "wkb-action",
									"repeat_interval":  "1m",
								},
							},
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
							"tags": []string{
								"wkb", "wangren", "sls"},
						},
					},
					"alert_name":   name + "_update",
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"display_name": "openapi-terraform",
					"description":  "create alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name + "_update",
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
						"description":  "create alert",
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

var AlicloudSlsAlertMap5844 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudSlsAlertBasicDependence5844(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "alert_name" {
  default = "openapi-terraform-alert"
}

variable "project_name" {
  default = "terraform-alert-test"
}

resource "alicloud_log_project" "defaultINsMgl" {
  description = "terraform-alert-test"
  name        = var.name
}


`, name)
}

// Case Alert_Terraform 6423
func TestAccAliCloudSlsAlert_basic6423(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap6423)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence6423)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
								{
									"store_type": "meta",
									"store":      "mytest1",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "no_group",
									"fields": []string{
										"a", "b"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "qqq",
								},
								{
									"type":      "inner_join",
									"condition": "fefefe",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-shanghai-intranet.log.aliyuncs.com",
									"project":     "wkb-wangren",
									"event_store": "alert",
									"role_arn":    "acs:ram::1654218965343050:role/aliyunlogetlrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"template_id": "sls.app.ack.autoscaler.cluster_unhealthy",
									"type":        "sys",
									"version":     "1.0",
									"lang":        "cn",
									"tokens": map[string]interface{}{
										"\"a\"": "b",
									},
									"annotations": map[string]interface{}{
										"\"x\"": "y",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 3",
									"count_condition": "__count__ < 3",
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic",
									"action_policy_id": "wkb-action",
									"repeat_interval":  "1m",
								},
							},
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
							"tags": []string{
								"wkb", "wangren", "sls"},
						},
					},
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"schedule": []map[string]interface{}{
						{
							"type":           "FixedRate",
							"run_immdiately": "true",
							"interval":       "1m",
							"time_zone":      "+0800",
							"delay":          "10",
						},
					},
					"display_name": "openapi-terraform",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "create alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "create alert",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ENABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":      "default",
							"version":   "2.0",
							"dashboard": "internal-alert-analysis",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select count(*)",
									"time_span_type": "Custom",
									"start":          "-5m",
									"end":            "-2m",
									"store_type":     "metric",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "aler",
									"region":         "cn-hangzhou",
									"power_sql_mode": "auto",
									"chart_title":    "wkb-chart-1",
									"dashboard_id":   "wkb-dashboard-1",
									"ui":             "{\\\"app\\\": \\\"trace\\\"}",
								},
								{
									"store_type": "meta",
									"store":      " user.es.authtest",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "custom",
									"fields": []string{
										"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "cross_join",
									"condition": "$0.id == $1.id",
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a1",
									"value": "b1",
								},
								{
									"key":   "a11",
									"value": "b11",
								},
								{
									"key":   "a111",
									"value": "b111",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x1",
									"value": "y1",
								},
								{
									"key":   "x11",
									"value": "y11",
								},
								{
									"key":   "x111",
									"value": "y111",
								},
							},
							"auto_annotation": "false",
							"send_resolved":   "true",
							"threshold":       "2",
							"no_data_fire":    "true",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-beijing-intranet.log.aliyuncs.com",
									"project":     "${alicloud_log_project.defaultINsMgl.name}",
									"event_store": "test",
									"role_arn":    " acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "true",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "10",
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt < 120",
										},
									},
								},
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt < 200",
										},
									},
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic)",
									"action_policy_id": "wkb-test",
									"repeat_interval":  "1h",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"type":        "user",
									"version":     "2.0",
									"lang":        "en",
									"template_id": "sls.app.ack.autoscaler.instance_expired",
									"tokens": map[string]interface{}{
										"\"a1\"": "b1",
									},
									"annotations": map[string]interface{}{
										"\"x1\"": "y1",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 5",
									"count_condition": "__count__ > 10",
								},
							},
							"mute_until":       "1",
							"no_data_severity": "4",
							"tags": []string{
								"wkb", "wangren"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"type":           "FixedRate",
							"run_immdiately": "true",
							"interval":       "5m",
							"delay":          "20",
							"time_zone":      "+0700",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "openapi-terraform-step1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "openapi-terraform-step1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":      "default",
							"version":   "2.0",
							"dashboard": "internal-alert-analysis",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "wkb-wangren",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"role_arn":       "acs:ram::1654218965343050:role/alert-monitor-role",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "custom",
									"fields": []string{},
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "2",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 6",
											"count_condition": "cnt < 30",
										},
									},
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"tags": []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "false",
							"cron_expression": "0 0/1 * * *",
							"delay":           "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "openapi-terraform-step1-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "openapi-terraform-step1-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "this is alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "this is alert",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "DISABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "DISABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "ENABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "ENABLED",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
								{
									"store_type": "meta",
									"store":      "mytest1",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "no_group",
									"fields": []string{
										"a", "b"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "qqq",
								},
								{
									"type":      "inner_join",
									"condition": "fefefe",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-shanghai-intranet.log.aliyuncs.com",
									"project":     "wkb-wangren",
									"event_store": "alert",
									"role_arn":    "acs:ram::1654218965343050:role/aliyunlogetlrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"template_id": "sls.app.ack.autoscaler.cluster_unhealthy",
									"type":        "sys",
									"version":     "1.0",
									"lang":        "cn",
									"tokens": map[string]interface{}{
										"\"a\"": "b",
									},
									"annotations": map[string]interface{}{
										"\"x\"": "y",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 3",
									"count_condition": "__count__ < 3",
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic",
									"action_policy_id": "wkb-action",
									"repeat_interval":  "1m",
								},
							},
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
							"tags": []string{
								"wkb", "wangren", "sls"},
						},
					},
					"alert_name":   name + "_update",
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"schedule": []map[string]interface{}{
						{
							"type":           "FixedRate",
							"run_immdiately": "true",
							"interval":       "1m",
							"time_zone":      "+0800",
							"delay":          "10",
						},
					},
					"display_name": "openapi-terraform",
					"description":  "create alert",
					"status":       "ENABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name + "_update",
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
						"description":  "create alert",
						"status":       "ENABLED",
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

var AlicloudSlsAlertMap6423 = map[string]string{
	"status":      acctest.CHECKSET,
	"create_time": acctest.CHECKSET,
}

func AlicloudSlsAlertBasicDependence6423(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

variable "alert_name" {
  default = "openapi-terraform-alert-new"
}

variable "project_name" {
  default = "terraform-alert-test-new"
}

resource "alicloud_log_project" "defaultINsMgl" {
  description = "terraform-alert-test"
  name        = var.name
}


`, name)
}

// Case Alert_Terraform_Schedule 5844  twin
func TestAccAliCloudSlsAlert_basic5844_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
								{
									"store_type": "meta",
									"store":      "mytest1",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "no_group",
									"fields": []string{
										"a", "b"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "qqq",
								},
								{
									"type":      "inner_join",
									"condition": "fefefe",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-shanghai-intranet.log.aliyuncs.com",
									"project":     "wkb-wangren",
									"event_store": "alert",
									"role_arn":    "acs:ram::1654218965343050:role/aliyunlogetlrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"template_id": "sls.app.ack.autoscaler.cluster_unhealthy",
									"type":        "sys",
									"version":     "1.0",
									"lang":        "cn",
									"tokens": map[string]interface{}{
										"\"a\"": "b",
									},
									"annotations": map[string]interface{}{
										"\"x\"": "y",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 3",
									"count_condition": "__count__ < 3",
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic",
									"action_policy_id": "wkb-action",
									"repeat_interval":  "1m",
								},
							},
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
							"tags": []string{
								"wkb", "wangren", "sls"},
						},
					},
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"display_name": "openapi-terraform",
					"description":  "create alert",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
						"description":  "create alert",
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

// Case Alert_Terraform 6423  twin
func TestAccAliCloudSlsAlert_basic6423_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap6423)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence6423)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
								{
									"store_type": "meta",
									"store":      "mytest1",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type": "no_group",
									"fields": []string{
										"a", "b"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "qqq",
								},
								{
									"type":      "inner_join",
									"condition": "fefefe",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
							},
							"auto_annotation": "true",
							"send_resolved":   "false",
							"threshold":       "1",
							"no_data_fire":    "false",
							"sink_event_store": []map[string]interface{}{
								{
									"enabled":     "true",
									"endpoint":    "cn-shanghai-intranet.log.aliyuncs.com",
									"project":     "wkb-wangren",
									"event_store": "alert",
									"role_arn":    "acs:ram::1654218965343050:role/aliyunlogetlrole",
								},
							},
							"sink_cms": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"sink_alerthub": []map[string]interface{}{
								{
									"enabled": "false",
								},
							},
							"template_configuration": []map[string]interface{}{
								{
									"template_id": "sls.app.ack.autoscaler.cluster_unhealthy",
									"type":        "sys",
									"version":     "1.0",
									"lang":        "cn",
									"tokens": map[string]interface{}{
										"\"a\"": "b",
									},
									"annotations": map[string]interface{}{
										"\"x\"": "y",
									},
								},
							},
							"condition_configuration": []map[string]interface{}{
								{
									"condition":       "cnt > 3",
									"count_condition": "__count__ < 3",
								},
							},
							"policy_configuration": []map[string]interface{}{
								{
									"alert_policy_id":  "sls.builtin.dynamic",
									"action_policy_id": "wkb-action",
									"repeat_interval":  "1m",
								},
							},
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
							"tags": []string{
								"wkb", "wangren", "sls"},
						},
					},
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"schedule": []map[string]interface{}{
						{
							"type":           "FixedRate",
							"run_immdiately": "true",
							"interval":       "1m",
							"time_zone":      "+0800",
							"delay":          "10",
						},
					},
					"display_name": "openapi-terraform",
					"description":  "create alert",
					"status":       "ENABLED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
						"description":  "create alert",
						"status":       "ENABLED",
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

// TestAccAliCloudSlsAlert_typeListOrder_memberOrder_annotations covers the
// configuration.annotations TypeList path. Pattern B (nested): step1 sets the full
// configuration block with all nested lists in order A (>=2 distinct
// members each); step2 reorders ONLY configuration.annotations, keeping every other
// attribute byte-identical to step1.
func TestAccAliCloudSlsAlert_typeListOrder_memberOrder_annotations(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x2",
									"value": "y2",
								},
								{
									"key":   "x",
									"value": "y",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsAlert_typeListOrder_memberOrder_groupConfigurationFields covers the
// configuration.group_configuration.fields TypeList path. Pattern B (nested): step1 sets the full
// configuration block with all nested lists in order A (>=2 distinct
// members each); step2 reorders ONLY configuration.group_configuration.fields, keeping every other
// attribute byte-identical to step1.
func TestAccAliCloudSlsAlert_typeListOrder_memberOrder_groupConfigurationFields(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"c", "b", "a"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsAlert_typeListOrder_memberOrder_joinConfigurations covers the
// configuration.join_configurations TypeList path. Pattern B (nested): step1 sets the full
// configuration block with all nested lists in order A (>=2 distinct
// members each); step2 reorders ONLY configuration.join_configurations, keeping every other
// attribute byte-identical to step1.
func TestAccAliCloudSlsAlert_typeListOrder_memberOrder_joinConfigurations(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "cross_join",
									"condition": "bb",
								},
								{
									"type":      "no_join",
									"condition": "aa",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsAlert_typeListOrder_memberOrder_labels covers the
// configuration.labels TypeList path. Pattern B (nested): step1 sets the full
// configuration block with all nested lists in order A (>=2 distinct
// members each); step2 reorders ONLY configuration.labels, keeping every other
// attribute byte-identical to step1.
func TestAccAliCloudSlsAlert_typeListOrder_memberOrder_labels(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "c",
									"value": "d",
								},
								{
									"key":   "a",
									"value": "b",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsAlert_typeListOrder_memberOrder_queryList covers the
// configuration.query_list TypeList path. Pattern B (nested): step1 sets the full
// configuration block with all nested lists in order A (>=2 distinct
// members each); step2 reorders ONLY configuration.query_list, keeping every other
// attribute byte-identical to step1.
func TestAccAliCloudSlsAlert_typeListOrder_memberOrder_queryList(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsAlert_typeListOrder_memberOrder_severityConfigurations covers the
// configuration.severity_configurations TypeList path. Pattern B (nested): step1 sets the full
// configuration block with all nested lists in order A (>=2 distinct
// members each); step2 reorders ONLY configuration.severity_configurations, keeping every other
// attribute byte-identical to step1.
func TestAccAliCloudSlsAlert_typeListOrder_memberOrder_severityConfigurations(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
		},
	})
}

// TestAccAliCloudSlsAlert_typeListOrder_memberOrder_tags covers the
// configuration.tags TypeList path. Pattern B (nested): step1 sets the full
// configuration block with all nested lists in order A (>=2 distinct
// members each); step2 reorders ONLY configuration.tags, keeping every other
// attribute byte-identical to step1.
func TestAccAliCloudSlsAlert_typeListOrder_memberOrder_tags(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_sls_alert.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudSlsAlertMap5844)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return sls.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeSlsAlert")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sslsalert%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudSlsAlertBasicDependence5844)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"wkb", "wangren", "sls"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_name":   name,
					"project_name": "${alicloud_log_project.defaultINsMgl.name}",
					"display_name": "openapi-terraform",
					"schedule": []map[string]interface{}{
						{
							"type":            "Cron",
							"run_immdiately":  "true",
							"time_zone":       "+0800",
							"delay":           "10",
							"cron_expression": "0/5 * * * *",
						},
					},
					"configuration": []map[string]interface{}{
						{
							"type":    "tpl",
							"version": "2",
							"query_list": []map[string]interface{}{
								{
									"query":          "* | select *",
									"time_span_type": "Relative",
									"start":          "-15m",
									"end":            "now",
									"store_type":     "log",
									"project":        "${alicloud_log_project.defaultINsMgl.name}",
									"store":          "alert",
									"region":         "cn-shanghai",
									"power_sql_mode": "disable",
									"chart_title":    "wkb-chart",
									"dashboard_id":   "wkb-dashboard",
									"ui":             "{}",
									"role_arn":       "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole",
								},
								{
									"store_type": "meta",
									"store":      "user.rds_ip_whitelist",
								},
							},
							"group_configuration": []map[string]interface{}{
								{
									"type":   "no_group",
									"fields": []string{"a", "b", "c"},
								},
							},
							"join_configurations": []map[string]interface{}{
								{
									"type":      "no_join",
									"condition": "aa",
								},
								{
									"type":      "cross_join",
									"condition": "bb",
								},
							},
							"severity_configurations": []map[string]interface{}{
								{
									"severity": "6",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 1",
											"count_condition": "cnt > 0",
										},
									},
								},
								{
									"severity": "8",
									"eval_condition": []map[string]interface{}{
										{
											"condition":       "__count__ > 2",
											"count_condition": "cnt > 1",
										},
									},
								},
							},
							"labels": []map[string]interface{}{
								{
									"key":   "a",
									"value": "b",
								},
								{
									"key":   "c",
									"value": "d",
								},
							},
							"annotations": []map[string]interface{}{
								{
									"key":   "x",
									"value": "y",
								},
								{
									"key":   "x2",
									"value": "y2",
								},
							},
							"tags":             []string{"sls", "wangren", "wkb"},
							"auto_annotation":  "true",
							"send_resolved":    "false",
							"threshold":        "1",
							"no_data_fire":     "false",
							"dashboard":        "internal-alert",
							"mute_until":       "0",
							"no_data_severity": "6",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":   name,
						"project_name": acctest.CHECKSET,
						"display_name": "openapi-terraform",
					}),
				),
			},
		},
	})
}

// Test Sls Alert. <<< Resource test cases, automatically generated.
