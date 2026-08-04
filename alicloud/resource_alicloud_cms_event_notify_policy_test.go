// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Cms EventNotifyPolicy. >>> Resource test cases, automatically generated.
// Case resource_EventNotifyPolicy_list_test 12980
func TestAccAliCloudCmsEventNotifyPolicy_basic12980(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12980)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12980)
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
					"description": "Event notification policy managed by Terraform",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Notify strategy for list test",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Event notification policy managed by Terraform",
						"workspace":   CHECKSET,
						"name":        name,
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

var AlicloudCmsEventNotifyPolicyMap12980 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12980(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_rca_template_test 12981
func TestAccAliCloudCmsEventNotifyPolicy_basic12981(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12981)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12981)
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
					"description": "Terraform test for EnableRca and CustomTemplateEntries fields",
					"response_plan": []map[string]interface{}{
						{
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Test notify strategy with RCA and custom templates",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"digital_employee_name": "test-employee-001",
									"enable_rca":            "true",
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
							"custom_template_entries": []map[string]interface{}{
								{
									"template_uuid": "test-template-uuid-001",
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform test for EnableRca and CustomTemplateEntries fields",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated: EnableRca disabled, template changed",
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Test notify strategy with RCA and custom templates",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"digital_employee_name": "test-employee-002",
									"enable_rca":            "false",
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
							"custom_template_entries": []map[string]interface{}{
								{
									"template_uuid": "test-template-uuid-002",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated: EnableRca disabled, template changed",
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

var AlicloudCmsEventNotifyPolicyMap12981 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12981(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_response_plan 12982
func TestAccAliCloudCmsEventNotifyPolicy_basic12982(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12982)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12982)
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
					"description": "Terraform ResponsePlan field coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"action-restore-001"},
									"alert_action_ids": []string{
										"action-alert-001"},
								},
							},
							"escalation_id": []string{
								"esc-test-001"},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Minimal notify strategy for ResponsePlan test",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform ResponsePlan field coverage test",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated ResponsePlan field coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "acknowledged",
									"repeat_interval":    "60",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"action-restore-002"},
									"alert_action_ids": []string{
										"action-alert-002", "action-alert-003"},
								},
							},
							"escalation_id": []string{
								"esc-test-001", "esc-test-002"},
							"auto_recover_seconds": "1200",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated ResponsePlan field coverage test",
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

var AlicloudCmsEventNotifyPolicyMap12982 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12982(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_lifecycle_test 12983
func TestAccAliCloudCmsEventNotifyPolicy_basic12983(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12983)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12983)
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
					"description": "Event notification policy managed by Terraform",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Primary notify strategy for lifecycle testing",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Event notification policy managed by Terraform",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated description for lifecycle test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "60",
								},
							},
							"auto_recover_seconds": "1200",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Primary notify strategy for lifecycle testing",
							"ignore_restored_notification": "true",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "alertname"},
									"period_min":  "10",
									"times":       "2",
									"silence_sec": "600",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group-updated"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "false",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "OR",
									"conditions": []map[string]interface{}{
										{
											"field": "alertname",
											"op":    "CONTAIN",
											"value": "memory",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated description for lifecycle test",
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

var AlicloudCmsEventNotifyPolicyMap12983 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12983(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_enable_disable_test 12984
func TestAccAliCloudCmsEventNotifyPolicy_basic12984(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12984)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12984)
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
					"description": "Terraform Enable/Disable operations test for EventNotifyPolicy",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Notify strategy for enable-disable test",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "alertname"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform Enable/Disable operations test for EventNotifyPolicy",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated - test regular update operation",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "60",
								},
							},
							"auto_recover_seconds": "1200",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Notify strategy for enable-disable test",
							"ignore_restored_notification": "true",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "alertname"},
									"period_min":  "10",
									"times":       "2",
									"silence_sec": "600",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "false",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "WARNING",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated - test regular update operation",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enabled": "true",
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

var AlicloudCmsEventNotifyPolicyMap12984 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12984(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_non_private_attrs 12985
func TestAccAliCloudCmsEventNotifyPolicy_basic12985(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12985)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12985)
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
					"description": "Terraform non-private attrs coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Strategy for non-private attrs test",
							"ignore_restored_notification": "true",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"filter_setting": []map[string]interface{}{
								{
									"expression": "1 AND 2",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "source",
											"op":    "EQ",
											"value": "cms",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform non-private attrs coverage test",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Strategy for non-private attrs test updated",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Tokyo",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"filter_setting": []map[string]interface{}{
								{
									"expression": "1 OR 2",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "CONTAIN",
											"value": "WARNING",
										},
										{
											"field": "source",
											"op":    "CONTAIN",
											"value": "arms",
										},
									},
								},
							},
						},
					},
					"name": name + "update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "update",
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

var AlicloudCmsEventNotifyPolicyMap12985 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12985(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_enum_channel_types 12986
func TestAccAliCloudCmsEventNotifyPolicy_basic12986(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12986)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12986)
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
					"description": "Terraform enum ChannelType coverage test for EventNotifyPolicy",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Notify strategy for ChannelType enum coverage",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-enum-test-ding"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-enum-test-ding2"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform enum ChannelType coverage test for EventNotifyPolicy",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated ChannelType enum coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "60",
								},
							},
							"auto_recover_seconds": "1200",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Notify strategy for ChannelType enum coverage",
							"ignore_restored_notification": "true",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "alertname"},
									"period_min":  "10",
									"times":       "2",
									"silence_sec": "600",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-enum-test-ding-updated"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-enum-test-ding2-updated"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "false",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "memory",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated ChannelType enum coverage test",
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

var AlicloudCmsEventNotifyPolicyMap12986 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12986(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_filter_setting 12987
func TestAccAliCloudCmsEventNotifyPolicy_basic12987(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12987)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12987)
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
					"description": "Terraform FilterSetting coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Minimal strategy for filter testing",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"filter_setting": []map[string]interface{}{
										{
											"relation":   "AND",
											"expression": "1 AND 2",
											"conditions": []map[string]interface{}{
												{
													"field": "source",
													"op":    "EQ",
													"value": "cms",
												},
												{
													"field": "severity",
													"op":    "EQ",
													"value": "WARNING",
												},
											},
										},
									},
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"workspace_filter_setting": []map[string]interface{}{
								{
									"workspace_uuids": []string{
										"${alicloud_cms_workspace.default.id}"},
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"relation": "OR",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "cpu",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform FilterSetting coverage test",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Terraform FilterSetting coverage test - extreme conditions",
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Minimal strategy for filter testing",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"filter_setting": []map[string]interface{}{
										{
											"relation":   "OR",
											"expression": "1",
											"conditions": []map[string]interface{}{
												{
													"field": "severity",
													"op":    "CONTAIN",
													"value": "CRITICAL",
												},
											},
										},
									},
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-test-group"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"workspace_filter_setting": []map[string]interface{}{
								{
									"workspace_uuids": []string{
										"${alicloud_cms_workspace.default.id}"},
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "cpu_high",
										},
										{
											"field": "namespace",
											"op":    "EQ",
											"value": "acs_ecs",
										},
										{
											"field": "source",
											"op":    "EQ",
											"value": "cms",
										},
										{
											"field": "instance",
											"op":    "EQ",
											"value": "i-001",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform FilterSetting coverage test - extreme conditions",
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

var AlicloudCmsEventNotifyPolicyMap12987 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12987(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_array_reduce_empty 12988
func TestAccAliCloudCmsEventNotifyPolicy_basic12988(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12988)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12988)
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
					"description": "Terraform array reduce and empty coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"restore-action-test-001", "restore-action-test-002", "restore-action-test-003"},
									"alert_action_ids": []string{
										"alert-action-test-001", "alert-action-test-002", "alert-action-test-003"},
								},
							},
							"escalation_id": []string{
								"esc-id-test-001", "esc-id-test-002", "esc-id-test-003"},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Array reduce empty test strategy",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "instance", "alertname"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-reduce-group-1"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "480",
											"end_time_in_minute":   "1200",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-reduce-group-2"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "540",
											"end_time_in_minute":   "1080",
											"day_in_week": []string{
												"1", "2", "3"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-reduce-group-3"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"workspace_filter_setting": []map[string]interface{}{
								{
									"workspace_uuids": []string{
										"ws-uuid-test-001", "ws-uuid-test-002", "ws-uuid-test-003"},
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "cpu_high",
										},
										{
											"field": "namespace",
											"op":    "EQ",
											"value": "acs_ecs",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform array reduce and empty coverage test",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated array reduce and empty coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "60",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"restore-action-test-001"},
									"alert_action_ids": []string{
										"alert-action-test-001"},
								},
							},
							"escalation_id": []string{
								"esc-id-test-001"},
							"auto_recover_seconds": "1200",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Array reduce empty test strategy",
							"ignore_restored_notification": "true",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity"},
									"period_min":  "10",
									"times":       "2",
									"silence_sec": "600",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-reduce-group-1"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "480",
											"end_time_in_minute":   "1200",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-reduce-group-2"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "false",
							"workspace_filter_setting": []map[string]interface{}{
								{
									"workspace_uuids": []string{
										"ws-uuid-test-001", "ws-uuid-test-002"},
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "cpu_high",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated array reduce and empty coverage test",
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

var AlicloudCmsEventNotifyPolicyMap12988 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12988(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_array_empty_extended 12989
func TestAccAliCloudCmsEventNotifyPolicy_basic12989(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12989)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12989)
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
					"description": "Terraform extended array coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"restore-action-ext-001"},
									"alert_action_ids": []string{
										"alert-action-ext-001", "alert-action-ext-002", "alert-action-ext-003"},
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Extended array test strategy",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "alertname"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-ext-group-1", "tf-ext-group-2", "tf-ext-group-3"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "480",
											"end_time_in_minute":   "1200",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-ext-group-4", "tf-ext-group-5"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "540",
											"end_time_in_minute":   "1080",
											"day_in_week": []string{
												"1", "2", "3"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-ext-group-6", "tf-ext-group-7"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"workspace_filter_setting": []map[string]interface{}{
								{
									"workspace_uuids": []string{
										"ws-uuid-ext-001", "ws-uuid-ext-002", "ws-uuid-ext-003"},
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "cpu_high",
										},
										{
											"field": "namespace",
											"op":    "EQ",
											"value": "acs_ecs",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform extended array coverage test",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Updated extended array coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"restore-action-ext-001"},
									"alert_action_ids": []string{
										"alert-action-ext-001"},
								},
							},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Extended array test strategy",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "alertname"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-ext-group-1"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "480",
											"end_time_in_minute":   "1200",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-ext-group-4"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "540",
											"end_time_in_minute":   "1080",
											"day_in_week": []string{
												"1", "2", "3"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-ext-group-6", "tf-ext-group-7"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"workspace_filter_setting": []map[string]interface{}{
								{
									"workspace_uuids": []string{
										"ws-uuid-ext-001"},
								},
							},
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Updated extended array coverage test",
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

var AlicloudCmsEventNotifyPolicyMap12989 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12989(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Case resource_EventNotifyPolicy_array_multi_element 12990
func TestAccAliCloudCmsEventNotifyPolicy_basic12990(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_cms_event_notify_policy.default"
	ra := resourceAttrInit(resourceId, AlicloudCmsEventNotifyPolicyMap12990)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CmsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeCmsEventNotifyPolicy")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfacccms%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudCmsEventNotifyPolicyBasicDependence12990)
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
					"description": "Terraform array multi-element coverage test",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "30",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"tf-restore-action-1", "tf-restore-action-2"},
									"alert_action_ids": []string{
										"tf-alert-action-1", "tf-alert-action-2"},
								},
							},
							"escalation_id": []string{
								"tf-escalation-id-1", "tf-escalation-id-2"},
							"auto_recover_seconds": "600",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Multi-element array test strategy",
							"ignore_restored_notification": "false",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "alertname"},
									"period_min":  "5",
									"times":       "1",
									"silence_sec": "300",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2", "3", "4", "5"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-array-group-1", "tf-array-group-2", "tf-array-group-3"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
										{
											"receivers": []string{
												"tf-array-contact-1", "tf-array-contact-2", "tf-array-contact-3"},
											"channel_type": "CONTACT",
											// The server normalises the order of enabledSubChannels, so the
											// configuration lists them in the order it returns them.
											"enabled_sub_channels": []string{
												"EMAIL", "SMS"},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "480",
											"end_time_in_minute":   "1200",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-array-webhook-1", "tf-array-webhook-2", "tf-array-webhook-3"},
											"channel_type":         "WEBHOOK",
											"enabled_sub_channels": []string{},
										},
										{
											"receivers": []string{
												"tf-array-ding-1", "tf-array-ding-2", "tf-array-ding-3"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "true",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "cpu_high",
										},
									},
								},
							},
						},
					},
					"workspace": "${alicloud_cms_workspace.default.id}",
					"name":      name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Terraform array multi-element coverage test",
						"workspace":   CHECKSET,
						"name":        name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Merged - growth to 3 Routes, DayInWeek boundary [1,2], Conditions 4 elements",
					"response_plan": []map[string]interface{}{
						{
							"repeat_notify_setting": []map[string]interface{}{
								{
									"end_incident_state": "resolved",
									"repeat_interval":    "60",
								},
							},
							"pushing_setting": []map[string]interface{}{
								{
									"restore_action_ids": []string{
										"tf-restore-action-1", "tf-restore-action-2", "tf-restore-action-3"},
									"alert_action_ids": []string{
										"tf-alert-action-1", "tf-alert-action-2", "tf-alert-action-3"},
								},
							},
							"escalation_id": []string{
								"tf-escalation-id-1", "tf-escalation-id-2", "tf-escalation-id-3"},
							"auto_recover_seconds": "1200",
						},
					},
					"notify_strategy": []map[string]interface{}{
						{
							"description":                  "Multi-element array test strategy",
							"ignore_restored_notification": "true",
							"grouping_setting": []map[string]interface{}{
								{
									"grouping_keys": []string{
										"severity", "instance", "alertname"},
									"period_min":  "10",
									"times":       "2",
									"silence_sec": "600",
								},
							},
							"routes": []map[string]interface{}{
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "0",
											"end_time_in_minute":   "1439",
											"day_in_week": []string{
												"1", "2"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-array-group-1", "tf-array-group-2", "tf-array-group-3"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
										{
											"receivers": []string{
												"tf-array-contact-1", "tf-array-contact-2", "tf-array-contact-3"},
											"channel_type": "CONTACT",
											"enabled_sub_channels": []string{
												"EMAIL"},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "480",
											"end_time_in_minute":   "1200",
											"day_in_week": []string{
												"1", "2", "3", "4", "5", "6", "0"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-array-webhook-1", "tf-array-webhook-2", "tf-array-webhook-3"},
											"channel_type":         "WEBHOOK",
											"enabled_sub_channels": []string{},
										},
										{
											"receivers": []string{
												"tf-array-ding-1", "tf-array-ding-2", "tf-array-ding-3"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
									},
								},
								{
									"effect_time_range": []map[string]interface{}{
										{
											"time_zone":            "Asia/Shanghai",
											"start_time_in_minute": "540",
											"end_time_in_minute":   "1080",
											"day_in_week": []string{
												"0", "1", "2", "3", "4"},
										},
									},
									"channels": []map[string]interface{}{
										{
											"receivers": []string{
												"tf-array-ding-4", "tf-array-ding-5", "tf-array-ding-6"},
											"channel_type":         "DING",
											"enabled_sub_channels": []string{},
										},
										{
											"receivers": []string{
												"tf-array-webhook-4", "tf-array-webhook-5", "tf-array-webhook-6"},
											"channel_type":         "WEBHOOK",
											"enabled_sub_channels": []string{},
										},
									},
								},
							},
						},
					},
					"subscription": []map[string]interface{}{
						{
							"subscribe_legacy_event": "false",
							"filter_setting": []map[string]interface{}{
								{
									"relation": "AND",
									"conditions": []map[string]interface{}{
										{
											"field": "severity",
											"op":    "EQ",
											"value": "CRITICAL",
										},
										{
											"field": "alertname",
											"op":    "EQ",
											"value": "cpu_high",
										},
										{
											"field": "namespace",
											"op":    "EQ",
											"value": "acs_ecs",
										},
										{
											"field": "source",
											"op":    "EQ",
											"value": "cms",
										},
									},
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": "Merged - growth to 3 Routes, DayInWeek boundary [1,2], Conditions 4 elements",
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

var AlicloudCmsEventNotifyPolicyMap12990 = map[string]string{
	"create_time": CHECKSET,
	"version":     CHECKSET,
	"user_id":     CHECKSET,
	"update_time": CHECKSET,
}

func AlicloudCmsEventNotifyPolicyBasicDependence12990(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

provider "alicloud" {
  region = "cn-hangzhou"
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

// Test Cms EventNotifyPolicy. <<< Resource test cases, automatically generated.
