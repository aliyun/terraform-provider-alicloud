// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCmsAlertRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAlertRuleV2Create,
		Read:   resourceAliCloudCmsAlertRuleV2Read,
		Update: resourceAliCloudCmsAlertRuleV2Update,
		Delete: resourceAliCloudCmsAlertRuleV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_integration_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"arms_integration_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"condition_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operator": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"compare_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"aggregate": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"yoy_time_value": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"threshold": {
										Type:     schema.TypeFloat,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"yoy_time_unit": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
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
							Required: true,
							ForceNew: true,
						},
						"relation": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"aggregate": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"threshold_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"threshold": {
										Type:     schema.TypeFloat,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
					},
				},
			},
			"content_template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datasource_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"legacy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
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
				Optional: true,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notify_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"utc_offset": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"notify_strategies": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"active_days": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"active_end_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"silence_time_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"active_start_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"channels": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"identifiers": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										ForceNew: true,
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
				Optional: true,
				Computed: true,
				MaxItems: 1,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"field": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"label_filters": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operator": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
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
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"entity_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"measure_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"measure_code": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"group_by": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"window_secs": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"expr": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"entity_domain": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"relation_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"service_id_list": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enable_data_complete_check": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"entity_fields": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"field": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"schedule_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"interval_secs": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							ForceNew: true,
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
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAliCloudCmsAlertRuleV2Create(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/manageAlertRules")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("alert_rule_v2_id"); ok {
		request["uuid"] = v
	}
	query["RegionId"] = StringPointer(client.RegionId)

	if v := d.Get("notify_config"); !IsNil(v) {
		request["notifyConfig"] = expandCmsAlertRuleV2NotifyConfig(v)
	}

	if v := d.Get("query_config"); !IsNil(v) {
		request["queryConfig"] = expandCmsAlertRuleV2QueryConfig(v)
	}

	if v := d.Get("condition_config"); !IsNil(v) {
		request["conditionConfig"] = expandCmsAlertRuleV2ConditionConfig(v)
	}

	armsIntegrationConfig := make(map[string]interface{})

	if v := d.Get("arms_integration_config"); !IsNil(v) {
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			armsIntegrationConfig["enabled"] = enabled1
		}

		request["armsIntegrationConfig"] = armsIntegrationConfig
	}

	if v, ok := d.GetOkExists("enabled"); ok {
		request["enabled"] = v
	}
	if v, ok := d.GetOk("labels"); ok {
		request["labels"] = v
	}
	actionIntegrationConfig := make(map[string]interface{})

	if v := d.Get("action_integration_config"); !IsNil(v) {
		actions1, _ := jsonpath.Get("$[0].actions", v)
		if actionsList, ok := actions1.([]interface{}); ok && len(actionsList) > 0 {
			actionIntegrationConfig["actions"] = actions1
		}
		enabled3, _ := jsonpath.Get("$[0].enabled", v)
		if enabled3 != nil && enabled3 != "" {
			actionIntegrationConfig["enabled"] = enabled3
		}

		request["actionIntegrationConfig"] = actionIntegrationConfig
	}

	datasourceConfig := make(map[string]interface{})

	if v := d.Get("datasource_config"); !IsNil(v) {
		instanceId1, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId1 != nil && instanceId1 != "" {
			datasourceConfig["instanceId"] = instanceId1
		}
		legacyRaw5, _ := jsonpath.Get("$[0].legacy_raw", v)
		if legacyRaw5 != nil && legacyRaw5 != "" {
			datasourceConfig["legacyRaw"] = legacyRaw5
		}
		legacyType5, _ := jsonpath.Get("$[0].legacy_type", v)
		if legacyType5 != nil && legacyType5 != "" {
			datasourceConfig["legacyType"] = legacyType5
		}
		productCategory1, _ := jsonpath.Get("$[0].product_category", v)
		if productCategory1 != nil && productCategory1 != "" {
			datasourceConfig["productCategory"] = productCategory1
		}
		type7, _ := jsonpath.Get("$[0].type", v)
		if type7 != nil && type7 != "" {
			datasourceConfig["type"] = type7
		}
		regionId1, _ := jsonpath.Get("$[0].region_id", v)
		if regionId1 != nil && regionId1 != "" {
			datasourceConfig["regionId"] = regionId1
		}

		request["datasourceConfig"] = datasourceConfig
	}

	if v, ok := d.GetOk("annotations"); ok {
		request["annotations"] = v
	}
	if v, ok := d.GetOk("display_name"); ok {
		request["displayName"] = v
	}
	scheduleConfig := make(map[string]interface{})

	if v := d.Get("schedule_config"); !IsNil(v) {
		type9, _ := jsonpath.Get("$[0].type", v)
		if type9 != nil && type9 != "" {
			scheduleConfig["type"] = type9
		}
		intervalSecs1, _ := jsonpath.Get("$[0].interval_secs", v)
		if intervalSecsInt, ok := intervalSecs1.(int); ok && intervalSecsInt != 0 {
			scheduleConfig["intervalSecs"] = intervalSecs1
		}

		request["scheduleConfig"] = scheduleConfig
	}

	if v, ok := d.GetOk("content_template"); ok {
		request["contentTemplate"] = v
	}
	if v, ok := d.GetOk("workspace"); ok {
		request["workspace"] = v
	}
	request["action"] = "CREATE"
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_alert_rule_v2", action, AlibabaCloudSdkGoERROR)
	}

	id, _ := jsonpath.Get("$.data.alertRule.uuid", response)
	d.SetId(fmt.Sprint(id))

	return resourceAliCloudCmsAlertRuleV2Read(d, meta)
}

func resourceAliCloudCmsAlertRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsServiceV2 := CmsServiceV2{client}

	objectRaw, err := cmsServiceV2.DescribeCmsAlertRuleV2(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_alert_rule_v2 DescribeCmsAlertRuleV2 Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("annotations", objectRaw["annotations"])
	d.Set("content_template", objectRaw["contentTemplate"])
	d.Set("created_at", objectRaw["createdAt"])
	d.Set("datasource_type", objectRaw["datasourceType"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("enabled", objectRaw["enabled"])
	if labelsRaw, ok := objectRaw["labels"].(map[string]interface{}); ok {
		userLabels := make(map[string]interface{})
		for labelKey, labelValue := range labelsRaw {
			if strings.HasPrefix(labelKey, "_cms_") {
				continue
			}
			userLabels[labelKey] = labelValue
		}
		d.Set("labels", userLabels)
	} else {
		d.Set("labels", objectRaw["labels"])
	}
	d.Set("notify_strategy_id", objectRaw["notifyStrategyId"])
	d.Set("observe_resource_global_scope", objectRaw["observeResourceGlobalScope"])
	d.Set("observe_resource_type", objectRaw["observeResourceType"])
	d.Set("partition_key", objectRaw["partitionKey"])
	d.Set("severity_levels", objectRaw["severityLevels"])
	d.Set("status", objectRaw["status"])
	d.Set("updated_at", objectRaw["updatedAt"])
	d.Set("workspace", objectRaw["workspace"])
	d.Set("alert_rule_v2_id", objectRaw["uuid"])

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
	if err := d.Set("action_integration_config", actionIntegrationConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("arms_integration_config", armsIntegrationConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("condition_config", conditionConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("datasource_config", datasourceConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("notify_config", notifyConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("query_config", queryConfigMaps); err != nil {
		return err
	}
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
	if err := d.Set("schedule_config", scheduleConfigMaps); err != nil {
		return err
	}

	d.Set("alert_rule_v2_id", d.Id())

	return nil
}

func resourceAliCloudCmsAlertRuleV2Update(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false

	var err error
	action := fmt.Sprintf("/manageAlertRules")
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	request["uuid"] = d.Id()
	query["RegionId"] = StringPointer(client.RegionId)
	if d.HasChange("notify_config") {
		update = true
	}
	if v := d.Get("notify_config"); !IsNil(v) || d.HasChange("notify_config") {
		request["notifyConfig"] = expandCmsAlertRuleV2NotifyConfig(v)
	}

	if d.HasChange("query_config") {
		update = true
	}
	if v := d.Get("query_config"); !IsNil(v) || d.HasChange("query_config") {
		request["queryConfig"] = expandCmsAlertRuleV2QueryConfig(v)
	}

	if d.HasChange("condition_config") {
		update = true
	}
	if v := d.Get("condition_config"); !IsNil(v) || d.HasChange("condition_config") {
		request["conditionConfig"] = expandCmsAlertRuleV2ConditionConfig(v)
	}

	if d.HasChange("arms_integration_config") {
		update = true
	}
	armsIntegrationConfig := make(map[string]interface{})

	if v := d.Get("arms_integration_config"); !IsNil(v) || d.HasChange("arms_integration_config") {
		enabled1, _ := jsonpath.Get("$[0].enabled", v)
		if enabled1 != nil && enabled1 != "" {
			armsIntegrationConfig["enabled"] = enabled1
		}

		request["armsIntegrationConfig"] = armsIntegrationConfig
	}

	if d.HasChange("enabled") {
		update = true
	}
	if v, ok := d.GetOkExists("enabled"); ok || d.HasChange("enabled") {
		request["enabled"] = v
	}
	if d.HasChange("labels") {
		update = true
	}
	if v, ok := d.GetOk("labels"); ok || d.HasChange("labels") {
		request["labels"] = v
	}
	if d.HasChange("action_integration_config") {
		update = true
	}
	actionIntegrationConfig := make(map[string]interface{})

	if v := d.Get("action_integration_config"); !IsNil(v) || d.HasChange("action_integration_config") {
		actions1, _ := jsonpath.Get("$[0].actions", v)
		if actionsList, ok := actions1.([]interface{}); ok && len(actionsList) > 0 {
			actionIntegrationConfig["actions"] = actions1
		}
		enabled3, _ := jsonpath.Get("$[0].enabled", v)
		if enabled3 != nil && enabled3 != "" {
			actionIntegrationConfig["enabled"] = enabled3
		}

		request["actionIntegrationConfig"] = actionIntegrationConfig
	}

	if d.HasChange("datasource_config") {
		update = true
	}
	datasourceConfig := make(map[string]interface{})

	if v := d.Get("datasource_config"); !IsNil(v) || d.HasChange("datasource_config") {
		instanceId1, _ := jsonpath.Get("$[0].instance_id", v)
		if instanceId1 != nil && instanceId1 != "" {
			datasourceConfig["instanceId"] = instanceId1
		}
		legacyRaw5, _ := jsonpath.Get("$[0].legacy_raw", v)
		if legacyRaw5 != nil && legacyRaw5 != "" {
			datasourceConfig["legacyRaw"] = legacyRaw5
		}
		legacyType5, _ := jsonpath.Get("$[0].legacy_type", v)
		if legacyType5 != nil && legacyType5 != "" {
			datasourceConfig["legacyType"] = legacyType5
		}
		productCategory1, _ := jsonpath.Get("$[0].product_category", v)
		if productCategory1 != nil && productCategory1 != "" {
			datasourceConfig["productCategory"] = productCategory1
		}
		type7, _ := jsonpath.Get("$[0].type", v)
		if type7 != nil && type7 != "" {
			datasourceConfig["type"] = type7
		}
		regionId1, _ := jsonpath.Get("$[0].region_id", v)
		if regionId1 != nil && regionId1 != "" {
			datasourceConfig["regionId"] = regionId1
		}

		request["datasourceConfig"] = datasourceConfig
	}

	if d.HasChange("annotations") {
		update = true
	}
	if v, ok := d.GetOk("annotations"); ok || d.HasChange("annotations") {
		request["annotations"] = v
	}
	if d.HasChange("display_name") {
		update = true
	}
	if v, ok := d.GetOk("display_name"); ok || d.HasChange("display_name") {
		request["displayName"] = v
	}
	if d.HasChange("schedule_config") {
		update = true
	}
	scheduleConfig := make(map[string]interface{})

	if v := d.Get("schedule_config"); !IsNil(v) || d.HasChange("schedule_config") {
		type9, _ := jsonpath.Get("$[0].type", v)
		if type9 != nil && type9 != "" {
			scheduleConfig["type"] = type9
		}
		intervalSecs1, _ := jsonpath.Get("$[0].interval_secs", v)
		if intervalSecsInt, ok := intervalSecs1.(int); ok && intervalSecsInt != 0 {
			scheduleConfig["intervalSecs"] = intervalSecs1
		}

		request["scheduleConfig"] = scheduleConfig
	}

	if d.HasChange("content_template") {
		update = true
	}
	if v, ok := d.GetOk("content_template"); ok || d.HasChange("content_template") {
		request["contentTemplate"] = v
	}
	if d.HasChange("workspace") {
		update = true
	}
	if v, ok := d.GetOk("workspace"); ok || d.HasChange("workspace") {
		request["workspace"] = v
	}
	request["action"] = "UPDATE"
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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

	return resourceAliCloudCmsAlertRuleV2Read(d, meta)
}

func resourceAliCloudCmsAlertRuleV2Delete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := fmt.Sprintf("/manageAlertRules")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["uuid"] = d.Id()
	query["RegionId"] = StringPointer(client.RegionId)

	if v, ok := d.GetOk("workspace"); ok {
		request["workspace"] = v
	}
	request["action"] = "DELETE"
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaPost("Cms", "2024-03-30", action, query, nil, body, true)
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
		if NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func expandCmsAlertRuleV2ObjectList(l []interface{}, stringKeys map[string]string, intKeys map[string]string, floatKeys map[string]string, listKeys map[string]string) []interface{} {
	items := make([]interface{}, 0, len(l))
	for _, item := range l {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		out := make(map[string]interface{})
		for tfKey, apiKey := range stringKeys {
			if s, ok := itemMap[tfKey].(string); ok && s != "" {
				out[apiKey] = s
			}
		}
		for tfKey, apiKey := range intKeys {
			if n, ok := itemMap[tfKey].(int); ok && n != 0 {
				out[apiKey] = n
			}
		}
		for tfKey, apiKey := range floatKeys {
			if n, ok := itemMap[tfKey].(float64); ok {
				out[apiKey] = n
			}
		}
		for tfKey, apiKey := range listKeys {
			if lv, ok := itemMap[tfKey].([]interface{}); ok && len(lv) > 0 {
				out[apiKey] = lv
			}
		}
		if len(out) > 0 {
			items = append(items, out)
		}
	}
	return items
}

func expandCmsAlertRuleV2NotifyConfig(v interface{}) map[string]interface{} {
	notifyConfig := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return notifyConfig
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return notifyConfig
	}
	for tfKey, apiKey := range map[string]string{"type": "type", "utc_offset": "utcOffset", "active_start_time": "activeStartTime", "active_end_time": "activeEndTime"} {
		if s, ok := m[tfKey].(string); ok && s != "" {
			notifyConfig[apiKey] = s
		}
	}
	if n, ok := m["silence_time_secs"].(int); ok && n != 0 {
		notifyConfig["silenceTimeSecs"] = n
	}
	if l, ok := m["active_days"].([]interface{}); ok && len(l) > 0 {
		notifyConfig["activeDays"] = l
	}
	if l, ok := m["notify_strategies"].([]interface{}); ok && len(l) > 0 {
		notifyConfig["notifyStrategies"] = l
	}
	if l, ok := m["channels"].([]interface{}); ok && len(l) > 0 {
		notifyConfig["channels"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"type": "type"}, nil, nil,
			map[string]string{"identifiers": "identifiers"})
	}
	return notifyConfig
}

func expandCmsAlertRuleV2QueryConfig(v interface{}) map[string]interface{} {
	queryConfig := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return queryConfig
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return queryConfig
	}
	queryType, _ := m["type"].(string)
	for tfKey, apiKey := range map[string]string{"type": "type", "expr": "expr", "prom_ql": "promQl", "metric_set": "metricSet", "metric": "metric", "entity_domain": "entityDomain", "entity_type": "entityType", "namespace": "namespace", "relation_type": "relationType", "group_id": "groupId", "legacy_type": "legacyType", "legacy_raw": "legacyRaw"} {
		if s, ok := m[tfKey].(string); ok && s != "" {
			queryConfig[apiKey] = s
		}
	}
	if queryType == "PROMETHEUS_SINGLE_QUERY" {
		if b, ok := m["enable_data_complete_check"].(bool); ok && b {
			queryConfig["enableDataCompleteCheck"] = b
		}
	}
	if l, ok := m["service_id_list"].([]interface{}); ok && len(l) > 0 {
		queryConfig["serviceIdList"] = l
	}
	if l, ok := m["dimensions"].([]interface{}); ok && len(l) > 0 {
		queryConfig["dimensions"] = l
	}
	if l, ok := m["entity_filters"].([]interface{}); ok && len(l) > 0 {
		queryConfig["entityFilters"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"field": "field", "operator": "operator", "value": "value"}, nil, nil, nil)
	}
	if l, ok := m["label_filters"].([]interface{}); ok && len(l) > 0 {
		queryConfig["labelFilters"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"name": "name", "operator": "operator", "value": "value"}, nil, nil, nil)
	}
	if l, ok := m["entity_fields"].([]interface{}); ok && len(l) > 0 {
		queryConfig["entityFields"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"field": "field", "value": "value"}, nil, nil, nil)
	}
	if l, ok := m["filter_list"].([]interface{}); ok && len(l) > 0 {
		queryConfig["filterList"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"key": "key", "type": "type", "value": "value"}, nil, nil, nil)
	}
	if l, ok := m["measure_list"].([]interface{}); ok && len(l) > 0 {
		queryConfig["measureList"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"measure_code": "measureCode"},
			map[string]string{"window_secs": "windowSecs"}, nil,
			map[string]string{"group_by": "groupBy"})
	}
	return queryConfig
}

func expandCmsAlertRuleV2ConditionConfig(v interface{}) map[string]interface{} {
	conditionConfig := make(map[string]interface{})
	raw, ok := v.([]interface{})
	if !ok || len(raw) == 0 || raw[0] == nil {
		return conditionConfig
	}
	m, ok := raw[0].(map[string]interface{})
	if !ok {
		return conditionConfig
	}
	for tfKey, apiKey := range map[string]string{"type": "type", "severity": "severity", "operator": "operator", "aggregate": "aggregate", "relation": "relation", "escalation_type": "escalationType", "no_data_policy": "noDataPolicy", "yoy_time_unit": "yoyTimeUnit", "legacy_type": "legacyType", "legacy_raw": "legacyRaw"} {
		if s, ok := m[tfKey].(string); ok && s != "" {
			conditionConfig[apiKey] = s
		}
	}
	if n, ok := m["duration_secs"].(int); ok && n != 0 {
		conditionConfig["durationSecs"] = n
	}
	if n, ok := m["yoy_time_value"].(int); ok && n != 0 {
		conditionConfig["yoyTimeValue"] = n
	}
	if n, ok := m["threshold"].(float64); ok && n != 0 {
		conditionConfig["threshold"] = n
	}
	if l, ok := m["threshold_list"].([]interface{}); ok && len(l) > 0 {
		conditionConfig["thresholdList"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"severity": "severity"}, nil,
			map[string]string{"threshold": "threshold"}, nil)
	}
	if l, ok := m["compare_list"].([]interface{}); ok && len(l) > 0 {
		conditionConfig["compareList"] = expandCmsAlertRuleV2ObjectList(l,
			map[string]string{"aggregate": "aggregate", "operator": "operator", "yoy_time_unit": "yoyTimeUnit"},
			map[string]string{"yoy_time_value": "yoyTimeValue"},
			map[string]string{"threshold": "threshold"}, nil)
	}
	if l, ok := m["simple_escalation"].([]interface{}); ok && len(l) > 0 && l[0] != nil {
		if sm, ok := l[0].(map[string]interface{}); ok {
			simpleEscalation := make(map[string]interface{})
			if s, ok := sm["metric_name"].(string); ok && s != "" {
				simpleEscalation["metricName"] = s
			}
			if n, ok := sm["period"].(int); ok && n != 0 {
				simpleEscalation["period"] = n
			}
			if el, ok := sm["escalations"].([]interface{}); ok && len(el) > 0 {
				simpleEscalation["escalations"] = expandCmsAlertRuleV2ObjectList(el,
					map[string]string{"comparison_operator": "comparisonOperator", "statistics": "statistics", "threshold": "threshold", "severity": "severity", "pre_condition": "preCondition"},
					map[string]string{"times": "times"}, nil, nil)
			}
			if len(simpleEscalation) > 0 {
				conditionConfig["simpleEscalation"] = simpleEscalation
			}
		}
	}
	if l, ok := m["composite_escalation"].([]interface{}); ok && len(l) > 0 && l[0] != nil {
		if cm, ok := l[0].(map[string]interface{}); ok {
			compositeEscalation := make(map[string]interface{})
			for tfKey, apiKey := range map[string]string{"relation": "relation", "severity": "severity"} {
				if s, ok := cm[tfKey].(string); ok && s != "" {
					compositeEscalation[apiKey] = s
				}
			}
			if n, ok := cm["times"].(int); ok && n != 0 {
				compositeEscalation["times"] = n
			}
			if el, ok := cm["escalations"].([]interface{}); ok && len(el) > 0 {
				compositeEscalation["escalations"] = expandCmsAlertRuleV2ObjectList(el,
					map[string]string{"comparison_operator": "comparisonOperator", "metric_name": "metricName", "statistics": "statistics", "threshold": "threshold", "pre_condition": "preCondition"},
					map[string]string{"period": "period"}, nil, nil)
			}
			if len(compositeEscalation) > 0 {
				conditionConfig["compositeEscalation"] = compositeEscalation
			}
		}
	}
	if l, ok := m["express_escalation"].([]interface{}); ok && len(l) > 0 && l[0] != nil {
		if em, ok := l[0].(map[string]interface{}); ok {
			expressEscalation := make(map[string]interface{})
			for tfKey, apiKey := range map[string]string{"raw_expression": "rawExpression", "severity": "severity"} {
				if s, ok := em[tfKey].(string); ok && s != "" {
					expressEscalation[apiKey] = s
				}
			}
			if n, ok := em["times"].(int); ok && n != 0 {
				expressEscalation["times"] = n
			}
			if len(expressEscalation) > 0 {
				conditionConfig["expressEscalation"] = expressEscalation
			}
		}
	}
	if l, ok := m["prometheus"].([]interface{}); ok && len(l) > 0 && l[0] != nil {
		if pm, ok := l[0].(map[string]interface{}); ok {
			prometheus := make(map[string]interface{})
			for tfKey, apiKey := range map[string]string{"prom_ql": "promQl", "severity": "severity"} {
				if s, ok := pm[tfKey].(string); ok && s != "" {
					prometheus[apiKey] = s
				}
			}
			if n, ok := pm["times"].(int); ok && n != 0 {
				prometheus["times"] = n
			}
			if len(prometheus) > 0 {
				conditionConfig["prometheus"] = prometheus
			}
		}
	}
	return conditionConfig
}
