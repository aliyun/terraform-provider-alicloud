// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test CloudMonitorService MetricAlarmRule. >>> Resource test cases, automatically generated.
// Case resourceCase_info_labels 12863
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12863(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12863)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12863)
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
					"contact_groups":       "云账号报警联系人",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":       "云账号报警联系人",
						"metric_alarm_rule_id": CHECKSET,
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_effective_interval": "00:00-23:59",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_effective_interval": "00:00-23:59",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"value": "testValue",
							"key":   "testKey",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email_subject": "test alarm",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email_subject": "test alarm",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "https://www.aliyun.com/webhook",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "https://www.aliyun.com/webhook",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"resources": "[{\\\"resource\\\":\\\"acs:ecs:cn-hangzhou:*:instance/*\\\"}]",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"resources": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12863 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12863(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_composite 12864
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12864(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12864)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12864)
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
					"contact_groups":       "云账号报警联系人",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":                "3",
							"level":                "CRITICAL",
							"expression_list_join": "&&",
							"expression_list": []map[string]interface{}{
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"statistics":          "Average",
									"threshold":           "80",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"statistics":          "Average",
									"threshold":           "85",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":       "云账号报警联系人",
						"metric_alarm_rule_id": CHECKSET,
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":                "3",
							"level":                "CRITICAL",
							"expression_list_join": "&&",
							"expression_list": []map[string]interface{}{
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "80",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "85",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12864 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12864(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_expression 12865
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12865(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12865)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12865)
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
					"contact_groups":       "云账号报警联系人",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":          "3",
							"expression_raw": "$Average > 80",
							"level":          "CRITICAL",
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":       "云账号报警联系人",
						"metric_alarm_rule_id": CHECKSET,
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":          "3",
							"expression_raw": "$Average > 80",
							"level":          "CRITICAL",
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12865 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12865(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_targets 12866
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12866(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12866)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12866)
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
					"contact_groups":       "云账号报警联系人",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"targets": []map[string]interface{}{
						{
							"level":     "WARN",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-1",
						},
						{
							"level":     "INFO",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-2",
						},
						{
							"level":     "CRITICAL",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message2",
							"target_id": "test-target-3",
						},
					},
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":       "云账号报警联系人",
						"metric_alarm_rule_id": CHECKSET,
						"targets.#":            "3",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"level":     "WARN",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-1",
						},
						{
							"level":     "INFO",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-2",
						},
						{
							"level":     "CRITICAL",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message2",
							"target_id": "test-target-3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"targets": []map[string]interface{}{
						{
							"level":       "WARN",
							"arn":         "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id":   "test-target-1",
							"json_params": "{\\\"key\\\":\\\"value\\\"}",
						},
						{
							"level":     "INFO",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-2",
						},
						{
							"level":     "CRITICAL",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message2",
							"target_id": "test-target-3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"targets.0.json_params": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "43200",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"targets": []map[string]interface{}{
						{
							"level":     "WARN",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-1",
						},
						{
							"level":     "INFO",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-2",
						},
						{
							"level":     "CRITICAL",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message2",
							"target_id": "test-target-3",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"targets.#":            "3",
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"targets"},
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12866 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12866(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_preconditions 12867
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12867(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12867)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12867)
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
					"contact_groups":       "云账号报警联系人",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":       "云账号报警联系人",
						"metric_alarm_rule_id": CHECKSET,
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12867 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12867(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_prometheus 12868
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12868(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12868)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12868)
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
					"contact_groups":       "云账号报警联系人",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"prometheus": []map[string]interface{}{
						{
							"annotations": []map[string]interface{}{
								{
									"value": "High request rate on {{ $labels.instance }}",
									"key":   "summary",
								},
								{
									"value": "Request rate is elevated",
									"key":   "description",
								},
								{
									"value": "critical",
									"key":   "severity",
								},
							},
							"prom_ql": "avg(rate(http_requests_total[5m])) by (instance) > 100",
							"times":   "3",
							"level":   "2",
						},
					},
					"namespace":          "acs_prometheus",
					"metric_name":        "IOPSUsageOfDN",
					"resources":          "",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":       "云账号报警联系人",
						"metric_alarm_rule_id": CHECKSET,
						"namespace":            "acs_prometheus",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            "",
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus": []map[string]interface{}{
						{
							"annotations": []map[string]interface{}{
								{
									"value": "High request rate on {{ $labels.instance }}",
									"key":   "summary",
								},
								{
									"value": "Request rate is elevated",
									"key":   "description",
								},
								{
									"value": "critical",
									"key":   "severity",
								},
							},
							"prom_ql": "avg(rate(http_requests_total[5m])) by (instance) > 100",
							"times":   "3",
							"level":   "2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"prometheus": []map[string]interface{}{
						{
							"annotations": []map[string]interface{}{
								{
									"value": "Request rate is elevated on {{ $labels.instance }}",
									"key":   "description",
								},
								{
									"value": "warning",
									"key":   "severity",
								},
							},
							"prom_ql": "avg(rate(http_requests_total[5m])) by (instance) > 50",
							"times":   "5",
							"level":   "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"prometheus": []map[string]interface{}{
						{
							"annotations": []map[string]interface{}{
								{
									"value": "High request rate on {{ $labels.instance }}",
									"key":   "summary",
								},
								{
									"value": "Request rate is elevated",
									"key":   "description",
								},
								{
									"value": "critical",
									"key":   "severity",
								},
							},
							"prom_ql": "avg(rate(http_requests_total[5m])) by (instance) > 100",
							"times":   "3",
							"level":   "2",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_prometheus",
					"metric_name":        "IOPSUsageOfDN",
					"resources":          "",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_prometheus",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            "",
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12868 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12868(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_info_init 12869
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12869(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12869)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12869)
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
					"contact_groups":        "云账号报警联系人",
					"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"labels": []map[string]interface{}{
						{
							"value": "prod",
							"key":   "env",
						},
						{
							"value": "platform",
							"key":   "team",
						},
						{
							"value": "hangzhou",
							"key":   "region",
						},
					},
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"email_subject":      "test-email-subject",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "3",
									"statistics":          "Average",
									"threshold":           "50",
								},
							},
						},
					},
					"webhook":   "https://www.aliyun.com/webhook",
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
					"interval":  "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":        "云账号报警联系人",
						"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"metric_alarm_rule_id":  CHECKSET,
						"labels.#":              "3",
						"namespace":             "acs_drds",
						"metric_name":           "IOPSUsageOfDN",
						"email_subject":         "test-email-subject",
						"webhook":               "https://www.aliyun.com/webhook",
						"resources":             CHECKSET,
						"rule_name":             CHECKSET,
						"interval":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"value": "prod",
							"key":   "env",
						},
						{
							"value": "platform",
							"key":   "team",
						},
						{
							"value": "hangzhou",
							"key":   "region",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email_subject": "test-email-subject",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email_subject": "test-email-subject",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "https://www.aliyun.com/webhook",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "https://www.aliyun.com/webhook",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"labels": []map[string]interface{}{
						{
							"value": "prod",
							"key":   "env",
						},
						{
							"value": "platform",
							"key":   "team",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"labels.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"email_subject": "updated-email-subject",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"email_subject": "updated-email-subject",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"webhook": "https://www.aliyun.com/webhook2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"webhook": "https://www.aliyun.com/webhook2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":                "true",
					"send_ok":               "true",
					"contact_groups":        "云账号报警联系人",
					"silence_time":          "86400",
					"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":                "60",
					"labels": []map[string]interface{}{
						{
							"value": "prod",
							"key":   "env",
						},
						{
							"value": "platform",
							"key":   "team",
						},
						{
							"value": "hangzhou",
							"key":   "region",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"email_subject":      "test-email-subject",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "3",
									"statistics":          "Average",
									"threshold":           "50",
								},
							},
						},
					},
					"webhook":   "https://www.aliyun.com/webhook",
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
					"interval":  "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                "true",
						"send_ok":               "true",
						"contact_groups":        "云账号报警联系人",
						"silence_time":          CHECKSET,
						"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"metric_alarm_rule_id":  CHECKSET,
						"period":                CHECKSET,
						"labels.#":              "3",
						"effective_interval":    "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":        "KEEP_LAST_STATE",
						"namespace":             "acs_drds",
						"metric_name":           "IOPSUsageOfDN",
						"email_subject":         "test-email-subject",
						"webhook":               "https://www.aliyun.com/webhook",
						"resources":             CHECKSET,
						"rule_name":             CHECKSET,
						"interval":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12869 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12869(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_warn_init 12870
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12870(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12870)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12870)
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
					"contact_groups":       "云账号报警联系人",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"contact_groups":       "云账号报警联系人",
						"metric_alarm_rule_id": CHECKSET,
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"silence_time": "86400",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"period": "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"period": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"no_data_policy": "KEEP_LAST_STATE",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"no_data_policy": "KEEP_LAST_STATE",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlicloudCloudMonitorServiceMetricAlarmRuleMap12870 = map[string]string{
	"source_type": CHECKSET,
}

func AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12870(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alicloud_account" "this" {}

`, name)
}

// Case resourceCase_info_labels 12863  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12863_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12863)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12863)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_composite 12864  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12864_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12864)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12864)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":                "3",
							"level":                "CRITICAL",
							"expression_list_join": "&&",
							"expression_list": []map[string]interface{}{
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "80",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "85",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_expression 12865  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12865_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12865)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12865)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":          "3",
							"expression_raw": "$Average > 80",
							"level":          "CRITICAL",
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_targets 12866  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12866_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12866)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12866)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"targets": []map[string]interface{}{
						{
							"level":     "WARN",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-1",
						},
						{
							"level":     "INFO",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-2",
						},
						{
							"level":     "CRITICAL",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message2",
							"target_id": "test-target-3",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"targets.#":            "3",
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"targets"},
			},
		},
	})
}

// Case resourceCase_preconditions 12867  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12867_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12867)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12867)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_prometheus 12868  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12868_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12868)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12868)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"prometheus": []map[string]interface{}{
						{
							"annotations": []map[string]interface{}{
								{
									"value": "High request rate on {{ $labels.instance }}",
									"key":   "summary",
								},
								{
									"value": "Request rate is elevated",
									"key":   "description",
								},
								{
									"value": "critical",
									"key":   "severity",
								},
							},
							"prom_ql": "avg(rate(http_requests_total[5m])) by (instance) > 100",
							"times":   "3",
							"level":   "2",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_prometheus",
					"metric_name":        "IOPSUsageOfDN",
					"resources":          "",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_prometheus",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            "",
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_info_init 12869  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12869_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12869)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12869)
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
					"status":                "true",
					"send_ok":               "true",
					"contact_groups":        "云账号报警联系人",
					"silence_time":          "86400",
					"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":                "60",
					"labels": []map[string]interface{}{
						{
							"value": "prod",
							"key":   "env",
						},
						{
							"value": "platform",
							"key":   "team",
						},
						{
							"value": "hangzhou",
							"key":   "region",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"email_subject":      "test-email-subject",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "3",
									"statistics":          "Average",
									"threshold":           "50",
								},
							},
						},
					},
					"webhook":   "https://www.aliyun.com/webhook",
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
					"interval":  "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                "true",
						"send_ok":               "true",
						"contact_groups":        "云账号报警联系人",
						"silence_time":          CHECKSET,
						"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"metric_alarm_rule_id":  CHECKSET,
						"period":                CHECKSET,
						"labels.#":              "3",
						"effective_interval":    "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":        "KEEP_LAST_STATE",
						"namespace":             "acs_drds",
						"metric_name":           "IOPSUsageOfDN",
						"email_subject":         "test-email-subject",
						"webhook":               "https://www.aliyun.com/webhook",
						"resources":             CHECKSET,
						"rule_name":             CHECKSET,
						"interval":              CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_warn_init 12870  twin
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12870_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12870)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12870)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_info_labels 12863  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12863_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12863)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12863)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
					"escalations": []map[string]interface{}{
						{
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "3",
									"statistics":          "Average",
									"threshold":           "50",
								},
							},
						},
					},
					"no_effective_interval": "00:00-23:59",
					"labels": []map[string]interface{}{
						{
							"value": "testValue",
							"key":   "testKey",
						},
					},
					"email_subject": "test alarm",
					"webhook":       "https://www.aliyun.com/webhook",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok":               "false",
						"no_effective_interval": "00:00-23:59",
						"labels.#":              "1",
						"email_subject":         "test alarm",
						"webhook":               "https://www.aliyun.com/webhook",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_composite 12864  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12864_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12864)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12864)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":                "3",
							"level":                "CRITICAL",
							"expression_list_join": "&&",
							"expression_list": []map[string]interface{}{
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "80",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "85",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
					"composite_expression": []map[string]interface{}{
						{
							"times":                "5",
							"level":                "WARN",
							"expression_list_join": "||",
							"expression_list": []map[string]interface{}{
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "60",
								},
								{
									"metric_name":         "IOPSUsageOfDN",
									"comparison_operator": "GreaterThanThreshold",
									"period":              "60",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_expression 12865  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12865_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12865)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12865)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"composite_expression": []map[string]interface{}{
						{
							"times":          "3",
							"expression_raw": "$Average > 80",
							"level":          "CRITICAL",
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
					"composite_expression": []map[string]interface{}{
						{
							"times":          "5",
							"expression_raw": "$Average > 90",
							"level":          "WARN",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_targets 12866  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12866_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12866)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12866)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"targets": []map[string]interface{}{
						{
							"level":     "WARN",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-1",
						},
						{
							"level":     "INFO",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message",
							"target_id": "test-target-2",
						},
						{
							"level":     "CRITICAL",
							"arn":       "acs:mns:cn-hangzhou:${data.alicloud_account.this.id}:/queues/test/message2",
							"target_id": "test-target-3",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"targets.#":            "3",
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok":      "false",
					"silence_time": "43200",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok":      "false",
						"silence_time": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"targets"},
			},
		},
	})
}

// Case resourceCase_preconditions 12867  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12867_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12867)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12867)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "3",
									"statistics":          "Average",
									"threshold":           "50",
								},
							},
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_prometheus 12868  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12868_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12868)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12868)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"prometheus": []map[string]interface{}{
						{
							"annotations": []map[string]interface{}{
								{
									"value": "High request rate on {{ $labels.instance }}",
									"key":   "summary",
								},
								{
									"value": "Request rate is elevated",
									"key":   "description",
								},
								{
									"value": "critical",
									"key":   "severity",
								},
							},
							"prom_ql": "avg(rate(http_requests_total[5m])) by (instance) > 100",
							"times":   "3",
							"level":   "2",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_prometheus",
					"metric_name":        "IOPSUsageOfDN",
					"resources":          "",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_prometheus",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            "",
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
					"prometheus": []map[string]interface{}{
						{
							"annotations": []map[string]interface{}{
								{
									"value": "Request rate is elevated on {{ $labels.instance }}",
									"key":   "description",
								},
								{
									"value": "warning",
									"key":   "severity",
								},
							},
							"prom_ql": "avg(rate(http_requests_total[5m])) by (instance) > 50",
							"times":   "5",
							"level":   "3",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_info_init 12869  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12869_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12869)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12869)
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
					"status":                "true",
					"send_ok":               "true",
					"contact_groups":        "云账号报警联系人",
					"silence_time":          "86400",
					"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":                "60",
					"labels": []map[string]interface{}{
						{
							"value": "prod",
							"key":   "env",
						},
						{
							"value": "platform",
							"key":   "team",
						},
						{
							"value": "hangzhou",
							"key":   "region",
						},
					},
					"effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":     "KEEP_LAST_STATE",
					"namespace":          "acs_drds",
					"metric_name":        "IOPSUsageOfDN",
					"email_subject":      "test-email-subject",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "3",
									"statistics":          "Average",
									"threshold":           "50",
								},
							},
						},
					},
					"webhook":   "https://www.aliyun.com/webhook",
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
					"interval":  "60",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":                "true",
						"send_ok":               "true",
						"contact_groups":        "云账号报警联系人",
						"silence_time":          CHECKSET,
						"no_effective_interval": "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"metric_alarm_rule_id":  CHECKSET,
						"period":                CHECKSET,
						"labels.#":              "3",
						"effective_interval":    "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":        "KEEP_LAST_STATE",
						"namespace":             "acs_drds",
						"metric_name":           "IOPSUsageOfDN",
						"email_subject":         "test-email-subject",
						"webhook":               "https://www.aliyun.com/webhook",
						"resources":             CHECKSET,
						"rule_name":             CHECKSET,
						"interval":              CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
					"labels": []map[string]interface{}{
						{
							"value": "prod",
							"key":   "env",
						},
						{
							"value": "platform",
							"key":   "team",
						},
					},
					"email_subject": "updated-email-subject",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "90",
								},
							},
							"info": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "60",
								},
							},
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "70",
								},
							},
						},
					},
					"webhook": "https://www.aliyun.com/webhook2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok":       "false",
						"labels.#":      "2",
						"email_subject": "updated-email-subject",
						"webhook":       "https://www.aliyun.com/webhook2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Case resourceCase_warn_init 12870  raw
func TestAccAliCloudCloudMonitorServiceMetricAlarmRule_basic12870_raw(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cloud_monitor_service_metric_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudCloudMonitorServiceMetricAlarmRuleMap12870)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CloudMonitorServiceServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCloudMonitorServiceMetricAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccloudmonitorservice%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCloudMonitorServiceMetricAlarmRuleBasicDependence12870)
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
					"status":               "true",
					"send_ok":              "true",
					"contact_groups":       "云账号报警联系人",
					"silence_time":         "86400",
					"metric_alarm_rule_id": "SystemDefault_acs_drds_IOPSUsageOfDo",
					"period":               "60",
					"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
					"no_data_policy":       "KEEP_LAST_STATE",
					"namespace":            "acs_drds",
					"metric_name":          "IOPSUsageOfDN",
					"escalations": []map[string]interface{}{
						{
							"warn": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
					"resources": "[{\\\"resource\\\":\\\"_ALL\\\"}]",
					"rule_name": "SystemDefault_acs_drds_IOPSUsageOfDN",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status":               "true",
						"send_ok":              "true",
						"contact_groups":       "云账号报警联系人",
						"silence_time":         CHECKSET,
						"metric_alarm_rule_id": CHECKSET,
						"period":               CHECKSET,
						"effective_interval":   "00:00-23:59 +0800 dayofweek 1,2,3,4,5,6,7",
						"no_data_policy":       "KEEP_LAST_STATE",
						"namespace":            "acs_drds",
						"metric_name":          "IOPSUsageOfDN",
						"resources":            CHECKSET,
						"rule_name":            CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"send_ok": "false",
					"escalations": []map[string]interface{}{
						{
							"critical": []map[string]interface{}{
								{
									"comparison_operator": "GreaterThanThreshold",
									"times":               "5",
									"statistics":          "Average",
									"threshold":           "80",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"send_ok": "false",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test CloudMonitorService MetricAlarmRule. <<< Resource test cases, automatically generated.
