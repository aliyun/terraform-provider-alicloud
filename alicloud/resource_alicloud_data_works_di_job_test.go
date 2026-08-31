package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test DataWorks DiJob. >>> Resource test cases, automatically generated.
// Case diJob_coverage_TF验收_成都 8957
func TestAccAliCloudDataWorksDiJob_basic8957(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_di_job.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksDiJobMap8957)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksDiJob")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksDiJobBasicDependence8957)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			//testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-chengdu"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"description":    "aaaa",
					"project_id":     "${alicloud_data_works_project.defaultMMHL8U.id}",
					"job_name":       "zhenyuan_test_case",
					"migration_type": "api_FullAndRealtimeIncremental",
					"source_data_source_settings": []map[string]interface{}{
						{
							"data_source_name": "dw_mysql",
							"data_source_properties": []map[string]interface{}{
								{
									"encoding": "utf-8",
									"timezone": "Asia/Shanghai",
								},
							},
						},
					},
					"destination_data_source_type": "Hologres",
					"table_mappings": []map[string]interface{}{
						{
							"source_object_selection_rules": []map[string]interface{}{
								{
									"action":          "Include",
									"expression":      "dw_mysql",
									"expression_type": "Exact",
									"object_type":     "Datasource",
								},
								{
									"action":          "Include",
									"expression":      "test_db1",
									"expression_type": "Exact",
									"object_type":     "Database",
								},
								{
									"action":          "Include",
									"expression":      "lsc_test01",
									"expression_type": "Exact",
									"object_type":     "Table",
								},
							},
							"transformation_rules": []map[string]interface{}{
								{
									"rule_name":        "my_table_rename_rule",
									"rule_action_type": "Rename",
									"rule_target_type": "Table",
								},
							},
						},
					},
					"source_data_source_type": "MySQL",
					"resource_settings": []map[string]interface{}{
						{
							"offline_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "2",
									"resource_group_identifier": "S_res_group_524257424564736_1716799673667",
								},
							},
							"realtime_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "2",
									"resource_group_identifier": "S_res_group_524257424564736_1716799673667",
								},
							},
							"schedule_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "2",
									"resource_group_identifier": "S_res_group_524257424564736_1716799673667",
								},
							},
						},
					},
					"transformation_rules": []map[string]interface{}{
						{
							"rule_action_type": "Rename",
							"rule_expression":  "{\\\"expression\\\":\\\"table2\\\"}",
							"rule_name":        "my_table_rename_rule",
							"rule_target_type": "Table",
						},
					},
					"destination_data_source_settings": []map[string]interface{}{
						{
							"data_source_name": "dw_test_holo",
						},
					},
					"job_settings": []map[string]interface{}{
						{
							"column_data_type_settings": []map[string]interface{}{
								{
									"destination_data_type": "bigint",
									"source_data_type":      "longtext",
								},
							},
							"ddl_handling_settings": []map[string]interface{}{
								{
									"action": "Ignore",
									"type":   "CreateTable",
								},
							},
							"runtime_settings": []map[string]interface{}{
								{
									"name":  "runtime.realtime.concurrent",
									"value": "1",
								},
							},
							"channel_settings": "1",
							"cycle_schedule_settings": []map[string]interface{}{
								{
									"cycle_migration_type": "2",
									"schedule_parameters":  "3",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":                        "aaaa",
						"project_id":                         CHECKSET,
						"job_name":                           "zhenyuan_test_case",
						"migration_type":                     "api_FullAndRealtimeIncremental",
						"source_data_source_settings.#":      "1",
						"destination_data_source_type":       "Hologres",
						"table_mappings.#":                   "1",
						"source_data_source_type":            "MySQL",
						"transformation_rules.#":             "1",
						"destination_data_source_settings.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "bbbbb",
					"table_mappings": []map[string]interface{}{
						{
							"source_object_selection_rules": []map[string]interface{}{
								{
									"action":          "Include",
									"expression":      "3333",
									"expression_type": "Exact",
									"object_type":     "3333",
								},
								{
									"action":          "Include",
									"expression":      "3333",
									"expression_type": "Exact",
									"object_type":     "3333",
								},
								{
									"action":          "Include",
									"expression":      "lsc_test01",
									"expression_type": "Exact",
									"object_type":     "3333",
								},
							},
							"transformation_rules": []map[string]interface{}{
								{
									"rule_name":        "my_define_pk",
									"rule_action_type": "DefinePrimaryKey",
									"rule_target_type": "333",
								},
							},
						},
					},
					"resource_settings": []map[string]interface{}{
						{
							"offline_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "3",
									"resource_group_identifier": "S_res_group_524257424564736_1703732138445",
								},
							},
							"realtime_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "3",
									"resource_group_identifier": "S_res_group_524257424564736_1703732138445",
								},
							},
							"schedule_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "3",
									"resource_group_identifier": "S_res_group_524257424564736_1703732138445",
								},
							},
						},
					},
					"transformation_rules": []map[string]interface{}{
						{
							"rule_action_type": "3333",
							"rule_expression":  "{\\\"expression\\\":\\\"public\\\"}",
							"rule_name":        "3333",
							"rule_target_type": "3333",
						},
						{
							"rule_action_type": "Rename",
							"rule_expression":  "{\\\"expression\\\":\\\"table2\\\"}",
							"rule_name":        "my_table_rename_rule",
							"rule_target_type": "4444444",
						},
						{
							"rule_action_type": "DefinePrimaryKey",
							"rule_expression":  "{\\\"columns\\\":[\\\"id\\\"]}",
							"rule_name":        "my_define_pk",
							"rule_target_type": "555555",
						},
					},
					"job_settings": []map[string]interface{}{
						{
							"column_data_type_settings": []map[string]interface{}{
								{
									"destination_data_type": "text",
									"source_data_type":      "3333",
								},
							},
							"ddl_handling_settings": []map[string]interface{}{
								{
									"action": "33333",
									"type":   "AddColumn",
								},
								{
									"action": "Ignore",
									"type":   "ModifyColumn",
								},
								{
									"action": "Ignore",
									"type":   "DropColumn",
								},
							},
							"runtime_settings": []map[string]interface{}{
								{
									"name":  "3333",
									"value": "3333",
								},
								{
									"name":  "runtime.realtime.concurrent",
									"value": "2",
								},
								{
									"name":  "runtime.offline.enable.error.record",
									"value": "false",
								},
							},
							"channel_settings": "3333",
							"cycle_schedule_settings": []map[string]interface{}{
								{
									"schedule_parameters": "3333",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "bbbbb",
						"table_mappings.#":       "1",
						"transformation_rules.#": "3",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": "cccc",
					"table_mappings": []map[string]interface{}{
						{
							"source_object_selection_rules": []map[string]interface{}{
								{
									"action":          "Include",
									"expression":      "dw_mysql",
									"expression_type": "Exact",
									"object_type":     "Datasource",
								},
								{
									"action":          "Include",
									"expression":      "t2",
									"expression_type": "Exact",
									"object_type":     "Database",
								},
								{
									"action":          "Include",
									"expression":      "t3",
									"expression_type": "Exact",
									"object_type":     "Table",
								},
							},
							"transformation_rules": []map[string]interface{}{
								{
									"rule_name":        "my_define_pk",
									"rule_action_type": "DefinePrimaryKey",
									"rule_target_type": "Table",
								},
							},
						},
						{
							"source_object_selection_rules": []map[string]interface{}{
								{
									"action":          "Include",
									"expression":      "xxx",
									"expression_type": "Exact",
									"object_type":     "xxx",
								},
								{
									"action":          "Include",
									"expression":      "xxx",
									"expression_type": "Exact",
									"object_type":     "xxx",
								},
								{
									"action":          "Include",
									"expression":      "xxx",
									"expression_type": "Exact",
									"object_type":     "xxx",
								},
							},
							"transformation_rules": []map[string]interface{}{
								{
									"rule_name":        "my_define_pk",
									"rule_action_type": "DefinePrimaryKey",
									"rule_target_type": "Table",
								},
							},
						},
						{
							"source_object_selection_rules": []map[string]interface{}{
								{
									"action":          "Include",
									"expression":      "xxx",
									"expression_type": "Exact",
									"object_type":     "xxx",
								},
								{
									"action":          "Include",
									"expression":      "xxx",
									"expression_type": "Exact",
									"object_type":     "xxx",
								},
								{
									"action":          "Include",
									"expression":      "xxx",
									"expression_type": "Exact",
									"object_type":     "xxx",
								},
							},
							"transformation_rules": []map[string]interface{}{
								{
									"rule_name":        "my_define_pk",
									"rule_action_type": "DefinePrimaryKey",
									"rule_target_type": "Table",
								},
							},
						},
					},
					"resource_settings": []map[string]interface{}{
						{
							"offline_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "1",
									"resource_group_identifier": "S_res_group_524257424564736_1716799673667",
								},
							},
							"realtime_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "1",
									"resource_group_identifier": "S_res_group_524257424564736_1716799673667",
								},
							},
							"schedule_resource_settings": []map[string]interface{}{
								{
									"requested_cu":              "1",
									"resource_group_identifier": "S_res_group_524257424564736_1716799673667",
								},
							},
						},
					},
					"transformation_rules": []map[string]interface{}{
						{
							"rule_action_type": "DefinePrimaryKey",
							"rule_expression":  "{\\\"columns\\\":[\\\"id\\\"]}",
							"rule_name":        "my_define_pk",
							"rule_target_type": "Table",
						},
					},
					"job_settings": []map[string]interface{}{
						{
							"column_data_type_settings": []map[string]interface{}{
								{
									"destination_data_type": "text",
									"source_data_type":      "longtext",
								},
								{
									"destination_data_type": "timestamptz",
									"source_data_type":      "date",
								},
								{
									"destination_data_type": "boolean",
									"source_data_type":      "bool",
								},
							},
							"ddl_handling_settings": []map[string]interface{}{
								{
									"action": "222222",
									"type":   "AddColumn",
								},
								{
									"action": "22222",
									"type":   "DropColumn",
								},
							},
							"runtime_settings": []map[string]interface{}{
								{
									"name":  "runtime.offline.concurrent",
									"value": "2",
								},
								{
									"name":  "runtime.realtime.concurrent",
									"value": "5",
								},
							},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":            "cccc",
						"table_mappings.#":       "3",
						"transformation_rules.#": "1",
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

var AlicloudDataWorksDiJobMap8957 = map[string]string{
	"di_job_id": CHECKSET,
}

func AlicloudDataWorksDiJobBasicDependence8957(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alicloud_data_works_project" "defaultMMHL8U" {
  description  = "浅鹏测试"
  project_name = var.name
  display_name = "qianpeng_terraform_test1"
  pai_task_enabled = true
}


`, name)
}

// Test DataWorks DiJob. <<< Resource test cases, automatically generated.
