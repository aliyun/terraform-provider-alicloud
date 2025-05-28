// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudSlsAlerts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudSlsAlertRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alerts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"configuration": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity_configurations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"severity": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"eval_condition": {
													Type:     schema.TypeList,
													Computed: true,
													MaxItems: 1,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"condition": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"count_condition": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"auto_annotation": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"no_data_fire": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"sink_cms": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"dashboard": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mute_until": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"template_configuration": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"annotations": {
													Type:     schema.TypeMap,
													Computed: true,
												},
												"version": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"lang": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"template_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"tokens": {
													Type:     schema.TypeMap,
													Computed: true,
												},
											},
										},
									},
									"labels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
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
									"group_configuration": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"fields": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"no_data_severity": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"annotations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
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
									"condition_configuration": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"count_condition": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"join_configurations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"policy_configuration": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"alert_policy_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"action_policy_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"repeat_interval": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"sink_event_store": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"project": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"event_store": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"endpoint": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"role_arn": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"query_list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"query": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"time_span_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"start": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"store": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"power_sql_mode": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"dashboard_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"role_arn": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"store_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"project": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"ui": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"end": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"chart_title": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"sink_alerthub": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"tags": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"send_resolved": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"threshold": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schedule": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"time_zone": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"run_immdiately": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"cron_expression": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"delay": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"interval": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
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
				ForceNew: true,
			},
		},
	}
}

func dataSourceAliCloudSlsAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		nameRegex = r
	}

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
	action := fmt.Sprintf("/alerts")
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("GET", "2020-12-30", "ListAlerts", action), query, nil, nil, hostMap, false)

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

	resp, _ := jsonpath.Get("$.results[*]", response)

	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(fmt.Sprint(item["name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(*hostMap["project"], ":", item["name"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = fmt.Sprint(*hostMap["project"], ":", objectRaw["name"])

		mapping["description"] = objectRaw["description"]
		mapping["display_name"] = objectRaw["displayName"]
		mapping["alert_name"] = objectRaw["name"]

		configurationMaps := make([]map[string]interface{}, 0)
		configurationMap := make(map[string]interface{})
		configurationRaw := make(map[string]interface{})
		if objectRaw["configuration"] != nil {
			configurationRaw = objectRaw["configuration"].(map[string]interface{})
		}
		if len(configurationRaw) > 0 {
			configurationMap["auto_annotation"] = configurationRaw["autoAnnotation"]
			configurationMap["dashboard"] = configurationRaw["dashboard"]
			configurationMap["mute_until"] = configurationRaw["muteUntil"]
			configurationMap["no_data_fire"] = configurationRaw["noDataFire"]
			configurationMap["no_data_severity"] = configurationRaw["noDataSeverity"]
			configurationMap["send_resolved"] = configurationRaw["sendResolved"]
			configurationMap["threshold"] = configurationRaw["threshold"]
			configurationMap["type"] = configurationRaw["type"]
			configurationMap["version"] = configurationRaw["version"]

			annotationsRaw := configurationRaw["annotations"]
			annotationsMaps := make([]map[string]interface{}, 0)
			if annotationsRaw != nil {
				for _, annotationsChildRaw := range annotationsRaw.([]interface{}) {
					annotationsMap := make(map[string]interface{})
					annotationsChildRaw := annotationsChildRaw.(map[string]interface{})
					annotationsMap["key"] = annotationsChildRaw["key"]
					annotationsMap["value"] = annotationsChildRaw["value"]

					annotationsMaps = append(annotationsMaps, annotationsMap)
				}
			}
			configurationMap["annotations"] = annotationsMaps
			conditionConfigurationMaps := make([]map[string]interface{}, 0)
			conditionConfigurationMap := make(map[string]interface{})
			conditionConfigurationRaw := make(map[string]interface{})
			if configurationRaw["conditionConfiguration"] != nil {
				conditionConfigurationRaw = configurationRaw["conditionConfiguration"].(map[string]interface{})
			}
			if len(conditionConfigurationRaw) > 0 {
				conditionConfigurationMap["condition"] = conditionConfigurationRaw["condition"]
				conditionConfigurationMap["count_condition"] = conditionConfigurationRaw["countCondition"]

				conditionConfigurationMaps = append(conditionConfigurationMaps, conditionConfigurationMap)
			}
			configurationMap["condition_configuration"] = conditionConfigurationMaps
			groupConfigurationMaps := make([]map[string]interface{}, 0)
			groupConfigurationMap := make(map[string]interface{})
			groupConfigurationRaw := make(map[string]interface{})
			if configurationRaw["groupConfiguration"] != nil {
				groupConfigurationRaw = configurationRaw["groupConfiguration"].(map[string]interface{})
			}
			if len(groupConfigurationRaw) > 0 {
				groupConfigurationMap["type"] = groupConfigurationRaw["type"]

				fieldsRaw := make([]interface{}, 0)
				if groupConfigurationRaw["fields"] != nil {
					fieldsRaw = groupConfigurationRaw["fields"].([]interface{})
				}

				groupConfigurationMap["fields"] = fieldsRaw
				groupConfigurationMaps = append(groupConfigurationMaps, groupConfigurationMap)
			}
			configurationMap["group_configuration"] = groupConfigurationMaps
			joinConfigurationsRaw := configurationRaw["joinConfigurations"]
			joinConfigurationsMaps := make([]map[string]interface{}, 0)
			if joinConfigurationsRaw != nil {
				for _, joinConfigurationsChildRaw := range joinConfigurationsRaw.([]interface{}) {
					joinConfigurationsMap := make(map[string]interface{})
					joinConfigurationsChildRaw := joinConfigurationsChildRaw.(map[string]interface{})
					joinConfigurationsMap["condition"] = joinConfigurationsChildRaw["condition"]
					joinConfigurationsMap["type"] = joinConfigurationsChildRaw["type"]

					joinConfigurationsMaps = append(joinConfigurationsMaps, joinConfigurationsMap)
				}
			}
			configurationMap["join_configurations"] = joinConfigurationsMaps
			labelsRaw := configurationRaw["labels"]
			labelsMaps := make([]map[string]interface{}, 0)
			if labelsRaw != nil {
				for _, labelsChildRaw := range labelsRaw.([]interface{}) {
					labelsMap := make(map[string]interface{})
					labelsChildRaw := labelsChildRaw.(map[string]interface{})
					labelsMap["key"] = labelsChildRaw["key"]
					labelsMap["value"] = labelsChildRaw["value"]

					labelsMaps = append(labelsMaps, labelsMap)
				}
			}
			configurationMap["labels"] = labelsMaps
			policyConfigurationMaps := make([]map[string]interface{}, 0)
			policyConfigurationMap := make(map[string]interface{})
			policyConfigurationRaw := make(map[string]interface{})
			if configurationRaw["policyConfiguration"] != nil {
				policyConfigurationRaw = configurationRaw["policyConfiguration"].(map[string]interface{})
			}
			if len(policyConfigurationRaw) > 0 {
				policyConfigurationMap["action_policy_id"] = policyConfigurationRaw["actionPolicyId"]
				policyConfigurationMap["alert_policy_id"] = policyConfigurationRaw["alertPolicyId"]
				policyConfigurationMap["repeat_interval"] = policyConfigurationRaw["repeatInterval"]

				policyConfigurationMaps = append(policyConfigurationMaps, policyConfigurationMap)
			}
			configurationMap["policy_configuration"] = policyConfigurationMaps
			queryListRaw := configurationRaw["queryList"]
			queryListMaps := make([]map[string]interface{}, 0)
			if queryListRaw != nil {
				for _, queryListChildRaw := range queryListRaw.([]interface{}) {
					queryListMap := make(map[string]interface{})
					queryListChildRaw := queryListChildRaw.(map[string]interface{})
					queryListMap["chart_title"] = queryListChildRaw["chartTitle"]
					queryListMap["dashboard_id"] = queryListChildRaw["dashboardId"]
					queryListMap["end"] = queryListChildRaw["end"]
					queryListMap["power_sql_mode"] = queryListChildRaw["powerSqlMode"]
					queryListMap["project"] = queryListChildRaw["project"]
					queryListMap["query"] = queryListChildRaw["query"]
					queryListMap["region"] = queryListChildRaw["region"]
					queryListMap["role_arn"] = queryListChildRaw["roleArn"]
					queryListMap["start"] = queryListChildRaw["start"]
					queryListMap["store"] = queryListChildRaw["store"]
					queryListMap["store_type"] = queryListChildRaw["storeType"]
					queryListMap["time_span_type"] = queryListChildRaw["timeSpanType"]
					queryListMap["ui"] = queryListChildRaw["ui"]

					queryListMaps = append(queryListMaps, queryListMap)
				}
			}
			configurationMap["query_list"] = queryListMaps
			severityConfigurationsRaw := configurationRaw["severityConfigurations"]
			severityConfigurationsMaps := make([]map[string]interface{}, 0)
			if severityConfigurationsRaw != nil {
				for _, severityConfigurationsChildRaw := range severityConfigurationsRaw.([]interface{}) {
					severityConfigurationsMap := make(map[string]interface{})
					severityConfigurationsChildRaw := severityConfigurationsChildRaw.(map[string]interface{})
					severityConfigurationsMap["severity"] = severityConfigurationsChildRaw["severity"]

					evalConditionMaps := make([]map[string]interface{}, 0)
					evalConditionMap := make(map[string]interface{})
					evalConditionRaw := make(map[string]interface{})
					if severityConfigurationsChildRaw["evalCondition"] != nil {
						evalConditionRaw = severityConfigurationsChildRaw["evalCondition"].(map[string]interface{})
					}
					if len(evalConditionRaw) > 0 {
						evalConditionMap["condition"] = evalConditionRaw["condition"]
						evalConditionMap["count_condition"] = evalConditionRaw["countCondition"]

						evalConditionMaps = append(evalConditionMaps, evalConditionMap)
					}
					severityConfigurationsMap["eval_condition"] = evalConditionMaps
					severityConfigurationsMaps = append(severityConfigurationsMaps, severityConfigurationsMap)
				}
			}
			configurationMap["severity_configurations"] = severityConfigurationsMaps
			sinkAlerthubMaps := make([]map[string]interface{}, 0)
			sinkAlerthubMap := make(map[string]interface{})
			sinkAlerthubRaw := make(map[string]interface{})
			if configurationRaw["sinkAlerthub"] != nil {
				sinkAlerthubRaw = configurationRaw["sinkAlerthub"].(map[string]interface{})
			}
			if len(sinkAlerthubRaw) > 0 {
				sinkAlerthubMap["enabled"] = sinkAlerthubRaw["enabled"]

				sinkAlerthubMaps = append(sinkAlerthubMaps, sinkAlerthubMap)
			}
			configurationMap["sink_alerthub"] = sinkAlerthubMaps
			sinkCmsMaps := make([]map[string]interface{}, 0)
			sinkCmsMap := make(map[string]interface{})
			sinkCmsRaw := make(map[string]interface{})
			if configurationRaw["sinkCms"] != nil {
				sinkCmsRaw = configurationRaw["sinkCms"].(map[string]interface{})
			}
			if len(sinkCmsRaw) > 0 {
				sinkCmsMap["enabled"] = sinkCmsRaw["enabled"]

				sinkCmsMaps = append(sinkCmsMaps, sinkCmsMap)
			}
			configurationMap["sink_cms"] = sinkCmsMaps
			sinkEventStoreMaps := make([]map[string]interface{}, 0)
			sinkEventStoreMap := make(map[string]interface{})
			sinkEventStoreRaw := make(map[string]interface{})
			if configurationRaw["sinkEventStore"] != nil {
				sinkEventStoreRaw = configurationRaw["sinkEventStore"].(map[string]interface{})
			}
			if len(sinkEventStoreRaw) > 0 {
				sinkEventStoreMap["enabled"] = sinkEventStoreRaw["enabled"]
				sinkEventStoreMap["endpoint"] = sinkEventStoreRaw["endpoint"]
				sinkEventStoreMap["event_store"] = sinkEventStoreRaw["eventStore"]
				sinkEventStoreMap["project"] = sinkEventStoreRaw["project"]
				sinkEventStoreMap["role_arn"] = sinkEventStoreRaw["roleArn"]

				sinkEventStoreMaps = append(sinkEventStoreMaps, sinkEventStoreMap)
			}
			configurationMap["sink_event_store"] = sinkEventStoreMaps
			tagsRaw := make([]interface{}, 0)
			if configurationRaw["tags"] != nil {
				tagsRaw = configurationRaw["tags"].([]interface{})
			}

			configurationMap["tags"] = tagsRaw
			templateConfigurationMaps := make([]map[string]interface{}, 0)
			templateConfigurationMap := make(map[string]interface{})
			templateConfigurationRaw := make(map[string]interface{})
			if configurationRaw["templateConfiguration"] != nil {
				templateConfigurationRaw = configurationRaw["templateConfiguration"].(map[string]interface{})
			}
			if len(templateConfigurationRaw) > 0 {
				templateConfigurationMap["annotations"] = templateConfigurationRaw["aonotations"]
				templateConfigurationMap["lang"] = templateConfigurationRaw["lang"]
				templateConfigurationMap["template_id"] = templateConfigurationRaw["id"]
				templateConfigurationMap["tokens"] = templateConfigurationRaw["tokens"]
				templateConfigurationMap["type"] = templateConfigurationRaw["type"]
				templateConfigurationMap["version"] = templateConfigurationRaw["version"]

				templateConfigurationMaps = append(templateConfigurationMaps, templateConfigurationMap)
			}
			configurationMap["template_configuration"] = templateConfigurationMaps
			configurationMaps = append(configurationMaps, configurationMap)
		}
		mapping["configuration"] = configurationMaps
		scheduleMaps := make([]map[string]interface{}, 0)
		scheduleMap := make(map[string]interface{})
		scheduleRaw := make(map[string]interface{})
		if objectRaw["schedule"] != nil {
			scheduleRaw = objectRaw["schedule"].(map[string]interface{})
		}
		if len(scheduleRaw) > 0 {
			scheduleMap["cron_expression"] = scheduleRaw["cronExpression"]
			scheduleMap["delay"] = scheduleRaw["delay"]
			scheduleMap["interval"] = scheduleRaw["interval"]
			scheduleMap["run_immdiately"] = scheduleRaw["runImmediately"]
			scheduleMap["time_zone"] = scheduleRaw["timeZone"]
			scheduleMap["type"] = scheduleRaw["type"]

			scheduleMaps = append(scheduleMaps, scheduleMap)
		}
		mapping["schedule"] = scheduleMaps

		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, objectRaw["name"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if err := d.Set("alerts", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
