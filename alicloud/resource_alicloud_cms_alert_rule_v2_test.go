// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms AlertRuleV2. >>> Resource test cases, automatically generated.
// Case resource_AlertRuleV2_UMODEL_test 12929
func TestAccAliCloudCmsAlertRuleV2_basic12929(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12929)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12929)
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
					"content_template": "umodel test alert on $${metric}",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"metric":        "CPUUtilization",
							"label_filters": []map[string]interface{}{
								{
									"operator": "=",
									"value":    "web-server",
									"name":     "app",
								},
								{
									"operator": "=",
									"value":    "production",
									"name":     "env",
								},
							},
							"metric_set": "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-umodel-10",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": CHECKSET,
						"display_name":     "regression-umodel-10",
						"enabled":          "true",
						"workspace":        "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content_template": "umodel test alert updated on $${metric}",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-integration-001"},
							"enabled": "true",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"metric":        "CPUUtilization",
							"label_filters": []map[string]interface{}{
								{
									"operator": "=",
									"value":    "api-gateway",
									"name":     "app",
								},
								{
									"operator": "=",
									"value":    "staging",
									"name":     "env",
								},
								{
									"operator": "=",
									"value":    "cn-hangzhou",
									"name":     "region",
								},
							},
							"metric_set": "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-umodel-10-updated",
					"enabled":      "false",
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "120",
							"threshold":     "95",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": CHECKSET,
						"display_name":     "regression-umodel-10-updated",
						"enabled":          "false",
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

var AlicloudCmsAlertRuleV2Map12929 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12929(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_IntegrationConfig_test 12939
func TestAccAliCloudCmsAlertRuleV2_basic12939(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12939)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12939)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-id-001", "action-id-002", "action-id-003"},
							"enabled": "true",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 1",
						},
					},
					"display_name": "integration-config-roundtrip-test",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "integration-config-roundtrip-test",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-id-004"},
							"enabled": "false",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 2",
						},
					},
					"display_name": "integration-config-roundtrip-test-s1",
					"enabled":      "false",
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "integration-config-roundtrip-test-s1",
						"enabled":      "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-id-005", "action-id-006", "action-id-007"},
							"enabled": "true",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"display_name": "integration-config-roundtrip-test-s2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "integration-config-roundtrip-test-s2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-id-008", "action-id-009", "action-id-010"},
							"enabled": "true",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 1",
						},
					},
					"display_name": "integration-config-roundtrip-test-s3",
					"enabled":      "true",
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "integration-config-roundtrip-test-s3",
						"enabled":      "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{},
							"enabled": "true",
						},
					},
					"display_name": "integration-config-roundtrip-test-s4-empty",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "integration-config-roundtrip-test-s4-empty",
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

var AlicloudCmsAlertRuleV2Map12939 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12939(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_Labels_Verify_test 12938
func TestAccAliCloudCmsAlertRuleV2_basic12938(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12938)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12938)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 1",
						},
					},
					"annotations": map[string]interface{}{
						"\"summary\"":     "CPU usage high",
						"\"description\"": "Node CPU is above 90%",
					},
					"display_name": "cloudspec_labels_verify_test",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"labels": map[string]interface{}{
						"\"severity\"": "critical",
						"\"env\"":      "test",
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "cloudspec_labels_verify_test",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 2",
						},
					},
					"annotations": map[string]interface{}{
						"\"summary\"":     "CPU usage updated",
						"\"description\"": "Updated CPU alert",
						"\"runbook\"":     "https://example.com",
					},
					"display_name": "cloudspec_labels_verify_test_updated",
					"labels": map[string]interface{}{
						"\"severity\"": "warning",
						"\"team\"":     "dev",
						"\"env\"":      "staging",
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "cloudspec_labels_verify_test_updated",
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

var AlicloudCmsAlertRuleV2Map12938 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12938(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_APM_Aggregate_test 12937
func TestAccAliCloudCmsAlertRuleV2_basic12937(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12937)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12937)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "app_service_stats",
									"group_by": []string{
										"serviceName", "callType"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-agg-1", "arms-agg-2"},
						},
					},
					"display_name": "regression-apm-aggregate",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-aggregate",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "app_service_stats",
									"group_by": []string{
										"serviceName", "statusCode"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-agg-1", "arms-agg-2"},
						},
					},
					"display_name": "regression-apm-aggregate-max",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "LT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "MAX",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "200",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-aggregate-max",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "app_service_stats",
									"group_by": []string{
										"serviceName", "interfaceName"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-agg-1", "arms-agg-2"},
						},
					},
					"display_name": "regression-apm-aggregate-p50",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "P50",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "INFO",
									"threshold": "50",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-aggregate-p50",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "app_service_stats",
									"group_by": []string{
										"serviceName", "regionId"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-agg-1", "arms-agg-2"},
						},
					},
					"display_name": "regression-apm-aggregate-p75",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "P75",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "ERROR",
									"threshold": "75",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-aggregate-p75",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "app_service_stats",
									"group_by": []string{
										"serviceName", "callType"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-agg-1", "arms-agg-2"},
						},
					},
					"display_name": "regression-apm-aggregate-p90",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "LT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "P90",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "90",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-aggregate-p90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "app_service_stats",
									"group_by": []string{
										"serviceName", "statusCode"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-agg-1", "arms-agg-2"},
						},
					},
					"display_name": "regression-apm-aggregate-p99",
					"enabled":      "false",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "P99",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "99",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-aggregate-p99",
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

var AlicloudCmsAlertRuleV2Map12937 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12937(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_QueryConfig_Advanced_test 12936
func TestAccAliCloudCmsAlertRuleV2_basic12936(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12936)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12936)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"entity_filters": []map[string]interface{}{
								{
									"operator": "=",
									"field":    "regionId",
									"value":    "cn-hangzhou",
								},
								{
									"operator": "=",
									"field":    "instanceId",
									"value":    "i-bp1abc123def",
								},
							},
							"metric":     "CPUUtilization",
							"metric_set": "acs_ecs_dashboard",
							"entity_fields": []map[string]interface{}{
								{
									"field": "name",
									"value": "my-instance",
								},
							},
						},
					},
					"display_name": "regression-queryconfig-advanced",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-queryconfig-advanced",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"entity_filters": []map[string]interface{}{
								{
									"operator": "=",
									"field":    "regionId",
									"value":    "cn-shanghai",
								},
								{
									"operator": "=",
									"field":    "instanceId",
									"value":    "i-bp2ghi456jkl",
								},
								{
									"operator": "=",
									"field":    "status",
									"value":    "Running",
								},
							},
							"metric":     "memory_usedutilization",
							"metric_set": "acs_ecs_dashboard",
							"entity_fields": []map[string]interface{}{
								{
									"field": "name",
									"value": "my-ecs-instance",
								},
								{
									"field": "hostName",
									"value": "web-server-01",
								},
							},
						},
					},
					"display_name": "regression-queryconfig-advanced-updated",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test-updated"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
							"threshold":     "95",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-queryconfig-advanced-updated",
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

var AlicloudCmsAlertRuleV2Map12936 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12936(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_APM_CompareList_test 12935
func TestAccAliCloudCmsAlertRuleV2_basic12935(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12935)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12935)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "http_response_time",
									"group_by": []string{
										"serviceName"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
							},
							"service_id_list": []string{
								"arms-svc-compare-1", "arms-svc-compare-2"},
						},
					},
					"display_name": "regression-apm-comparelist",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"compare_list": []map[string]interface{}{
								{
									"operator":       "GT",
									"aggregate":      "AVG",
									"yoy_time_value": "1",
									"threshold":      "80",
									"yoy_time_unit":  "hour",
								},
							},
							"type":     "APM_COMPOSITE_CONDITION",
							"relation": "OR",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-comparelist",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-apm-comparelist-updated",
					"condition_config": []map[string]interface{}{
						{
							"compare_list": []map[string]interface{}{
								{
									"operator":       "GT",
									"aggregate":      "P99",
									"yoy_time_value": "24",
									"threshold":      "95",
									"yoy_time_unit":  "hour",
								},
							},
							"type":     "APM_COMPOSITE_CONDITION",
							"relation": "AND",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-comparelist-updated",
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

var AlicloudCmsAlertRuleV2Map12935 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12935(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_APM_Operator_test 12934
func TestAccAliCloudCmsAlertRuleV2_basic12934(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12934)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12934)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "app_service_stats",
									"group_by": []string{
										"serviceName"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "EQ",
									"value": "http",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-app-001", "arms-app-002"},
						},
					},
					"display_name": "regression-apm-operator",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "COUNT",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-operator",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"display_name": "regression-apm-operator-gte",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "MIN",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "800",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-operator-gte",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"display_name": "regression-apm-operator-lt",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "LT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "CONTINUES",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "INFO",
									"threshold": "10",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-operator-lt",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-apm-operator-lte",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "LT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "COUNT",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "20",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-operator-lte",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-apm-operator-eq",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "EQ",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "MIN",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "ERROR",
									"threshold": "0",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-operator-eq",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"display_name": "regression-apm-operator-ne",
					"enabled":      "false",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "NE",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "CONTINUES",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "999",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-operator-ne",
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

var AlicloudCmsAlertRuleV2Map12934 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12934(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_Prometheus_Basic_test 12933
func TestAccAliCloudCmsAlertRuleV2_basic12933(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12933)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12933)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 1",
						},
					},
					"display_name": "regression-basic-prometheus",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-basic-prometheus",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 2",
						},
					},
					"display_name": "regression-basic-prometheus-updated",
					"enabled":      "false",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test-updated"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-basic-prometheus-updated",
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

var AlicloudCmsAlertRuleV2Map12933 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12933(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_SimpleEscalation_test 12932
func TestAccAliCloudCmsAlertRuleV2_basic12932(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12932)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12932)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"metric":        "memory_usedutilization",
							"metric_set":    "acs_ecs_dashboard",
						},
					},
					"display_name": "regression-umodel-simple-escalation",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GE",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
							"threshold":     "20",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-umodel-simple-escalation",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-umodel-simple-escalation-updated",
					"condition_config": []map[string]interface{}{
						{
							"operator":      "LE",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "180",
							"threshold":     "90",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-umodel-simple-escalation-updated",
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

var AlicloudCmsAlertRuleV2Map12932 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12932(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_ExpressEscalation_test 12931
func TestAccAliCloudCmsAlertRuleV2_basic12931(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12931)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12931)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
							"region_id":   "cn-hangzhou",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.2",
						},
					},
					"display_name": "regression-prom-express-escalation",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-prom-express-escalation",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes < 0.1",
						},
					},
					"display_name": "regression-prom-express-escalation-updated",
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-prom-express-escalation-updated",
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

var AlicloudCmsAlertRuleV2Map12931 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12931(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_APM_Simple_test 12930
func TestAccAliCloudCmsAlertRuleV2_basic12930(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12930)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12930)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "responseTime",
									"group_by": []string{
										"serviceName", "endpoint"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type":  "EQ",
									"value": "500",
									"key":   "statusCode",
								},
								{
									"type":  "NE",
									"value": "/health",
									"key":   "endpoint",
								},
							},
							"service_id_list": []string{
								"svc-001", "svc-002"},
						},
					},
					"display_name": "regression-apm-simple-8",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "AVG",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "1000",
								},
								{
									"severity":  "WARNING",
									"threshold": "500",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-simple-8",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "p99Latency",
									"group_by": []string{
										"region", "host", "method"},
									"window_secs": "120",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type":  "EQ",
									"value": "POST",
									"key":   "method",
								},
							},
							"service_id_list": []string{
								"svc-001", "svc-003", "svc-004"},
						},
					},
					"display_name": "regression-apm-simple-8-modified",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "AVG",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "50",
								},
								{
									"severity":  "INFO",
									"threshold": "30",
								},
								{
									"severity":  "ERROR",
									"threshold": "80",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-simple-8-modified",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "p99Latency",
									"group_by": []string{
										"zone"},
									"window_secs": "120",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type":  "EQ",
									"value": "POST",
									"key":   "method",
								},
							},
							"service_id_list": []string{
								"svc-009"},
						},
					},
					"display_name": "regression-apm-simple-8-modified-2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-simple-8-modified-2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "p99Latency",
									"group_by": []string{
										"cluster", "pod"},
									"window_secs": "120",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type":  "EQ",
									"value": "POST",
									"key":   "method",
								},
							},
							"service_id_list": []string{
								"svc-010", "svc-011"},
						},
					},
					"display_name": "regression-apm-simple-8-modified-3",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-simple-8-modified-3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "p99Latency",
									"group_by":     []string{},
									"window_secs":  "120",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type":  "EQ",
									"value": "POST",
									"key":   "method",
								},
							},
							"service_id_list": []string{},
						},
					},
					"display_name": "regression-apm-simple-8-modified-4-empty",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-simple-8-modified-4-empty",
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

var AlicloudCmsAlertRuleV2Map12930 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12930(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_APM_YOY_test 12918
func TestAccAliCloudCmsAlertRuleV2_basic12918(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12918)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12918)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "http_status",
									"group_by": []string{
										"serviceName", "statusCode"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "NE",
									"value": "internal",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-svc-yoy-1", "arms-svc-yoy-2"},
						},
					},
					"display_name": "regression-apm-yoy",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":       "YOY_UP",
							"type":           "APM_SIMPLE_CONDITION",
							"aggregate":      "AVG",
							"yoy_time_value": "5",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "30",
								},
							},
							"yoy_time_unit": "minute",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-yoy",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "http_status",
									"group_by": []string{
										"serviceName", "statusCode"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "NE",
									"value": "internal",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-svc-yoy-3", "arms-svc-yoy-4", "arms-svc-yoy-5"},
						},
					},
					"display_name": "regression-apm-yoy-week",
					"condition_config": []map[string]interface{}{
						{
							"operator":       "YOY_DOWN",
							"type":           "APM_SIMPLE_CONDITION",
							"aggregate":      "SUM",
							"yoy_time_value": "2",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "20",
								},
							},
							"yoy_time_unit": "week",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-yoy-week",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "http_status",
									"group_by": []string{
										"serviceName", "statusCode"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
								{
									"type":  "NE",
									"value": "internal",
									"key":   "callType",
								},
							},
							"service_id_list": []string{
								"arms-svc-yoy-1", "arms-svc-yoy-2"},
						},
					},
					"display_name": "regression-apm-yoy-month",
					"enabled":      "false",
					"condition_config": []map[string]interface{}{
						{
							"operator":       "YOY_UP",
							"type":           "APM_SIMPLE_CONDITION",
							"aggregate":      "MAX",
							"yoy_time_value": "1",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "ERROR",
									"threshold": "50",
								},
							},
							"yoy_time_unit": "month",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-yoy-month",
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

var AlicloudCmsAlertRuleV2Map12918 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12918(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_Prometheus_Rich_test 12928
func TestAccAliCloudCmsAlertRuleV2_basic12928(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12928)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12928)
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
					"content_template": "告警：$${alertName} 触发",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
							"region_id":   "cn-hangzhou",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "rate(http_requests_total[5m]) > 100",
						},
					},
					"display_name": "regression-enriched-1",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5", "6", "7"},
							"active_end_time":   "23:59",
							"silence_time_secs": "86400",
							"active_start_time": "00:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"group-1", "group-2", "group-3"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": CHECKSET,
						"display_name":     "regression-enriched-1",
						"enabled":          "true",
						"workspace":        "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content_template": "更新告警：$${alertName}",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "rate(http_requests_total[5m]) > 200",
						},
					},
					"display_name": "regression-enriched-1-updated",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5"},
							"active_end_time":   "18:00",
							"silence_time_secs": "43200",
							"active_start_time": "09:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test-updated"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": CHECKSET,
						"display_name":     "regression-enriched-1-updated",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-enriched-1-m2",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5"},
							"active_end_time":   "18:00",
							"silence_time_secs": "43200",
							"active_start_time": "09:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"group-4", "group-5", "group-6"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-enriched-1-m2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-enriched-1-m3",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5"},
							"active_end_time":   "18:00",
							"silence_time_secs": "43200",
							"active_start_time": "09:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"group-4", "group-5"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-enriched-1-m3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-enriched-1-m4",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5"},
							"active_end_time":   "18:00",
							"silence_time_secs": "43200",
							"active_start_time": "09:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"group-7"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-enriched-1-m4",
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

var AlicloudCmsAlertRuleV2Map12928 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12928(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_lifecycle_test 12927
func TestAccAliCloudCmsAlertRuleV2_basic12927(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12927)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12927)
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
					"content_template": "alert triggered on $${metric}",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
						},
					},
					"display_name": "cspec-lifecycle-test-rule",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": CHECKSET,
						"display_name":     "cspec-lifecycle-test-rule",
						"enabled":          "true",
						"workspace":        "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content_template": "alert updated on $${metric}",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-integration-001"},
							"enabled": "true",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "true",
						},
					},
					"display_name": "cspec-lifecycle-test-rule-updated",
					"enabled":      "false",
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "120",
							"threshold":     "95",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": CHECKSET,
						"display_name":     "cspec-lifecycle-test-rule-updated",
						"enabled":          "false",
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

var AlicloudCmsAlertRuleV2Map12927 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12927(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_Prometheus_NotifyStrategies_test 12926
func TestAccAliCloudCmsAlertRuleV2_basic12926(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12926)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12926)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
							"region_id":   "cn-hangzhou",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 1",
						},
					},
					"display_name": "regression-notify-strategies",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "NOTIFY_POLICY",
							"notify_strategies": []string{
								"strategy-1"},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-strategies",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 2",
						},
					},
					"display_name": "regression-notify-strategies-m1",
					"notify_config": []map[string]interface{}{
						{
							"type": "NOTIFY_POLICY",
							"notify_strategies": []string{
								"strategy-id-003"},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-strategies-m1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-notify-strategies-m2",
					"notify_config": []map[string]interface{}{
						{
							"type": "NOTIFY_POLICY",
							"notify_strategies": []string{
								"strategy-4"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-strategies-m2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-notify-strategies-m3",
					"notify_config": []map[string]interface{}{
						{
							"type": "NOTIFY_POLICY",
							"notify_strategies": []string{
								"strategy-5"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-strategies-m3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-notify-strategies-m4",
					"notify_config": []map[string]interface{}{
						{
							"type":              "NOTIFY_POLICY",
							"notify_strategies": []string{},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-strategies-m4",
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

var AlicloudCmsAlertRuleV2Map12926 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12926(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case delete_minimal 12925
func TestAccAliCloudCmsAlertRuleV2_basic12925(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12925)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12925)
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
					"content_template": "minimal delete test",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
						},
					},
					"display_name": "cspec-delete-minimal-test",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": "minimal delete test",
						"display_name":     "cspec-delete-minimal-test",
						"enabled":          "true",
						"workspace":        "default-cms-1511928242963727-cn-hangzhou",
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

var AlicloudCmsAlertRuleV2Map12925 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12925(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_Prometheus_Array_test 12924
func TestAccAliCloudCmsAlertRuleV2_basic12924(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12924)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12924)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 1",
						},
					},
					"display_name": "regression-severity-array",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5", "6", "7"},
							"active_end_time":   "23:59",
							"active_start_time": "00:00",
							"silence_time_secs": "86400",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"group-1", "group-2", "group-3"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "INFO",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-severity-array",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up > 2",
						},
					},
					"display_name": "regression-severity-array-updated",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "5", "6", "7"},
							"active_end_time":   "20:00",
							"active_start_time": "08:00",
							"silence_time_secs": "43200",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"group-a", "group-b"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "ERROR",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-severity-array-updated",
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

var AlicloudCmsAlertRuleV2Map12924 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12924(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_NotifyChannels_test 12923
func TestAccAliCloudCmsAlertRuleV2_basic12923(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12923)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12923)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "cmsxx-prometheus-1",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "up == 0",
						},
					},
					"display_name": "regression-notify-channels",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "3", "5"},
							"active_end_time":   "23:59",
							"silence_time_secs": "86400",
							"active_start_time": "00:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"cms-alarm-group-1", "cms-alarm-group-2"},
								},
								{
									"type": "DINGTALK",
									"identifiers": []string{
										"dingtalk-webhook-url-1", "dingtalk-webhook-url-2"},
								},
								{
									"type": "WEBHOOK",
									"identifiers": []string{
										"https://webhook.example.com/alert", "https://webhook2.example.com/alert"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-channels",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-notify-channels-m1",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5", "6", "7"},
							"active_end_time":   "23:59",
							"silence_time_secs": "86400",
							"active_start_time": "00:00",
							"channels": []map[string]interface{}{
								{
									"type": "CONTACT",
									"identifiers": []string{
										"contact-1", "contact-2", "contact-3"},
								},
								{
									"type": "FEISHU",
									"identifiers": []string{
										"feishu-webhook-url-1", "feishu-webhook-url-2"},
								},
								{
									"type": "SLACK",
									"identifiers": []string{
										"slack-channel-url-1", "slack-channel-url-2"},
								},
								{
									"type": "WEIXIN",
									"identifiers": []string{
										"weixin-bot-key-1", "weixin-bot-key-2"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-channels-m1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"display_name": "regression-notify-channels-m2",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"2", "4", "6"},
							"active_end_time":   "20:00",
							"silence_time_secs": "43200",
							"active_start_time": "08:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"group-a", "group-b", "group-c"},
								},
								{
									"type": "DINGTALK",
									"identifiers": []string{
										"dt-url-1", "dt-url-2", "dt-url-3"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-channels-m2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"display_name": "regression-notify-channels-m3",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5"},
							"active_end_time":   "18:00",
							"silence_time_secs": "86400",
							"active_start_time": "09:00",
							"channels": []map[string]interface{}{
								{
									"type": "CONTACT",
									"identifiers": []string{
										"admin-1", "admin-2"},
								},
								{
									"type": "WEBHOOK",
									"identifiers": []string{
										"https://hook1.example.com", "https://hook2.example.com"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-channels-m3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-notify-channels-m4",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "2", "3", "4", "5", "6", "7"},
							"active_end_time":   "23:59",
							"silence_time_secs": "86400",
							"active_start_time": "00:00",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"ops-group-1", "ops-group-2"},
								},
								{
									"type": "DINGTALK",
									"identifiers": []string{
										"dt-1", "dt-2"},
								},
								{
									"type": "FEISHU",
									"identifiers": []string{
										"feishu-1", "feishu-2"},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-channels-m4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "90",
						},
					},
					"display_name": "regression-notify-channels-m5",
					"notify_config": []map[string]interface{}{
						{
							"utc_offset": "+08:00",
							"type":       "DIRECT_NOTIFY",
							"active_days": []string{
								"1", "3", "5", "7"},
							"active_end_time":   "22:00",
							"silence_time_secs": "86400",
							"active_start_time": "06:00",
							"channels": []map[string]interface{}{
								{
									"type": "CONTACT",
									"identifiers": []string{
										"dev-1", "dev-2", "dev-3"},
								},
								{
									"type": "SLACK",
									"identifiers": []string{
										"slack-1", "slack-2"},
								},
								{
									"type": "WEIXIN",
									"identifiers": []string{
										"wx-1", "wx-2", "wx-3"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "90",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-notify-channels-m5",
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

var AlicloudCmsAlertRuleV2Map12923 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12923(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_Prometheus_DataCheck_test 12922
func TestAccAliCloudCmsAlertRuleV2_basic12922(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12922)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12922)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":                       "PROMETHEUS_SINGLE_QUERY",
							"expr":                       "up > 1",
							"enable_data_complete_check": "true",
						},
					},
					"display_name": "regression-data-check",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-data-check",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "120",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type":                       "PROMETHEUS_SINGLE_QUERY",
							"expr":                       "up > 2",
							"enable_data_complete_check": "false",
						},
					},
					"display_name": "regression-data-check-updated",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test-updated"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-data-check-updated",
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

var AlicloudCmsAlertRuleV2Map12922 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12922(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_Prometheus_Condition_test 12921
func TestAccAliCloudCmsAlertRuleV2_basic12921(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12921)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12921)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type":        "PROMETHEUS",
							"instance_id": "prom-regression-test",
							"region_id":   "cn-hangzhou",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "rate(http_errors_total[5m]) / rate(http_requests_total[5m]) > 0.05",
						},
					},
					"display_name": "regression-prom-condition-rewrite",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "WARNING",
							"duration_secs": "120",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-prom-condition-rewrite",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"query_config": []map[string]interface{}{
						{
							"type": "PROMETHEUS_SINGLE_QUERY",
							"expr": "rate(http_errors_total[10m]) / rate(http_requests_total[10m]) > 0.1",
						},
					},
					"display_name": "regression-prom-condition-rewrite-updated",
					"condition_config": []map[string]interface{}{
						{
							"type":          "PROMETHEUS_SIMPLE_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "180",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-prom-condition-rewrite-updated",
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

var AlicloudCmsAlertRuleV2Map12921 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12921(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_EmptyArray_Repro_test 12920
func TestAccAliCloudCmsAlertRuleV2_basic12920(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12920)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12920)
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
					"content_template": "empty array repro test",
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "UMODEL",
						},
					},
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-001", "action-002"},
							"enabled": "true",
						},
					},
					"arms_integration_config": []map[string]interface{}{
						{
							"enabled": "false",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"entity_type":   "instance",
							"type":          "UMODEL_METRICSET_QUERY",
							"entity_domain": "ecs",
							"metric":        "CPUUtilization",
							"metric_set":    "acs_ecs_dashboard",
						},
					},
					"display_name": "empty-array-repro-init",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "90",
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": "empty array repro test",
						"display_name":     "empty-array-repro-init",
						"enabled":          "true",
						"workspace":        "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content_template": "empty array repro test modified",
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{
								"action-003"},
							"enabled": "true",
						},
					},
					"display_name": "empty-array-repro-modify-nonempty",
					"condition_config": []map[string]interface{}{
						{
							"operator":      "GT",
							"type":          "UMODEL_METRICSET_CONDITION",
							"severity":      "CRITICAL",
							"duration_secs": "60",
							"threshold":     "95",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": "empty array repro test modified",
						"display_name":     "empty-array-repro-modify-nonempty",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"content_template": "empty array repro test emptied",
					"action_integration_config": []map[string]interface{}{
						{
							"actions": []string{},
							"enabled": "true",
						},
					},
					"display_name": "empty-array-repro-modify-empty-array-repro",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"content_template": "empty array repro test emptied",
						"display_name":     "empty-array-repro-modify-empty-array-repro",
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

var AlicloudCmsAlertRuleV2Map12920 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12920(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Case resource_AlertRuleV2_APM_Composite_test 12919
func TestAccAliCloudCmsAlertRuleV2_basic12919(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_alert_rule_v2.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsAlertRuleV2Map12919)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsAlertRuleV2")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsAlertRuleV2BasicDependence12919)
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
					"schedule_config": []map[string]interface{}{
						{
							"type":          "FIXED",
							"interval_secs": "60",
						},
					},
					"datasource_config": []map[string]interface{}{
						{
							"type": "APM",
						},
					},
					"query_config": []map[string]interface{}{
						{
							"measure_list": []map[string]interface{}{
								{
									"measure_code": "error_rate",
									"group_by": []string{
										"serviceName"},
									"window_secs": "60",
								},
							},
							"type": "APM_MULTI_QUERY",
							"filter_list": []map[string]interface{}{
								{
									"type": "ALL",
									"key":  "interfaceName",
								},
							},
							"service_id_list": []string{
								"arms-svc-composite-1", "arms-svc-composite-2"},
						},
					},
					"display_name": "regression-apm-composite-rewrite",
					"enabled":      "true",
					"notify_config": []map[string]interface{}{
						{
							"type": "DIRECT_NOTIFY",
							"channels": []map[string]interface{}{
								{
									"type": "GROUP",
									"identifiers": []string{
										"regression-test"},
								},
							},
						},
					},
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "WARNING",
									"threshold": "5",
								},
							},
						},
					},
					"workspace": "default-cms-1511928242963727-cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-composite-rewrite",
						"enabled":      "true",
						"workspace":    "default-cms-1511928242963727-cn-hangzhou",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name": "regression-apm-composite-rewrite-updated",
					"condition_config": []map[string]interface{}{
						{
							"operator":  "GT",
							"type":      "APM_SIMPLE_CONDITION",
							"aggregate": "SUM",
							"threshold_list": []map[string]interface{}{
								{
									"severity":  "CRITICAL",
									"threshold": "10",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": "regression-apm-composite-rewrite-updated",
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

var AlicloudCmsAlertRuleV2Map12919 = map[string]string{
	"created_at":                    CHECKSET,
	"datasource_type":               CHECKSET,
	"status":                        CHECKSET,
	"observe_resource_global_scope": CHECKSET,
	"severity_levels":               CHECKSET,
	"updated_at":                    CHECKSET,
}

func AlicloudCmsAlertRuleV2BasicDependence12919(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


`, name)
}

// Test Cms AlertRuleV2. <<< Resource test cases, automatically generated.
