package alicloud

import (
	"fmt"
	"testing"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudLogAlert_basic(t *testing.T) {
	var alert *sls.Alert
	resourceId := "alicloud_log_alert.default"
	ra := resourceAttrInit(resourceId, logAlertMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &alert, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogalert-%d", rand)
	displayname := fmt.Sprintf("alert_displayname-%d", rand)
	content := "aliyun sls alert test"
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogAlertDependence)

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
					"project_name":      "${alicloud_log_project.default.name}",
					"alert_name":        "alert_name",
					"alert_displayname": displayname,
					"condition":         "count >100",
					"dashboard":         "terraform-dashboard",
					"query_list": []map[string]interface{}{
						{
							"logstore":    "${alicloud_log_store.default.name}",
							"chart_title": "chart_title",
							"start":       "-60s",
							"end":         "20s",
							"query":       "* AND aliyun",
						},
					},
					"notification_list": []map[string]interface{}{
						{
							"type":        "SMS",
							"mobile_list": []string{"18865521787", "123456678"},
							"content":     content,
						},
						{
							"type":       "Email",
							"email_list": []string{"nihao@alibaba-inc.com", "test@123.com"},
							"content":    content,
						},
						{
							"type":        "DingTalk",
							"service_uri": "www.aliyun.com",
							"content":     content,
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_name":        "alert_name",
						"condition":         "count >100",
						"alert_displayname": displayname,
						"dashboard":         "terraform-dashboard",
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
					"throttling": "1h",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"throttling": "1h",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"throttling": "60s",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"throttling": "60s",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"throttling": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"throttling": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_displayname": "update_alert_name_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_displayname": "update_alert_name_new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"condition": "count>999",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"condition": "count>999",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"dashboard": "dashboard_update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"dashboard": "dashboard_update",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"alert_displayname": "update_alert_name",
					"condition":         "count<100",
					"dashboard":         "terraform-dashboard-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_displayname": "update_alert_name",
						"condition":         "count<100",
						"dashboard":         "terraform-dashboard-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_list": []map[string]interface{}{
						{
							"logstore":    "${alicloud_log_store.default.name}",
							"chart_title": "chart_title",
							"start":       "-80s",
							"end":         "60s",
							"query":       "* AND aliyun_update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_list.#":       "1",
						"query_list.0.start": "-80s",
						"query_list.0.end":   "60s",
						"query_list.0.query": "* AND aliyun_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notification_list": []map[string]interface{}{
						{
							"type":        "SMS",
							"mobile_list": []string{"456456", "456456456"},
							"content":     "updatecontent",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_list.#":         "1",
						"notification_list.0.type":    "SMS",
						"notification_list.0.content": "updatecontent",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudLogAlert_basic_new(t *testing.T) {
	var alert *sls.Alert
	resourceId := "alicloud_log_alert.default"
	ra := resourceAttrInit(resourceId, logAlertMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &alert, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogalert-%d", rand)
	displayname := fmt.Sprintf("alert_displayname-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogAlertDependence)

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
					"version":           "2.0",
					"type":              "default",
					"mute_until":        "1632486684",
					"no_data_fire":      "false",
					"no_data_severity":  "8",
					"project_name":      "${alicloud_log_project.default.name}",
					"alert_name":        "alert_name",
					"alert_displayname": displayname,
					"send_resolved":     "true",
					"query_list": []map[string]interface{}{
						{
							"store":       "${alicloud_log_store.default.name}",
							"storeType":   "log",
							"region":      "cn-heyuan",
							"chart_title": "chart_title",
							"start":       "-60s",
							"end":         "20s",
							"query":       "* AND aliyun | select count(1) as cnt",
						},
						{
							"store":       "${alicloud_log_store.default.name}",
							"storeType":   "log",
							"region":      "cn-heyuan",
							"chart_title": "chart_title",
							"start":       "-60s",
							"end":         "20s",
							"query":       "error | select count(1) as error_cnt",
						},
					},
					"labels": []map[string]interface{}{
						{
							"key":   "env",
							"value": "test",
						},
					},
					"annotations": []map[string]interface{}{
						{
							"key":   "title",
							"value": "alert title",
						},
						{
							"key":   "desc",
							"value": "alert desc",
						},
						{
							"key":   "test_key",
							"value": "test value",
						},
					},
					"group_configuration": map[string]interface{}{
						"type":   "no_group",
						"fields": []string{},
					},
					"policy_configuration": map[string]interface{}{
						"alert_policy_id":  "sls.builtin.dynamic",
						"action_policy_id": "test_alert_policy",
						"repeat_interval":  "1h",
					},
					"join_configuration": map[string]interface{}{
						"type":      "cross_join",
						"condition": "",
					},
					"severity_configurations": []map[string]interface{}{
						{
							"severity": "8",
							"eval_condition": map[string]interface{}{
								"condition":       "cnt > 3",
								"count_condition": "__count__ > 3",
							},
						},
						{
							"severity": "6",
							"eval_condition": map[string]interface{}{
								"condition":       "cnt > 1",
								"count_condition": "__count__ > 1",
							},
						},
						{
							"severity": "2",
							"eval_condition": map[string]interface{}{
								"condition":       "",
								"count_condition": "__count__ > 0",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"version":                               "2.0",
						"type":                                  "default",
						"mute_until":                            "1632486684",
						"no_data_fire":                          "false",
						"no_data_severity":                      "8",
						"project_name":                          "${alicloud_log_project.default.name}",
						"alert_name":                            "alert_name",
						"alert_displayname":                     displayname,
						"send_resolved":                         "true",
						"query_list.#":                          "2",
						"query_list.0.store":                    "${alicloud_log_store.default.name}",
						"query_list.0.storeType":                "log",
						"query_list.0.region":                   "cn-heyuan",
						"query_list.0.chart_title":              "chart_title",
						"query_list.0.start":                    "-60s",
						"query_list.0.end":                      "20s",
						"query_list.0.query":                    "* AND aliyun | select count(1) as cnt",
						"query_list.1.store":                    "${alicloud_log_store.default.name}",
						"query_list.1.storeType":                "log",
						"query_list.1.region":                   "cn-heyuan",
						"query_list.1.chart_title":              "chart_title",
						"query_list.1.start":                    "-60s",
						"query_list.1.end":                      "20s",
						"query_list.1.query":                    "error | select count(1) as error_cnt",
						"labels.#":                              "1",
						"labels.0.key":                          "env",
						"labels.0.value":                        "test",
						"annotations.#":                         "3",
						"annotations.0.key":                     "title",
						"annotations.0.value":                   "alert title",
						"annotations.1.key":                     "desc",
						"annotations.1.value":                   "alert desc",
						"annotations.2.key":                     "test_key",
						"annotations.2.value":                   "test value",
						"group_configuration.type":              "no_group",
						"policy_configuration.alert_policy_id":  "sls.builtin.dynamic",
						"policy_configuration.action_policy_id": "test_alert_policy",
						"policy_configuration.repeat_interval":  "1h",
						"join_configuration.type":               "cross_join",
						"join_configuration.condition":          "",
						"severity_configurations.#":             "3",
						"severity_configurations.0.severity":    "8",
						"severity_configurations.0.eval_condition.condition":       "cnt > 3",
						"severity_configurations.0.eval_condition.count_condition": "__count__ > 3",
						"severity_configurations.1.severity":                       "6",
						"severity_configurations.1.eval_condition.condition":       "cnt > 1",
						"severity_configurations.1.eval_condition.count_condition": "__count__ > 1",
						"severity_configurations.2.severity":                       "2",
						"severity_configurations.2.eval_condition.condition":       "",
						"severity_configurations.2.eval_condition.count_condition": "__count__ > 0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mute_until": "1632488888",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mute_until": "1632488888",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_fire": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_fire": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_resolved": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_resolved": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"alert_displayname": "update_alert_name_new",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"alert_displayname": "update_alert_name_new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_list": []map[string]interface{}{
						{
							"store":       "${alicloud_log_store.default.name}",
							"chart_title": "chart_title",
							"start":       "-80s",
							"end":         "60s",
							"query":       "* AND aliyun_update",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"query_list.#":       "1",
						"query_list.0.start": "-80s",
						"query_list.0.end":   "60s",
						"query_list.0.query": "* AND aliyun_update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"key":   "env",
							"value": "test new",
						},
						{
							"key":   "team",
							"value": "test team",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#":       "2",
						"labels.0.key":   "env",
						"labels.0.value": "test new",
						"labels.1.key":   "team",
						"labels.1.value": "test team",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"annotations": []map[string]interface{}{
						{
							"key":   "title",
							"value": "alert title new",
						},
						{
							"key":   "desc",
							"value": "alert desc new",
						},
						{
							"key":   "test_key",
							"value": "test value new",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"annotations.#":       "2",
						"annotations.0.key":   "title",
						"annotations.0.value": "alert title new",
						"annotations.1.key":   "desc",
						"annotations.1.value": "alert desc new",
						"annotations.2.key":   "test_key",
						"annotations.2.value": "test value new",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"group_configuration": map[string]interface{}{
						"type":   "custom",
						"fields": []string{"a", "b", "c"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"group_configuration.type":     "custom",
						"group_configuration.fields.#": "3",
						"group_configuration.fields.0": "a",
						"group_configuration.fields.1": "b",
						"group_configuration.fields.2": "c",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"policy_configuration": map[string]interface{}{
						"alert_policy_id":  "sls.builtin.dynamic",
						"action_policy_id": "test_action_policy_new",
						"repeat_interval":  "3h",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"policy_configuration.alert_policy_id":  "sls.builtin.dynamic",
						"policy_configuration.action_policy_id": "test_action_policy_new",
						"policy_configuration.repeat_interval":  "3h",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"join_configuration": map[string]interface{}{
						"type":      "left_join",
						"condition": "$0.cnt == $1.cnt",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"join_configuration.type":      "left_join",
						"join_configuration.condition": "$0.cnt == $1.cnt",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"severity_configurations": []map[string]interface{}{
						{
							"severity": "8",
							"eval_condition": map[string]interface{}{
								"condition":       "cnt > 8",
								"count_condition": "__count__ > 4",
							},
						},
						{
							"severity": "6",
							"eval_condition": map[string]interface{}{
								"condition":       "cnt > 1",
								"count_condition": "__count__ > 1",
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"severity_configurations.#":                                "2",
						"severity_configurations.0.severity":                       "8",
						"severity_configurations.0.eval_condition.condition":       "cnt > 8",
						"severity_configurations.0.eval_condition.count_condition": "__count__ > 4",
						"severity_configurations.1.severity":                       "6",
						"severity_configurations.1.eval_condition.condition":       "cnt > 1",
						"severity_configurations.1.eval_condition.count_condition": "__count__ > 1",
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

var logAlertMap = map[string]string{
	"project_name": CHECKSET,
	"alert_name":   CHECKSET,
}

func resourceLogAlertDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alicloud_log_project" "default"{
	name = "${var.name}"
	description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  	project = "${alicloud_log_project.default.name}"
  	name = "${var.name}"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
`, name)
}
