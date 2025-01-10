// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudDataWorksDiJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudDataWorksDiJobCreate,
		Read:   resourceAliCloudDataWorksDiJobRead,
		Update: resourceAliCloudDataWorksDiJobUpdate,
		Delete: resourceAliCloudDataWorksDiJobDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_data_source_settings": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"destination_data_source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"di_job_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"job_settings": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ddl_handling_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"column_data_type_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"destination_data_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"source_data_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"runtime_settings": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"channel_settings": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cycle_schedule_settings": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cycle_migration_type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"schedule_parameters": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"migration_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"resource_settings": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_resource_settings": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_group_identifier": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"requested_cu": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
						"realtime_resource_settings": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_group_identifier": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"requested_cu": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
						"offline_resource_settings": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_group_identifier": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"requested_cu": {
										Type:     schema.TypeFloat,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"source_data_source_settings": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_source_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"data_source_properties": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timezone": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"encoding": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"source_data_source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"table_mappings": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transformation_rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_action_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"rule_target_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"rule_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"source_object_selection_rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"action": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"Include"}, false),
									},
									"expression": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"expression_type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"Exact"}, false),
									},
								},
							},
						},
					},
				},
			},
			"transformation_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_action_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_target_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudDataWorksDiJobCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateDIJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["ProjectId"] = d.Get("project_id")
	query["RegionId"] = client.RegionId
	query["ClientToken"] = buildClientToken(action)

	if v, ok := d.GetOk("destination_data_source_type"); ok {
		query["DestinationDataSourceType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		query["Description"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("source_data_source_type"); ok {
		query["SourceDataSourceType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("job_name"); ok {
		query["JobName"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("migration_type"); ok {
		query["MigrationType"] = StringPointer(v.(string))
	}

	if v, ok := d.GetOk("destination_data_source_settings"); ok {
		destinationDataSourceSettingsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["DataSourceName"] = dataLoopTmp["data_source_name"]
			destinationDataSourceSettingsMapsArray = append(destinationDataSourceSettingsMapsArray, dataLoopMap)
		}
		destinationDataSourceSettingsMapsJson, err := json.Marshal(destinationDataSourceSettingsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		query["DestinationDataSourceSettings"] = string(destinationDataSourceSettingsMapsJson)
	}

	if v, ok := d.GetOk("source_data_source_settings"); ok {
		sourceDataSourceSettingsMapsArray := make([]interface{}, 0)
		for _, dataLoop1 := range v.([]interface{}) {
			dataLoop1Tmp := dataLoop1.(map[string]interface{})
			dataLoop1Map := make(map[string]interface{})
			localData2 := make(map[string]interface{})
			encoding1, _ := jsonpath.Get("$[0].encoding", dataLoop1Tmp["data_source_properties"])
			if encoding1 != nil && encoding1 != "" {
				localData2["Encoding"] = encoding1
			}
			timezone1, _ := jsonpath.Get("$[0].timezone", dataLoop1Tmp["data_source_properties"])
			if timezone1 != nil && timezone1 != "" {
				localData2["Timezone"] = timezone1
			}
			dataLoop1Map["DataSourceProperties"] = localData2
			dataLoop1Map["DataSourceName"] = dataLoop1Tmp["data_source_name"]
			sourceDataSourceSettingsMapsArray = append(sourceDataSourceSettingsMapsArray, dataLoop1Map)
		}
		sourceDataSourceSettingsMapsJson, err := json.Marshal(sourceDataSourceSettingsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		query["SourceDataSourceSettings"] = string(sourceDataSourceSettingsMapsJson)
	}

	objectDataLocalMap := make(map[string]interface{})

	if v := d.Get("resource_settings"); v != nil {
		scheduleResourceSettings := make(map[string]interface{})
		requestedCu1, _ := jsonpath.Get("$[0].schedule_resource_settings[0].requested_cu", v)
		if requestedCu1 != nil && requestedCu1 != "" {
			scheduleResourceSettings["RequestedCu"] = requestedCu1
		}
		resourceGroupIdentifier1, _ := jsonpath.Get("$[0].schedule_resource_settings[0].resource_group_identifier", v)
		if resourceGroupIdentifier1 != nil && resourceGroupIdentifier1 != "" {
			scheduleResourceSettings["ResourceGroupIdentifier"] = resourceGroupIdentifier1
		}

		objectDataLocalMap["ScheduleResourceSettings"] = scheduleResourceSettings
		realtimeResourceSettings := make(map[string]interface{})
		requestedCu3, _ := jsonpath.Get("$[0].realtime_resource_settings[0].requested_cu", v)
		if requestedCu3 != nil && requestedCu3 != "" {
			realtimeResourceSettings["RequestedCu"] = requestedCu3
		}
		resourceGroupIdentifier3, _ := jsonpath.Get("$[0].realtime_resource_settings[0].resource_group_identifier", v)
		if resourceGroupIdentifier3 != nil && resourceGroupIdentifier3 != "" {
			realtimeResourceSettings["ResourceGroupIdentifier"] = resourceGroupIdentifier3
		}

		objectDataLocalMap["RealtimeResourceSettings"] = realtimeResourceSettings
		offlineResourceSettings := make(map[string]interface{})
		requestedCu5, _ := jsonpath.Get("$[0].offline_resource_settings[0].requested_cu", v)
		if requestedCu5 != nil && requestedCu5 != "" {
			offlineResourceSettings["RequestedCu"] = requestedCu5
		}
		resourceGroupIdentifier5, _ := jsonpath.Get("$[0].offline_resource_settings[0].resource_group_identifier", v)
		if resourceGroupIdentifier5 != nil && resourceGroupIdentifier5 != "" {
			offlineResourceSettings["ResourceGroupIdentifier"] = resourceGroupIdentifier5
		}

		objectDataLocalMap["OfflineResourceSettings"] = offlineResourceSettings

		objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
		if err != nil {
			return WrapError(err)
		}
		query["ResourceSettings"] = string(objectDataLocalMapJson)
	}

	if v, ok := d.GetOk("table_mappings"); ok {
		tableMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop3 := range v.([]interface{}) {
			dataLoop3Tmp := dataLoop3.(map[string]interface{})
			dataLoop3Map := make(map[string]interface{})
			localMaps1 := make([]interface{}, 0)
			localData4 := dataLoop3Tmp["transformation_rules"]
			for _, dataLoop4 := range localData4.([]interface{}) {
				dataLoop4Tmp := dataLoop4.(map[string]interface{})
				dataLoop4Map := make(map[string]interface{})
				dataLoop4Map["RuleActionType"] = dataLoop4Tmp["rule_action_type"]
				dataLoop4Map["RuleName"] = dataLoop4Tmp["rule_name"]
				dataLoop4Map["RuleTargetType"] = dataLoop4Tmp["rule_target_type"]
				localMaps1 = append(localMaps1, dataLoop4Map)
			}
			dataLoop3Map["TransformationRules"] = localMaps1
			localMaps2 := make([]interface{}, 0)
			localData5 := dataLoop3Tmp["source_object_selection_rules"]
			for _, dataLoop5 := range localData5.([]interface{}) {
				dataLoop5Tmp := dataLoop5.(map[string]interface{})
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["Expression"] = dataLoop5Tmp["expression"]
				dataLoop5Map["ExpressionType"] = dataLoop5Tmp["expression_type"]
				dataLoop5Map["ObjectType"] = dataLoop5Tmp["object_type"]
				dataLoop5Map["Action"] = dataLoop5Tmp["action"]
				localMaps2 = append(localMaps2, dataLoop5Map)
			}
			dataLoop3Map["SourceObjectSelectionRules"] = localMaps2
			tableMappingsMapsArray = append(tableMappingsMapsArray, dataLoop3Map)
		}
		tableMappingsMapsJson, err := json.Marshal(tableMappingsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		query["TableMappings"] = string(tableMappingsMapsJson)
	}

	if v, ok := d.GetOk("transformation_rules"); ok {
		transformationRulesMapsArray := make([]interface{}, 0)
		for _, dataLoop6 := range v.([]interface{}) {
			dataLoop6Tmp := dataLoop6.(map[string]interface{})
			dataLoop6Map := make(map[string]interface{})
			dataLoop6Map["RuleExpression"] = dataLoop6Tmp["rule_expression"]
			dataLoop6Map["RuleName"] = dataLoop6Tmp["rule_name"]
			dataLoop6Map["RuleActionType"] = dataLoop6Tmp["rule_action_type"]
			dataLoop6Map["RuleTargetType"] = dataLoop6Tmp["rule_target_type"]
			transformationRulesMapsArray = append(transformationRulesMapsArray, dataLoop6Map)
		}
		transformationRulesMapsJson, err := json.Marshal(transformationRulesMapsArray)
		if err != nil {
			return WrapError(err)
		}
		query["TransformationRules"] = string(transformationRulesMapsJson)
	}

	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("job_settings"); !IsNil(v) {
		if v, ok := d.GetOk("job_settings"); ok {
			localData7, err := jsonpath.Get("$[0].ddl_handling_settings", v)
			if err != nil {
				localData7 = make([]interface{}, 0)
			}
			localMaps3 := make([]interface{}, 0)
			for _, dataLoop7 := range localData7.([]interface{}) {
				dataLoop7Tmp := make(map[string]interface{})
				if dataLoop7 != nil {
					dataLoop7Tmp = dataLoop7.(map[string]interface{})
				}
				dataLoop7Map := make(map[string]interface{})
				dataLoop7Map["Action"] = dataLoop7Tmp["action"]
				dataLoop7Map["Type"] = dataLoop7Tmp["type"]
				localMaps3 = append(localMaps3, dataLoop7Map)
			}
			objectDataLocalMap1["DdlHandlingSettings"] = localMaps3
		}

		channelSettings1, _ := jsonpath.Get("$[0].channel_settings", v)
		if channelSettings1 != nil && channelSettings1 != "" {
			objectDataLocalMap1["ChannelSettings"] = channelSettings1
		}
		if v, ok := d.GetOk("job_settings"); ok {
			localData8, err := jsonpath.Get("$[0].column_data_type_settings", v)
			if err != nil {
				localData8 = make([]interface{}, 0)
			}
			localMaps4 := make([]interface{}, 0)
			for _, dataLoop8 := range localData8.([]interface{}) {
				dataLoop8Tmp := make(map[string]interface{})
				if dataLoop8 != nil {
					dataLoop8Tmp = dataLoop8.(map[string]interface{})
				}
				dataLoop8Map := make(map[string]interface{})
				dataLoop8Map["DestinationDataType"] = dataLoop8Tmp["destination_data_type"]
				dataLoop8Map["SourceDataType"] = dataLoop8Tmp["source_data_type"]
				localMaps4 = append(localMaps4, dataLoop8Map)
			}
			objectDataLocalMap1["ColumnDataTypeSettings"] = localMaps4
		}

		cycleScheduleSettings := make(map[string]interface{})
		cycleMigrationType1, _ := jsonpath.Get("$[0].cycle_schedule_settings[0].cycle_migration_type", v)
		if cycleMigrationType1 != nil && cycleMigrationType1 != "" {
			cycleScheduleSettings["CycleMigrationType"] = cycleMigrationType1
		}
		scheduleParameters1, _ := jsonpath.Get("$[0].cycle_schedule_settings[0].schedule_parameters", v)
		if scheduleParameters1 != nil && scheduleParameters1 != "" {
			cycleScheduleSettings["ScheduleParameters"] = scheduleParameters1
		}

		objectDataLocalMap1["CycleScheduleSettings"] = cycleScheduleSettings
		if v, ok := d.GetOk("job_settings"); ok {
			localData9, err := jsonpath.Get("$[0].runtime_settings", v)
			if err != nil {
				localData9 = make([]interface{}, 0)
			}
			localMaps5 := make([]interface{}, 0)
			for _, dataLoop9 := range localData9.([]interface{}) {
				dataLoop9Tmp := make(map[string]interface{})
				if dataLoop9 != nil {
					dataLoop9Tmp = dataLoop9.(map[string]interface{})
				}
				dataLoop9Map := make(map[string]interface{})
				dataLoop9Map["Name"] = dataLoop9Tmp["name"]
				dataLoop9Map["Value"] = dataLoop9Tmp["value"]
				localMaps5 = append(localMaps5, dataLoop9Map)
			}
			objectDataLocalMap1["RuntimeSettings"] = localMaps5
		}

		objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
		query["JobSettings"] = string(objectDataLocalMap1Json)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_data_works_di_job", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", query["ProjectId"], response["DIJobId"]))

	return resourceAliCloudDataWorksDiJobRead(d, meta)
}

func resourceAliCloudDataWorksDiJobRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dataWorksServiceV2 := DataWorksServiceV2{client}

	objectRaw, err := dataWorksServiceV2.DescribeDataWorksDiJob(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_data_works_di_job DescribeDataWorksDiJob Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	if objectRaw["Description"] != nil {
		d.Set("description", objectRaw["Description"])
	}
	if objectRaw["DestinationDataSourceType"] != nil {
		d.Set("destination_data_source_type", objectRaw["DestinationDataSourceType"])
	}
	if objectRaw["JobName"] != nil {
		d.Set("job_name", objectRaw["JobName"])
	}
	if objectRaw["MigrationType"] != nil {
		d.Set("migration_type", objectRaw["MigrationType"])
	}
	if objectRaw["SourceDataSourceType"] != nil {
		d.Set("source_data_source_type", objectRaw["SourceDataSourceType"])
	}
	if objectRaw["DIJobId"] != nil {
		d.Set("di_job_id", formatInt(objectRaw["DIJobId"]))
	}
	if objectRaw["ProjectId"] != nil {
		d.Set("project_id", objectRaw["ProjectId"])
	}

	destinationDataSourceSettings1Raw := objectRaw["DestinationDataSourceSettings"]
	destinationDataSourceSettingsMaps := make([]map[string]interface{}, 0)
	if destinationDataSourceSettings1Raw != nil {
		for _, destinationDataSourceSettingsChild1Raw := range destinationDataSourceSettings1Raw.([]interface{}) {
			destinationDataSourceSettingsMap := make(map[string]interface{})
			destinationDataSourceSettingsChild1Raw := destinationDataSourceSettingsChild1Raw.(map[string]interface{})
			destinationDataSourceSettingsMap["data_source_name"] = destinationDataSourceSettingsChild1Raw["DataSourceName"]

			destinationDataSourceSettingsMaps = append(destinationDataSourceSettingsMaps, destinationDataSourceSettingsMap)
		}
	}
	if objectRaw["DestinationDataSourceSettings"] != nil {
		if err := d.Set("destination_data_source_settings", destinationDataSourceSettingsMaps); err != nil {
			return err
		}
	}
	jobSettingsMaps := make([]map[string]interface{}, 0)
	jobSettingsMap := make(map[string]interface{})
	jobSettings1Raw := make(map[string]interface{})
	if objectRaw["JobSettings"] != nil {
		jobSettings1Raw = objectRaw["JobSettings"].(map[string]interface{})
	}
	if len(jobSettings1Raw) > 0 {
		jobSettingsMap["channel_settings"] = jobSettings1Raw["ChannelSettings"]

		columnDataTypeSettings1Raw := jobSettings1Raw["ColumnDataTypeSettings"]
		columnDataTypeSettingsMaps := make([]map[string]interface{}, 0)
		if columnDataTypeSettings1Raw != nil {
			for _, columnDataTypeSettingsChild1Raw := range columnDataTypeSettings1Raw.([]interface{}) {
				columnDataTypeSettingsMap := make(map[string]interface{})
				columnDataTypeSettingsChild1Raw := columnDataTypeSettingsChild1Raw.(map[string]interface{})
				columnDataTypeSettingsMap["destination_data_type"] = columnDataTypeSettingsChild1Raw["DestinationDataType"]
				columnDataTypeSettingsMap["source_data_type"] = columnDataTypeSettingsChild1Raw["SourceDataType"]

				columnDataTypeSettingsMaps = append(columnDataTypeSettingsMaps, columnDataTypeSettingsMap)
			}
		}
		jobSettingsMap["column_data_type_settings"] = columnDataTypeSettingsMaps
		cycleScheduleSettingsMaps := make([]map[string]interface{}, 0)
		cycleScheduleSettingsMap := make(map[string]interface{})
		cycleScheduleSettings1Raw := make(map[string]interface{})
		if jobSettings1Raw["CycleScheduleSettings"] != nil {
			cycleScheduleSettings1Raw = jobSettings1Raw["CycleScheduleSettings"].(map[string]interface{})
		}
		if len(cycleScheduleSettings1Raw) > 0 {
			cycleScheduleSettingsMap["cycle_migration_type"] = cycleScheduleSettings1Raw["CycleMigrationType"]
			cycleScheduleSettingsMap["schedule_parameters"] = cycleScheduleSettings1Raw["ScheduleParameters"]

			cycleScheduleSettingsMaps = append(cycleScheduleSettingsMaps, cycleScheduleSettingsMap)
		}
		jobSettingsMap["cycle_schedule_settings"] = cycleScheduleSettingsMaps
		ddlHandlingSettings1Raw := jobSettings1Raw["DdlHandlingSettings"]
		ddlHandlingSettingsMaps := make([]map[string]interface{}, 0)
		if ddlHandlingSettings1Raw != nil {
			for _, ddlHandlingSettingsChild1Raw := range ddlHandlingSettings1Raw.([]interface{}) {
				ddlHandlingSettingsMap := make(map[string]interface{})
				ddlHandlingSettingsChild1Raw := ddlHandlingSettingsChild1Raw.(map[string]interface{})
				ddlHandlingSettingsMap["action"] = ddlHandlingSettingsChild1Raw["Action"]
				ddlHandlingSettingsMap["type"] = ddlHandlingSettingsChild1Raw["Type"]

				ddlHandlingSettingsMaps = append(ddlHandlingSettingsMaps, ddlHandlingSettingsMap)
			}
		}
		jobSettingsMap["ddl_handling_settings"] = ddlHandlingSettingsMaps
		runtimeSettings1Raw := jobSettings1Raw["RuntimeSettings"]
		runtimeSettingsMaps := make([]map[string]interface{}, 0)
		if runtimeSettings1Raw != nil {
			for _, runtimeSettingsChild1Raw := range runtimeSettings1Raw.([]interface{}) {
				runtimeSettingsMap := make(map[string]interface{})
				runtimeSettingsChild1Raw := runtimeSettingsChild1Raw.(map[string]interface{})
				runtimeSettingsMap["name"] = runtimeSettingsChild1Raw["Name"]
				runtimeSettingsMap["value"] = runtimeSettingsChild1Raw["Value"]

				runtimeSettingsMaps = append(runtimeSettingsMaps, runtimeSettingsMap)
			}
		}
		jobSettingsMap["runtime_settings"] = runtimeSettingsMaps
		jobSettingsMaps = append(jobSettingsMaps, jobSettingsMap)
	}
	if objectRaw["JobSettings"] != nil {
		if err := d.Set("job_settings", jobSettingsMaps); err != nil {
			return err
		}
	}
	resourceSettingsMaps := make([]map[string]interface{}, 0)
	resourceSettingsMap := make(map[string]interface{})
	resourceSettings1Raw := make(map[string]interface{})
	if objectRaw["ResourceSettings"] != nil {
		resourceSettings1Raw = objectRaw["ResourceSettings"].(map[string]interface{})
	}
	if len(resourceSettings1Raw) > 0 {

		offlineResourceSettingsMaps := make([]map[string]interface{}, 0)
		offlineResourceSettingsMap := make(map[string]interface{})
		offlineResourceSettings1Raw := make(map[string]interface{})
		if resourceSettings1Raw["OfflineResourceSettings"] != nil {
			offlineResourceSettings1Raw = resourceSettings1Raw["OfflineResourceSettings"].(map[string]interface{})
		}
		if len(offlineResourceSettings1Raw) > 0 {
			offlineResourceSettingsMap["requested_cu"] = offlineResourceSettings1Raw["RequestedCu"]
			offlineResourceSettingsMap["resource_group_identifier"] = offlineResourceSettings1Raw["ResourceGroupIdentifier"]

			offlineResourceSettingsMaps = append(offlineResourceSettingsMaps, offlineResourceSettingsMap)
		}
		resourceSettingsMap["offline_resource_settings"] = offlineResourceSettingsMaps
		realtimeResourceSettingsMaps := make([]map[string]interface{}, 0)
		realtimeResourceSettingsMap := make(map[string]interface{})
		realtimeResourceSettings1Raw := make(map[string]interface{})
		if resourceSettings1Raw["RealtimeResourceSettings"] != nil {
			realtimeResourceSettings1Raw = resourceSettings1Raw["RealtimeResourceSettings"].(map[string]interface{})
		}
		if len(realtimeResourceSettings1Raw) > 0 {
			realtimeResourceSettingsMap["requested_cu"] = realtimeResourceSettings1Raw["RequestedCu"]
			realtimeResourceSettingsMap["resource_group_identifier"] = realtimeResourceSettings1Raw["ResourceGroupIdentifier"]

			realtimeResourceSettingsMaps = append(realtimeResourceSettingsMaps, realtimeResourceSettingsMap)
		}
		resourceSettingsMap["realtime_resource_settings"] = realtimeResourceSettingsMaps
		scheduleResourceSettingsMaps := make([]map[string]interface{}, 0)
		scheduleResourceSettingsMap := make(map[string]interface{})
		scheduleResourceSettings1Raw := make(map[string]interface{})
		if resourceSettings1Raw["ScheduleResourceSettings"] != nil {
			scheduleResourceSettings1Raw = resourceSettings1Raw["ScheduleResourceSettings"].(map[string]interface{})
		}
		if len(scheduleResourceSettings1Raw) > 0 {
			scheduleResourceSettingsMap["requested_cu"] = scheduleResourceSettings1Raw["RequestedCu"]
			scheduleResourceSettingsMap["resource_group_identifier"] = scheduleResourceSettings1Raw["ResourceGroupIdentifier"]

			scheduleResourceSettingsMaps = append(scheduleResourceSettingsMaps, scheduleResourceSettingsMap)
		}
		resourceSettingsMap["schedule_resource_settings"] = scheduleResourceSettingsMaps
		resourceSettingsMaps = append(resourceSettingsMaps, resourceSettingsMap)
	}
	if objectRaw["ResourceSettings"] != nil {
		if err := d.Set("resource_settings", resourceSettingsMaps); err != nil {
			return err
		}
	}
	sourceDataSourceSettings1Raw := objectRaw["SourceDataSourceSettings"]
	sourceDataSourceSettingsMaps := make([]map[string]interface{}, 0)
	if sourceDataSourceSettings1Raw != nil {
		for _, sourceDataSourceSettingsChild1Raw := range sourceDataSourceSettings1Raw.([]interface{}) {
			sourceDataSourceSettingsMap := make(map[string]interface{})
			sourceDataSourceSettingsChild1Raw := sourceDataSourceSettingsChild1Raw.(map[string]interface{})
			sourceDataSourceSettingsMap["data_source_name"] = sourceDataSourceSettingsChild1Raw["DataSourceName"]

			dataSourcePropertiesMaps := make([]map[string]interface{}, 0)
			dataSourcePropertiesMap := make(map[string]interface{})
			dataSourceProperties1Raw := make(map[string]interface{})
			if sourceDataSourceSettingsChild1Raw["DataSourceProperties"] != nil {
				dataSourceProperties1Raw = sourceDataSourceSettingsChild1Raw["DataSourceProperties"].(map[string]interface{})
			}
			if len(dataSourceProperties1Raw) > 0 {
				dataSourcePropertiesMap["encoding"] = dataSourceProperties1Raw["Encoding"]
				dataSourcePropertiesMap["timezone"] = dataSourceProperties1Raw["Timezone"]

				dataSourcePropertiesMaps = append(dataSourcePropertiesMaps, dataSourcePropertiesMap)
			}
			sourceDataSourceSettingsMap["data_source_properties"] = dataSourcePropertiesMaps
			sourceDataSourceSettingsMaps = append(sourceDataSourceSettingsMaps, sourceDataSourceSettingsMap)
		}
	}
	if objectRaw["SourceDataSourceSettings"] != nil {
		if err := d.Set("source_data_source_settings", sourceDataSourceSettingsMaps); err != nil {
			return err
		}
	}
	tableMappings1Raw := objectRaw["TableMappings"]
	tableMappingsMaps := make([]map[string]interface{}, 0)
	if tableMappings1Raw != nil {
		for _, tableMappingsChild1Raw := range tableMappings1Raw.([]interface{}) {
			tableMappingsMap := make(map[string]interface{})
			tableMappingsChild1Raw := tableMappingsChild1Raw.(map[string]interface{})

			sourceObjectSelectionRules1Raw := tableMappingsChild1Raw["SourceObjectSelectionRules"]
			sourceObjectSelectionRulesMaps := make([]map[string]interface{}, 0)
			if sourceObjectSelectionRules1Raw != nil {
				for _, sourceObjectSelectionRulesChild1Raw := range sourceObjectSelectionRules1Raw.([]interface{}) {
					sourceObjectSelectionRulesMap := make(map[string]interface{})
					sourceObjectSelectionRulesChild1Raw := sourceObjectSelectionRulesChild1Raw.(map[string]interface{})
					sourceObjectSelectionRulesMap["action"] = sourceObjectSelectionRulesChild1Raw["Action"]
					sourceObjectSelectionRulesMap["expression"] = sourceObjectSelectionRulesChild1Raw["Expression"]
					sourceObjectSelectionRulesMap["expression_type"] = sourceObjectSelectionRulesChild1Raw["ExpressionType"]
					sourceObjectSelectionRulesMap["object_type"] = sourceObjectSelectionRulesChild1Raw["ObjectType"]

					sourceObjectSelectionRulesMaps = append(sourceObjectSelectionRulesMaps, sourceObjectSelectionRulesMap)
				}
			}
			tableMappingsMap["source_object_selection_rules"] = sourceObjectSelectionRulesMaps
			transformationRules2Raw := tableMappingsChild1Raw["TransformationRules"]
			transformationRulesMaps := make([]map[string]interface{}, 0)
			if transformationRules2Raw != nil {
				for _, transformationRulesChild2Raw := range transformationRules2Raw.([]interface{}) {
					transformationRulesMap := make(map[string]interface{})
					transformationRulesChild2Raw := transformationRulesChild2Raw.(map[string]interface{})
					transformationRulesMap["rule_action_type"] = transformationRulesChild2Raw["RuleActionType"]
					transformationRulesMap["rule_name"] = transformationRulesChild2Raw["RuleName"]
					transformationRulesMap["rule_target_type"] = transformationRulesChild2Raw["RuleTargetType"]

					transformationRulesMaps = append(transformationRulesMaps, transformationRulesMap)
				}
			}
			tableMappingsMap["transformation_rules"] = transformationRulesMaps
			tableMappingsMaps = append(tableMappingsMaps, tableMappingsMap)
		}
	}
	if objectRaw["TableMappings"] != nil {
		if err := d.Set("table_mappings", tableMappingsMaps); err != nil {
			return err
		}
	}
	transformationRules3Raw := objectRaw["TransformationRules"]
	transformationRulesMaps := make([]map[string]interface{}, 0)
	if transformationRules3Raw != nil {
		for _, transformationRulesChild3Raw := range transformationRules3Raw.([]interface{}) {
			transformationRulesMap := make(map[string]interface{})
			transformationRulesChild3Raw := transformationRulesChild3Raw.(map[string]interface{})
			transformationRulesMap["rule_action_type"] = transformationRulesChild3Raw["RuleActionType"]
			transformationRulesMap["rule_expression"] = transformationRulesChild3Raw["RuleExpression"]
			transformationRulesMap["rule_name"] = transformationRulesChild3Raw["RuleName"]
			transformationRulesMap["rule_target_type"] = transformationRulesChild3Raw["RuleTargetType"]

			transformationRulesMaps = append(transformationRulesMaps, transformationRulesMap)
		}
	}
	if objectRaw["TransformationRules"] != nil {
		if err := d.Set("transformation_rules", transformationRulesMaps); err != nil {
			return err
		}
	}

	return nil
}

func resourceAliCloudDataWorksDiJobUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	parts := strings.Split(d.Id(), ":")
	action := "UpdateDIJob"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	query["DIJobId"] = parts[1]
	query["ProjectId"] = parts[0]
	query["RegionId"] = client.RegionId
	query["ClientToken"] = buildClientToken(action)
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			query["Description"] = StringPointer(v.(string))
		}
	}

	if d.HasChange("table_mappings") {
		update = true
	}
	if v, ok := d.GetOk("table_mappings"); ok || d.HasChange("table_mappings") {
		tableMappingsMapsArray := make([]interface{}, 0)
		for _, dataLoop := range v.([]interface{}) {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			localMaps := make([]interface{}, 0)
			localData1 := dataLoopTmp["source_object_selection_rules"]
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["Expression"] = dataLoop1Tmp["expression"]
				dataLoop1Map["ExpressionType"] = dataLoop1Tmp["expression_type"]
				dataLoop1Map["ObjectType"] = dataLoop1Tmp["object_type"]
				dataLoop1Map["Action"] = dataLoop1Tmp["action"]
				localMaps = append(localMaps, dataLoop1Map)
			}
			dataLoopMap["SourceObjectSelectionRules"] = localMaps
			localMaps1 := make([]interface{}, 0)
			localData2 := dataLoopTmp["transformation_rules"]
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := dataLoop2.(map[string]interface{})
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["RuleName"] = dataLoop2Tmp["rule_name"]
				dataLoop2Map["RuleActionType"] = dataLoop2Tmp["rule_action_type"]
				dataLoop2Map["RuleTargetType"] = dataLoop2Tmp["rule_target_type"]
				localMaps1 = append(localMaps1, dataLoop2Map)
			}
			dataLoopMap["TransformationRules"] = localMaps1
			tableMappingsMapsArray = append(tableMappingsMapsArray, dataLoopMap)
		}
		tableMappingsMapsJson, err := json.Marshal(tableMappingsMapsArray)
		if err != nil {
			return WrapError(err)
		}
		query["TableMappings"] = string(tableMappingsMapsJson)
	}

	if d.HasChange("job_settings") {
		update = true
		objectDataLocalMap := make(map[string]interface{})

		if v := d.Get("job_settings"); v != nil {
			channelSettings1, _ := jsonpath.Get("$[0].channel_settings", v)
			if channelSettings1 != nil && (d.HasChange("job_settings.0.channel_settings") || channelSettings1 != "") {
				objectDataLocalMap["ChannelSettings"] = channelSettings1
			}
			if v, ok := d.GetOk("job_settings"); ok {
				localData3, err := jsonpath.Get("$[0].column_data_type_settings", v)
				if err != nil {
					localData3 = make([]interface{}, 0)
				}
				localMaps2 := make([]interface{}, 0)
				for _, dataLoop3 := range localData3.([]interface{}) {
					dataLoop3Tmp := make(map[string]interface{})
					if dataLoop3 != nil {
						dataLoop3Tmp = dataLoop3.(map[string]interface{})
					}
					dataLoop3Map := make(map[string]interface{})
					dataLoop3Map["DestinationDataType"] = dataLoop3Tmp["destination_data_type"]
					dataLoop3Map["SourceDataType"] = dataLoop3Tmp["source_data_type"]
					localMaps2 = append(localMaps2, dataLoop3Map)
				}
				objectDataLocalMap["ColumnDataTypeSettings"] = localMaps2
			}

			cycleScheduleSettings := make(map[string]interface{})
			scheduleParameters1, _ := jsonpath.Get("$[0].cycle_schedule_settings[0].schedule_parameters", v)
			if scheduleParameters1 != nil && (d.HasChange("job_settings.0.cycle_schedule_settings.0.schedule_parameters") || scheduleParameters1 != "") {
				cycleScheduleSettings["ScheduleParameters"] = scheduleParameters1
			}

			objectDataLocalMap["CycleScheduleSettings"] = cycleScheduleSettings
			if v, ok := d.GetOk("job_settings"); ok {
				localData4, err := jsonpath.Get("$[0].ddl_handling_settings", v)
				if err != nil {
					localData4 = make([]interface{}, 0)
				}
				localMaps3 := make([]interface{}, 0)
				for _, dataLoop4 := range localData4.([]interface{}) {
					dataLoop4Tmp := make(map[string]interface{})
					if dataLoop4 != nil {
						dataLoop4Tmp = dataLoop4.(map[string]interface{})
					}
					dataLoop4Map := make(map[string]interface{})
					dataLoop4Map["Type"] = dataLoop4Tmp["type"]
					dataLoop4Map["Action"] = dataLoop4Tmp["action"]
					localMaps3 = append(localMaps3, dataLoop4Map)
				}
				objectDataLocalMap["DdlHandlingSettings"] = localMaps3
			}

			if v, ok := d.GetOk("job_settings"); ok {
				localData5, err := jsonpath.Get("$[0].runtime_settings", v)
				if err != nil {
					localData5 = make([]interface{}, 0)
				}
				localMaps4 := make([]interface{}, 0)
				for _, dataLoop5 := range localData5.([]interface{}) {
					dataLoop5Tmp := make(map[string]interface{})
					if dataLoop5 != nil {
						dataLoop5Tmp = dataLoop5.(map[string]interface{})
					}
					dataLoop5Map := make(map[string]interface{})
					dataLoop5Map["Name"] = dataLoop5Tmp["name"]
					dataLoop5Map["Value"] = dataLoop5Tmp["value"]
					localMaps4 = append(localMaps4, dataLoop5Map)
				}
				objectDataLocalMap["RuntimeSettings"] = localMaps4
			}

			objectDataLocalMapJson, err := json.Marshal(objectDataLocalMap)
			if err != nil {
				return WrapError(err)
			}
			query["JobSettings"] = string(objectDataLocalMapJson)
		}
	}

	if d.HasChange("transformation_rules") {
		update = true
		if v, ok := d.GetOk("transformation_rules"); ok || d.HasChange("transformation_rules") {
			transformationRulesMapsArray := make([]interface{}, 0)
			for _, dataLoop6 := range v.([]interface{}) {
				dataLoop6Tmp := dataLoop6.(map[string]interface{})
				dataLoop6Map := make(map[string]interface{})
				dataLoop6Map["RuleActionType"] = dataLoop6Tmp["rule_action_type"]
				dataLoop6Map["RuleExpression"] = dataLoop6Tmp["rule_expression"]
				dataLoop6Map["RuleName"] = dataLoop6Tmp["rule_name"]
				dataLoop6Map["RuleTargetType"] = dataLoop6Tmp["rule_target_type"]
				transformationRulesMapsArray = append(transformationRulesMapsArray, dataLoop6Map)
			}
			transformationRulesMapsJson, err := json.Marshal(transformationRulesMapsArray)
			if err != nil {
				return WrapError(err)
			}
			query["TransformationRules"] = string(transformationRulesMapsJson)
		}
	}

	if d.HasChange("resource_settings") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})

	if v := d.Get("resource_settings"); v != nil {
		scheduleResourceSettings := make(map[string]interface{})
		requestedCu1, _ := jsonpath.Get("$[0].schedule_resource_settings[0].requested_cu", v)
		if requestedCu1 != nil && (d.HasChange("resource_settings.0.schedule_resource_settings.0.requested_cu") || requestedCu1 != "") {
			scheduleResourceSettings["RequestedCu"] = requestedCu1
		}
		resourceGroupIdentifier1, _ := jsonpath.Get("$[0].schedule_resource_settings[0].resource_group_identifier", v)
		if resourceGroupIdentifier1 != nil && (d.HasChange("resource_settings.0.schedule_resource_settings.0.resource_group_identifier") || resourceGroupIdentifier1 != "") {
			scheduleResourceSettings["ResourceGroupIdentifier"] = resourceGroupIdentifier1
		}

		objectDataLocalMap1["ScheduleResourceSettings"] = scheduleResourceSettings
		realtimeResourceSettings := make(map[string]interface{})
		resourceGroupIdentifier3, _ := jsonpath.Get("$[0].realtime_resource_settings[0].resource_group_identifier", v)
		if resourceGroupIdentifier3 != nil && (d.HasChange("resource_settings.0.realtime_resource_settings.0.resource_group_identifier") || resourceGroupIdentifier3 != "") {
			realtimeResourceSettings["ResourceGroupIdentifier"] = resourceGroupIdentifier3
		}
		requestedCu3, _ := jsonpath.Get("$[0].realtime_resource_settings[0].requested_cu", v)
		if requestedCu3 != nil && (d.HasChange("resource_settings.0.realtime_resource_settings.0.requested_cu") || requestedCu3 != "") {
			realtimeResourceSettings["RequestedCu"] = requestedCu3
		}

		objectDataLocalMap1["RealtimeResourceSettings"] = realtimeResourceSettings
		offlineResourceSettings := make(map[string]interface{})
		resourceGroupIdentifier5, _ := jsonpath.Get("$[0].offline_resource_settings[0].resource_group_identifier", v)
		if resourceGroupIdentifier5 != nil && (d.HasChange("resource_settings.0.offline_resource_settings.0.resource_group_identifier") || resourceGroupIdentifier5 != "") {
			offlineResourceSettings["ResourceGroupIdentifier"] = resourceGroupIdentifier5
		}
		requestedCu5, _ := jsonpath.Get("$[0].offline_resource_settings[0].requested_cu", v)
		if requestedCu5 != nil && (d.HasChange("resource_settings.0.offline_resource_settings.0.requested_cu") || requestedCu5 != "") {
			offlineResourceSettings["RequestedCu"] = requestedCu5
		}

		objectDataLocalMap1["OfflineResourceSettings"] = offlineResourceSettings

		objectDataLocalMap1Json, err := json.Marshal(objectDataLocalMap1)
		if err != nil {
			return WrapError(err)
		}
		query["ResourceSettings"] = string(objectDataLocalMap1Json)
	}

	if update {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}

	return resourceAliCloudDataWorksDiJobRead(d, meta)
}

func resourceAliCloudDataWorksDiJobDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteDIJob"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request = make(map[string]interface{})
	query["DIJobId"] = parts[1]
	query["ProjectId"] = parts[0]
	query["RegionId"] = client.RegionId

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2024-05-18"), StringPointer("AK"), query, request, &runtime)

		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"500130"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
