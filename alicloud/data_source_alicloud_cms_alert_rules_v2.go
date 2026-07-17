// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAliCloudCmsAlertRulesV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCmsAlertRulesV2Read,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"filter_datasource_type_eq": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_display_name_contains": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_display_name_not_contains": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_enabled_eq": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter_labels_all_of_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_labels_all_of_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_labels_any_of_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_labels_any_of_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_notify_strategy_id_eq": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_observe_resource_global_scope_eq": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"filter_observe_resource_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_observe_resource_list_contains": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_observe_resource_type_eq": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_partition_key_eq": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_severity_levels_contains": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_status_eq": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_uuid_eq": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter_uuid_in": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_integration_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"actions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"alert_rule_v2_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"arms_integration_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"condition_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"legacy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"legacy_raw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"prometheus": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prom_ql": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"times": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"no_data_policy": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"simple_escalation": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"escalations": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"comparison_operator": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"times": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"pre_condition": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"severity": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"statistics": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"threshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"period": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"yoy_time_unit": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"compare_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"aggregate": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"yoy_time_value": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
												"yoy_time_unit": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"escalation_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"relation": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"composite_escalation": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"relation": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"times": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"escalations": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_name": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"comparison_operator": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"pre_condition": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"period": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"statistics": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"threshold": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"severity": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"duration_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"aggregate": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"express_escalation": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"raw_expression": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"times": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"severity": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"yoy_time_value": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"threshold_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"severity": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"threshold": {
													Type:     schema.TypeFloat,
													Computed: true,
												},
											},
										},
									},
									"threshold": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
								},
							},
						},
						"content_template": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datasource_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"legacy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"legacy_raw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"product_category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"datasource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"notify_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"utc_offset": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"notify_strategies": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"active_days": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
									"active_end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"silence_time_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"active_start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"channels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"identifiers": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
						"notify_strategy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"observe_resource_global_scope": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"observe_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"partition_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"legacy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"prom_ql": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"legacy_raw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"entity_filters": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"field": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"dimensions": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeMap, Elem: &schema.Schema{Type: schema.TypeString}},
									},
									"filter_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"label_filters": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metric_set": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"entity_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"measure_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"measure_code": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"group_by": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"window_secs": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"entity_domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"relation_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metric": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"service_id_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"enable_data_complete_check": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"entity_fields": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"field": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"schedule_config": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"interval_secs": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"severity_levels": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudCmsAlertRulesV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	// QueryAlertRules
	action := fmt.Sprintf("/queryAlertRules")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	query["RegionId"] = StringPointer(client.RegionId)
	filter := make(map[string]interface{})
	if v := d.Get("filter_datasource_type_eq"); !IsNil(v) {
		datasourceType := make(map[string]interface{})
		datasourceType["eq"] = v
		if len(datasourceType) > 0 {
			filter["datasourceType"] = datasourceType
		}
	}
	displayName := make(map[string]interface{})
	if v := d.Get("filter_display_name_contains"); !IsNil(v) {
		displayName["contains"] = v
	}
	if v := d.Get("filter_display_name_not_contains"); !IsNil(v) {
		displayName["notContains"] = v
	}
	if len(displayName) > 0 {
		filter["displayName"] = displayName
	}
	if v, ok := d.GetOkExists("filter_enabled_eq"); ok {
		enabled := make(map[string]interface{})
		enabled["eq"] = v
		if len(enabled) > 0 {
			filter["enabled"] = enabled
		}
	}
	if v := d.Get("filter_notify_strategy_id_eq"); !IsNil(v) {
		notifyStrategyId := make(map[string]interface{})
		notifyStrategyId["eq"] = v
		if len(notifyStrategyId) > 0 {
			filter["notifyStrategyId"] = notifyStrategyId
		}
	}
	if v, ok := d.GetOkExists("filter_observe_resource_global_scope_eq"); ok {
		observeResourceGlobalScope := make(map[string]interface{})
		observeResourceGlobalScope["eq"] = v
		if len(observeResourceGlobalScope) > 0 {
			filter["observeResourceGlobalScope"] = observeResourceGlobalScope
		}
	}
	if v, ok := d.GetOk("filter_observe_resource_instance_id"); ok {
		filter["observeResourceInstanceId"] = v
	}
	if v := d.Get("filter_observe_resource_list_contains"); !IsNil(v) {
		observeResourceList := make(map[string]interface{})
		observeResourceList["contains"] = []interface{}{v}
		if len(observeResourceList) > 0 {
			filter["observeResourceList"] = observeResourceList
		}
	}
	if v := d.Get("filter_observe_resource_type_eq"); !IsNil(v) {
		observeResourceType := make(map[string]interface{})
		observeResourceType["eq"] = v
		if len(observeResourceType) > 0 {
			filter["observeResourceType"] = observeResourceType
		}
	}
	if v := d.Get("filter_partition_key_eq"); !IsNil(v) {
		partitionKey := make(map[string]interface{})
		partitionKey["eq"] = v
		if len(partitionKey) > 0 {
			filter["partitionKey"] = partitionKey
		}
	}
	if v := d.Get("filter_severity_levels_contains"); !IsNil(v) {
		severityLevels := make(map[string]interface{})
		severityLevels["contains"] = []interface{}{v}
		if len(severityLevels) > 0 {
			filter["severityLevels"] = severityLevels
		}
	}
	if v := d.Get("filter_status_eq"); !IsNil(v) {
		status := make(map[string]interface{})
		status["eq"] = v
		if len(status) > 0 {
			filter["status"] = status
		}
	}
	uuid := make(map[string]interface{})
	if v := d.Get("filter_uuid_eq"); !IsNil(v) {
		uuid["eq"] = v
	}
	if v, ok := d.GetOk("filter_uuid_in"); ok {
		uuidIn := make([]interface{}, 0)
		for _, item := range strings.Split(v.(string), ",") {
			uuidIn = append(uuidIn, strings.TrimSpace(item))
		}
		uuid["in"] = uuidIn
	}
	if len(uuid) > 0 {
		filter["uuid"] = uuid
	}
	labels := make(map[string]interface{})
	labelsAllOf := make(map[string]interface{})
	if v := d.Get("filter_labels_all_of_key"); !IsNil(v) {
		labelsAllOf["key"] = v
	}
	if v := d.Get("filter_labels_all_of_value"); !IsNil(v) {
		labelsAllOf["value"] = v
	}
	if len(labelsAllOf) > 0 {
		labels["allOf"] = []interface{}{labelsAllOf}
	}
	labelsAnyOf := make(map[string]interface{})
	if v := d.Get("filter_labels_any_of_key"); !IsNil(v) {
		labelsAnyOf["key"] = v
	}
	if v := d.Get("filter_labels_any_of_value"); !IsNil(v) {
		labelsAnyOf["value"] = v
	}
	if len(labelsAnyOf) > 0 {
		labels["anyOf"] = []interface{}{labelsAnyOf}
	}
	if len(labels) > 0 {
		filter["labels"] = labels
	}

	if len(idsMap) > 0 {
		if _, ok := filter["uuid"]; !ok {
			uuidIn := make([]interface{}, 0, len(idsMap))
			for id := range idsMap {
				uuidIn = append(uuidIn, id)
			}
			filter["uuid"] = map[string]interface{}{"in": uuidIn}
		}
	}
	if len(filter) > 0 {
		request["filter"] = filter
	}

	if v, ok := d.GetOk("workspace"); ok {
		request["workspace"] = v
	}
	if len(request) == 0 {
		request["maxResults"] = PageSizeLarge
	}
	body = request
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	query["maxResults"] = StringPointer(strconv.Itoa(PageSizeLarge))
	for {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
			response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)

			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		resp, _ := jsonpath.Get("$.data.alertRules[*]", response)

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["uuid"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if nextToken, ok := response["nextToken"].(string); ok && nextToken != "" {
			query["nextToken"] = StringPointer(nextToken)
		} else {
			break
		}
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["uuid"]

		mapping["annotations"] = objectRaw["annotations"]
		mapping["content_template"] = objectRaw["contentTemplate"]
		mapping["created_at"] = objectRaw["createdAt"]
		mapping["datasource_type"] = objectRaw["datasourceType"]
		mapping["display_name"] = objectRaw["displayName"]
		mapping["enabled"] = objectRaw["enabled"]
		mapping["labels"] = objectRaw["labels"]
		mapping["notify_strategy_id"] = objectRaw["notifyStrategyId"]
		mapping["observe_resource_global_scope"] = objectRaw["observeResourceGlobalScope"]
		mapping["observe_resource_type"] = objectRaw["observeResourceType"]
		mapping["partition_key"] = objectRaw["partitionKey"]
		mapping["severity_levels"] = objectRaw["severityLevels"]
		mapping["status"] = objectRaw["status"]
		mapping["updated_at"] = objectRaw["updatedAt"]
		mapping["workspace"] = objectRaw["workspace"]
		mapping["alert_rule_v2_id"] = objectRaw["uuid"]

		actionIntegrationConfigMaps := make([]map[string]interface{}, 0)
		actionIntegrationConfigMap := make(map[string]interface{})
		actionIntegrationConfigRaw := make(map[string]interface{})
		if objectRaw["actionIntegrationConfig"] != nil {
			actionIntegrationConfigRaw = objectRaw["actionIntegrationConfig"].(map[string]interface{})
		}
		if len(actionIntegrationConfigRaw) > 0 {
			actionIntegrationConfigMap["enabled"] = actionIntegrationConfigRaw["enabled"]

			actionsRaw := make([]interface{}, 0)
			if actionIntegrationConfigRaw["actions"] != nil {
				actionsRaw = convertToInterfaceArray(actionIntegrationConfigRaw["actions"])
			}

			actionIntegrationConfigMap["actions"] = actionsRaw
			actionIntegrationConfigMaps = append(actionIntegrationConfigMaps, actionIntegrationConfigMap)
		}
		mapping["action_integration_config"] = actionIntegrationConfigMaps
		armsIntegrationConfigMaps := make([]map[string]interface{}, 0)
		armsIntegrationConfigMap := make(map[string]interface{})

		armsIntegrationConfigRaw := make(map[string]interface{})
		if objectRaw["armsIntegrationConfig"] != nil {
			armsIntegrationConfigRaw = objectRaw["armsIntegrationConfig"].(map[string]interface{})
		}
		if len(armsIntegrationConfigRaw) > 0 {
			armsIntegrationConfigMap["enabled"] = armsIntegrationConfigRaw["enabled"]

			armsIntegrationConfigMaps = append(armsIntegrationConfigMaps, armsIntegrationConfigMap)
		}
		mapping["arms_integration_config"] = armsIntegrationConfigMaps
		conditionConfigMaps := make([]map[string]interface{}, 0)
		conditionConfigMap := make(map[string]interface{})

		conditionConfigRaw := make(map[string]interface{})
		if objectRaw["conditionConfig"] != nil {
			conditionConfigRaw = objectRaw["conditionConfig"].(map[string]interface{})
		}
		if len(conditionConfigRaw) > 0 {
			conditionConfigMap["aggregate"] = conditionConfigRaw["aggregate"]
			conditionConfigMap["duration_secs"] = conditionConfigRaw["durationSecs"]
			conditionConfigMap["escalation_type"] = conditionConfigRaw["escalationType"]
			conditionConfigMap["legacy_raw"] = conditionConfigRaw["legacyRaw"]
			conditionConfigMap["legacy_type"] = conditionConfigRaw["legacyType"]
			conditionConfigMap["no_data_policy"] = conditionConfigRaw["noDataPolicy"]
			conditionConfigMap["operator"] = conditionConfigRaw["operator"]
			conditionConfigMap["relation"] = conditionConfigRaw["relation"]
			conditionConfigMap["severity"] = conditionConfigRaw["severity"]
			conditionConfigMap["threshold"] = conditionConfigRaw["threshold"]
			conditionConfigMap["type"] = conditionConfigRaw["type"]
			conditionConfigMap["yoy_time_unit"] = conditionConfigRaw["yoyTimeUnit"]
			conditionConfigMap["yoy_time_value"] = conditionConfigRaw["yoyTimeValue"]

			compareListRaw := conditionConfigRaw["compareList"]
			compareListMaps := make([]map[string]interface{}, 0)
			if compareListRaw != nil {
				for _, compareListChildRaw := range convertToInterfaceArray(compareListRaw) {
					compareListMap := make(map[string]interface{})

					compareListChildRaw := compareListChildRaw.(map[string]interface{})
					compareListMap["aggregate"] = compareListChildRaw["aggregate"]
					compareListMap["operator"] = compareListChildRaw["operator"]
					compareListMap["threshold"] = compareListChildRaw["threshold"]
					compareListMap["yoy_time_unit"] = compareListChildRaw["yoyTimeUnit"]
					compareListMap["yoy_time_value"] = compareListChildRaw["yoyTimeValue"]

					compareListMaps = append(compareListMaps, compareListMap)
				}
			}
			conditionConfigMap["compare_list"] = compareListMaps
			compositeEscalationMaps := make([]map[string]interface{}, 0)
			compositeEscalationMap := make(map[string]interface{})

			compositeEscalationRaw := make(map[string]interface{})
			if conditionConfigRaw["compositeEscalation"] != nil {
				compositeEscalationRaw = conditionConfigRaw["compositeEscalation"].(map[string]interface{})
			}
			if len(compositeEscalationRaw) > 0 {
				compositeEscalationMap["relation"] = compositeEscalationRaw["relation"]
				compositeEscalationMap["severity"] = compositeEscalationRaw["severity"]
				compositeEscalationMap["times"] = compositeEscalationRaw["times"]

				escalationsRaw := compositeEscalationRaw["escalations"]
				escalationsMaps := make([]map[string]interface{}, 0)
				if escalationsRaw != nil {
					for _, escalationsChildRaw := range convertToInterfaceArray(escalationsRaw) {
						escalationsMap := make(map[string]interface{})

						escalationsChildRaw := escalationsChildRaw.(map[string]interface{})
						escalationsMap["comparison_operator"] = escalationsChildRaw["comparisonOperator"]
						escalationsMap["metric_name"] = escalationsChildRaw["metricName"]
						escalationsMap["period"] = escalationsChildRaw["period"]
						escalationsMap["pre_condition"] = escalationsChildRaw["preCondition"]
						escalationsMap["statistics"] = escalationsChildRaw["statistics"]
						escalationsMap["threshold"] = escalationsChildRaw["threshold"]

						escalationsMaps = append(escalationsMaps, escalationsMap)
					}
				}
				compositeEscalationMap["escalations"] = escalationsMaps
				compositeEscalationMaps = append(compositeEscalationMaps, compositeEscalationMap)
			}
			conditionConfigMap["composite_escalation"] = compositeEscalationMaps
			expressEscalationMaps := make([]map[string]interface{}, 0)
			expressEscalationMap := make(map[string]interface{})

			expressEscalationRaw := make(map[string]interface{})
			if conditionConfigRaw["expressEscalation"] != nil {
				expressEscalationRaw = conditionConfigRaw["expressEscalation"].(map[string]interface{})
			}
			if len(expressEscalationRaw) > 0 {
				expressEscalationMap["raw_expression"] = expressEscalationRaw["rawExpression"]
				expressEscalationMap["severity"] = expressEscalationRaw["severity"]
				expressEscalationMap["times"] = expressEscalationRaw["times"]

				expressEscalationMaps = append(expressEscalationMaps, expressEscalationMap)
			}
			conditionConfigMap["express_escalation"] = expressEscalationMaps
			prometheusMaps := make([]map[string]interface{}, 0)
			prometheusMap := make(map[string]interface{})

			prometheusRaw := make(map[string]interface{})
			if conditionConfigRaw["prometheus"] != nil {
				prometheusRaw = conditionConfigRaw["prometheus"].(map[string]interface{})
			}
			if len(prometheusRaw) > 0 {
				prometheusMap["prom_ql"] = prometheusRaw["promQl"]
				prometheusMap["severity"] = prometheusRaw["severity"]
				prometheusMap["times"] = prometheusRaw["times"]

				prometheusMaps = append(prometheusMaps, prometheusMap)
			}
			conditionConfigMap["prometheus"] = prometheusMaps
			simpleEscalationMaps := make([]map[string]interface{}, 0)
			simpleEscalationMap := make(map[string]interface{})

			simpleEscalationRaw := make(map[string]interface{})
			if conditionConfigRaw["simpleEscalation"] != nil {
				simpleEscalationRaw = conditionConfigRaw["simpleEscalation"].(map[string]interface{})
			}
			if len(simpleEscalationRaw) > 0 {
				simpleEscalationMap["metric_name"] = simpleEscalationRaw["metricName"]
				simpleEscalationMap["period"] = simpleEscalationRaw["period"]

				escalationsRaw := simpleEscalationRaw["escalations"]
				escalationsMaps := make([]map[string]interface{}, 0)
				if escalationsRaw != nil {
					for _, escalationsChildRaw := range convertToInterfaceArray(escalationsRaw) {
						escalationsMap := make(map[string]interface{})

						escalationsChildRaw := escalationsChildRaw.(map[string]interface{})
						escalationsMap["comparison_operator"] = escalationsChildRaw["comparisonOperator"]
						escalationsMap["pre_condition"] = escalationsChildRaw["preCondition"]
						escalationsMap["severity"] = escalationsChildRaw["severity"]
						escalationsMap["statistics"] = escalationsChildRaw["statistics"]
						escalationsMap["threshold"] = escalationsChildRaw["threshold"]
						escalationsMap["times"] = escalationsChildRaw["times"]

						escalationsMaps = append(escalationsMaps, escalationsMap)
					}
				}
				simpleEscalationMap["escalations"] = escalationsMaps
				simpleEscalationMaps = append(simpleEscalationMaps, simpleEscalationMap)
			}
			conditionConfigMap["simple_escalation"] = simpleEscalationMaps

			thresholdListRaw := conditionConfigRaw["thresholdList"]
			thresholdListMaps := make([]map[string]interface{}, 0)
			if thresholdListRaw != nil {
				for _, thresholdListChildRaw := range convertToInterfaceArray(thresholdListRaw) {
					thresholdListMap := make(map[string]interface{})

					thresholdListChildRaw := thresholdListChildRaw.(map[string]interface{})
					thresholdListMap["severity"] = thresholdListChildRaw["severity"]
					thresholdListMap["threshold"] = thresholdListChildRaw["threshold"]

					thresholdListMaps = append(thresholdListMaps, thresholdListMap)
				}
			}
			conditionConfigMap["threshold_list"] = thresholdListMaps
			conditionConfigMaps = append(conditionConfigMaps, conditionConfigMap)
		}
		mapping["condition_config"] = conditionConfigMaps
		datasourceConfigMaps := make([]map[string]interface{}, 0)
		datasourceConfigMap := make(map[string]interface{})
		datasourceConfigRaw := make(map[string]interface{})
		if objectRaw["datasourceConfig"] != nil {
			datasourceConfigRaw = objectRaw["datasourceConfig"].(map[string]interface{})
		}
		if len(datasourceConfigRaw) > 0 {
			datasourceConfigMap["instance_id"] = datasourceConfigRaw["instanceId"]
			datasourceConfigMap["legacy_raw"] = datasourceConfigRaw["legacyRaw"]
			datasourceConfigMap["legacy_type"] = datasourceConfigRaw["legacyType"]
			datasourceConfigMap["product_category"] = datasourceConfigRaw["productCategory"]
			datasourceConfigMap["region_id"] = datasourceConfigRaw["regionId"]
			datasourceConfigMap["type"] = datasourceConfigRaw["type"]

			datasourceConfigMaps = append(datasourceConfigMaps, datasourceConfigMap)
		}
		mapping["datasource_config"] = datasourceConfigMaps
		notifyConfigMaps := make([]map[string]interface{}, 0)
		notifyConfigMap := make(map[string]interface{})
		notifyConfigRaw := make(map[string]interface{})
		if objectRaw["notifyConfig"] != nil {
			notifyConfigRaw = objectRaw["notifyConfig"].(map[string]interface{})
		}
		if len(notifyConfigRaw) > 0 {
			notifyConfigMap["active_end_time"] = notifyConfigRaw["activeEndTime"]
			notifyConfigMap["active_start_time"] = notifyConfigRaw["activeStartTime"]
			notifyConfigMap["silence_time_secs"] = notifyConfigRaw["silenceTimeSecs"]
			notifyConfigMap["type"] = notifyConfigRaw["type"]
			notifyConfigMap["utc_offset"] = notifyConfigRaw["utcOffset"]

			activeDaysRaw := make([]interface{}, 0)
			if notifyConfigRaw["activeDays"] != nil {
				activeDaysRaw = convertToInterfaceArray(notifyConfigRaw["activeDays"])
			}

			notifyConfigMap["active_days"] = activeDaysRaw
			channelsRaw := notifyConfigRaw["channels"]
			channelsMaps := make([]map[string]interface{}, 0)
			if channelsRaw != nil {
				for _, channelsChildRaw := range convertToInterfaceArray(channelsRaw) {
					channelsMap := make(map[string]interface{})
					channelsChildRaw := channelsChildRaw.(map[string]interface{})
					channelsMap["type"] = channelsChildRaw["type"]

					identifiersRaw := make([]interface{}, 0)
					if channelsChildRaw["identifiers"] != nil {
						identifiersRaw = convertToInterfaceArray(channelsChildRaw["identifiers"])
					}

					channelsMap["identifiers"] = identifiersRaw
					channelsMaps = append(channelsMaps, channelsMap)
				}
			}
			notifyConfigMap["channels"] = channelsMaps
			notifyStrategiesRaw := make([]interface{}, 0)
			if notifyConfigRaw["notifyStrategies"] != nil {
				notifyStrategiesRaw = convertToInterfaceArray(notifyConfigRaw["notifyStrategies"])
			}

			notifyConfigMap["notify_strategies"] = notifyStrategiesRaw
			notifyConfigMaps = append(notifyConfigMaps, notifyConfigMap)
		}
		mapping["notify_config"] = notifyConfigMaps
		queryConfigMaps := make([]map[string]interface{}, 0)
		queryConfigMap := make(map[string]interface{})
		queryConfigRaw := make(map[string]interface{})
		if objectRaw["queryConfig"] != nil {
			queryConfigRaw = objectRaw["queryConfig"].(map[string]interface{})
		}
		if len(queryConfigRaw) > 0 {
			queryConfigMap["enable_data_complete_check"] = queryConfigRaw["enableDataCompleteCheck"]
			queryConfigMap["entity_domain"] = queryConfigRaw["entityDomain"]
			queryConfigMap["entity_type"] = queryConfigRaw["entityType"]
			queryConfigMap["expr"] = queryConfigRaw["expr"]
			queryConfigMap["group_id"] = queryConfigRaw["groupId"]
			queryConfigMap["legacy_raw"] = queryConfigRaw["legacyRaw"]
			queryConfigMap["legacy_type"] = queryConfigRaw["legacyType"]
			queryConfigMap["metric"] = queryConfigRaw["metric"]
			queryConfigMap["metric_set"] = queryConfigRaw["metricSet"]
			queryConfigMap["namespace"] = queryConfigRaw["namespace"]
			queryConfigMap["prom_ql"] = queryConfigRaw["promQl"]
			queryConfigMap["relation_type"] = queryConfigRaw["relationType"]
			queryConfigMap["type"] = queryConfigRaw["type"]

			queryConfigMap["dimensions"] = queryConfigRaw["dimensions"]

			entityFieldsRaw := queryConfigRaw["entityFields"]
			entityFieldsMaps := make([]map[string]interface{}, 0)
			if entityFieldsRaw != nil {
				for _, entityFieldsChildRaw := range convertToInterfaceArray(entityFieldsRaw) {
					entityFieldsMap := make(map[string]interface{})

					entityFieldsChildRaw := entityFieldsChildRaw.(map[string]interface{})
					entityFieldsMap["field"] = entityFieldsChildRaw["field"]
					entityFieldsMap["value"] = entityFieldsChildRaw["value"]

					entityFieldsMaps = append(entityFieldsMaps, entityFieldsMap)
				}
			}
			queryConfigMap["entity_fields"] = entityFieldsMaps
			entityFiltersRaw := queryConfigRaw["entityFilters"]
			entityFiltersMaps := make([]map[string]interface{}, 0)
			if entityFiltersRaw != nil {
				for _, entityFiltersChildRaw := range convertToInterfaceArray(entityFiltersRaw) {
					entityFiltersMap := make(map[string]interface{})
					entityFiltersChildRaw := entityFiltersChildRaw.(map[string]interface{})
					entityFiltersMap["field"] = entityFiltersChildRaw["field"]
					entityFiltersMap["operator"] = entityFiltersChildRaw["operator"]
					entityFiltersMap["value"] = entityFiltersChildRaw["value"]

					entityFiltersMaps = append(entityFiltersMaps, entityFiltersMap)
				}
			}
			queryConfigMap["entity_filters"] = entityFiltersMaps

			filterListRaw := queryConfigRaw["filterList"]
			filterListMaps := make([]map[string]interface{}, 0)
			if filterListRaw != nil {
				for _, filterListChildRaw := range convertToInterfaceArray(filterListRaw) {
					filterListMap := make(map[string]interface{})

					filterListChildRaw := filterListChildRaw.(map[string]interface{})
					filterListMap["key"] = filterListChildRaw["key"]
					filterListMap["type"] = filterListChildRaw["type"]
					filterListMap["value"] = filterListChildRaw["value"]

					filterListMaps = append(filterListMaps, filterListMap)
				}
			}
			queryConfigMap["filter_list"] = filterListMaps

			labelFiltersRaw := queryConfigRaw["labelFilters"]
			labelFiltersMaps := make([]map[string]interface{}, 0)
			if labelFiltersRaw != nil {
				for _, labelFiltersChildRaw := range convertToInterfaceArray(labelFiltersRaw) {
					labelFiltersMap := make(map[string]interface{})

					labelFiltersChildRaw := labelFiltersChildRaw.(map[string]interface{})
					labelFiltersMap["name"] = labelFiltersChildRaw["name"]
					labelFiltersMap["operator"] = labelFiltersChildRaw["operator"]
					labelFiltersMap["value"] = labelFiltersChildRaw["value"]

					labelFiltersMaps = append(labelFiltersMaps, labelFiltersMap)
				}
			}
			queryConfigMap["label_filters"] = labelFiltersMaps

			measureListRaw := queryConfigRaw["measureList"]
			measureListMaps := make([]map[string]interface{}, 0)
			if measureListRaw != nil {
				for _, measureListChildRaw := range convertToInterfaceArray(measureListRaw) {
					measureListMap := make(map[string]interface{})

					measureListChildRaw := measureListChildRaw.(map[string]interface{})
					measureListMap["measure_code"] = measureListChildRaw["measureCode"]
					measureListMap["window_secs"] = measureListChildRaw["windowSecs"]

					groupByRaw := make([]interface{}, 0)
					if measureListChildRaw["groupBy"] != nil {
						groupByRaw = convertToInterfaceArray(measureListChildRaw["groupBy"])
					}

					measureListMap["group_by"] = groupByRaw
					measureListMaps = append(measureListMaps, measureListMap)
				}
			}
			queryConfigMap["measure_list"] = measureListMaps

			serviceIdListRaw := make([]interface{}, 0)
			if queryConfigRaw["serviceIdList"] != nil {
				serviceIdListRaw = convertToInterfaceArray(queryConfigRaw["serviceIdList"])
			}

			queryConfigMap["service_id_list"] = serviceIdListRaw
			queryConfigMaps = append(queryConfigMaps, queryConfigMap)
		}
		mapping["query_config"] = queryConfigMaps
		scheduleConfigMaps := make([]map[string]interface{}, 0)
		scheduleConfigMap := make(map[string]interface{})

		scheduleConfigRaw := make(map[string]interface{})
		if objectRaw["scheduleConfig"] != nil {
			scheduleConfigRaw = objectRaw["scheduleConfig"].(map[string]interface{})
		}
		if len(scheduleConfigRaw) > 0 {
			scheduleConfigMap["interval_secs"] = scheduleConfigRaw["intervalSecs"]
			scheduleConfigMap["type"] = scheduleConfigRaw["type"]

			scheduleConfigMaps = append(scheduleConfigMaps, scheduleConfigMap)
		}
		mapping["schedule_config"] = scheduleConfigMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
