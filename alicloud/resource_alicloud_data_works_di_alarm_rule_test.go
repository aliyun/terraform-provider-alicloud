package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks DiAlarmRule. >>> Resource test cases, automatically generated.
// Case 数据集成报警规则_TF验收_成都 8956
func TestAccAliCloudDataWorksDiAlarmRule_basic8956(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_di_alarm_rule.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksDiAlarmRuleMap8956)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksDiAlarmRule")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksDiAlarmRuleBasicDependence8956)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "Description",
					"trigger_conditions": []map[string]interface{}{
						{
							"ddl_report_tags": []string{
								"ALTERADDCOLUMN"},
							"threshold": "20",
							"duration":  "10",
							"severity":  "Warning",
						},
					},
					"metric_type": "DdlReport",
					"notification_settings": []map[string]interface{}{
						{
							"notification_channels": []map[string]interface{}{
								{
									"severity": "Warning",
									"channels": []string{
										"Ding"},
								},
							},
							"notification_receivers": []map[string]interface{}{
								{
									"receiver_type": "DingToken",
									"receiver_values": []string{
										"1107550004253538"},
								},
							},
							"inhibition_interval": "10",
						},
					},
					"di_job_id":          "${alicloud_data_works_di_job.defaultUW8inp.di_job_id}",
					"di_alarm_rule_name": name,
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "Description",
						"trigger_conditions.#": "1",
						"metric_type":          "DdlReport",
						"di_job_id":            CHECKSET,
						"di_alarm_rule_name":   name,
						"enabled":              "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "D1",
					"trigger_conditions": []map[string]interface{}{
						{
							"severity":  "Critical",
							"duration":  "20",
							"threshold": "5",
							"ddl_report_tags": []string{
								"ALTERDROPCOLUMN"},
						},
					},
					"notification_settings": []map[string]interface{}{
						{
							"notification_channels": []map[string]interface{}{
								{
									"severity": "Critical",
									"channels": []string{
										"Ding"},
								},
							},
							"notification_receivers": []map[string]interface{}{
								{
									"receiver_type": "DingToken",
									"receiver_values": []string{
										"208139295196247891"},
								},
							},
							"inhibition_interval": "440",
						},
					},
					"enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "D1",
						"trigger_conditions.#": "1",
						"enabled":              "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "D2",
					"trigger_conditions": []map[string]interface{}{
						{
							"severity":  "Critical",
							"duration":  "20",
							"threshold": "0",
						},
					},
					"metric_type": "Heartbeat",
					"notification_settings": []map[string]interface{}{
						{
							"notification_channels": []map[string]interface{}{
								{
									"severity": "Critical",
									"channels": []string{
										"Mail"},
								},
							},
							"notification_receivers": []map[string]interface{}{
								{
									"receiver_type": "AliyunUid",
									"receiver_values": []string{
										"208139295196247891"},
								},
							},
							"inhibition_interval": "100",
						},
					},
					"di_alarm_rule_name": name + "_update",
					"enabled":            "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":          "D2",
						"trigger_conditions.#": "1",
						"metric_type":          "Heartbeat",
						"di_alarm_rule_name":   name + "_update",
						"enabled":              "false",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"trigger_conditions": []map[string]interface{}{
						{
							"severity":  "Warning",
							"duration":  "10",
							"threshold": "0",
						},
						{
							"severity":  "Critical",
							"duration":  "20",
							"threshold": "0",
						},
					},
					"notification_settings": []map[string]interface{}{
						{
							"notification_channels": []map[string]interface{}{
								{
									"severity": "Warning",
									"channels": []string{
										"Mail"},
								},
								{
									"severity": "Critical",
									"channels": []string{
										"Mail", "Ding", "Feishu"},
								},
							},
							"notification_receivers": []map[string]interface{}{
								{
									"receiver_type": "AliyunUid",
									"receiver_values": []string{
										"208139295196247891"},
								},
								{
									"receiver_type": "DingToken",
									"receiver_values": []string{
										"208139295196247891"},
								},
								{
									"receiver_type": "FeishuToken",
									"receiver_values": []string{
										"208139295196247891"},
								},
							},
							"inhibition_interval": "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"trigger_conditions.#": "2",
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

var AlicloudDataWorksDiAlarmRuleMap8956 = map[string]string{
	"di_alarm_rule_id": CHECKSET,
}

func AlicloudDataWorksDiAlarmRuleBasicDependence8956(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaulteNv8bu" {
  project_name = var.name
  display_name = "qianpeng_alarm_test_project"
  description  = "qianpeng_alarm_test_project"
  pai_task_enabled = true
}

resource "alicloud_data_works_di_job" "defaultUW8inp" {
  description             = "xxxx"
  project_id              = alicloud_data_works_project.defaulteNv8bu.id
  job_name                = "xxx"
  migration_type          = "api_xxx"
  source_data_source_type = "xxx"
  resource_settings {
    offline_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "xx"
    }
    realtime_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "xx"
    }
    schedule_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "xx"
    }
  }
  job_settings {
    channel_settings = "xxxx"
    column_data_type_settings {
      destination_data_type = "xxxx"
      source_data_type      = "xxxx"
    }
    cycle_schedule_settings {
      cycle_migration_type = "xxxx"
      schedule_parameters  = "xxxx"
    }
  }
  source_data_source_settings {
    data_source_name = "xxxx"
    data_source_properties {
      encoding = "xxxx"
      timezone = "xxxx"
    }
  }
  destination_data_source_type = "xxxx"
  table_mappings {
    source_object_selection_rules {
      action          = "Include"
      expression      = "xxxx"
      expression_type = "Exact"
      object_type     = "xxxx"
    }
    source_object_selection_rules {
      action          = "Include"
      expression      = "xxxx"
      expression_type = "Exact"
      object_type     = "xxxx"
    }
    transformation_rules {
      rule_name        = "xxxx"
      rule_action_type = "xxxx"
      rule_target_type = "xxxx"
    }
  }
  transformation_rules {
    rule_action_type = "xxxx"
    rule_expression  = "xxxx"
    rule_name        = "xxxx"
    rule_target_type = "xxxx"
  }
  destination_data_source_settings {
    data_source_name = "xxx"
  }
}


`, name)
}

// Test DataWorks DiAlarmRule. <<< Resource test cases, automatically generated.
