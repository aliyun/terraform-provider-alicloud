// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudSlsAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudSlsAlertCreate,
		Read:   resourceAliCloudSlsAlertRead,
		Update: resourceAliCloudSlsAlertUpdate,
		Delete: resourceAliCloudSlsAlertDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"alert_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity_configurations": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"eval_condition": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"condition": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"count_condition": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"auto_annotation": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"no_data_fire": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"sink_cms": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"dashboard": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"mute_until": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"template_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"annotations": {
										Type:     schema.TypeMap,
										Optional: true,
									},
									"version": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"lang": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"template_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"tokens": {
										Type:     schema.TypeMap,
										Optional: true,
									},
								},
							},
						},
						"labels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"group_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"no_group", "custom", "labels_auto"}, true),
									},
									"fields": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"no_data_severity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"condition_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"count_condition": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"join_configurations": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: StringInSlice([]string{"cross_join", "inner_join", "left_join", "right_join", "full_join", "left_exclude", "right_exclude", "concat", "no_join"}, true),
									},
								},
							},
						},
						"policy_configuration": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alert_policy_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"action_policy_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"repeat_interval": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"sink_event_store": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"event_store": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"endpoint": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"role_arn": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"query_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"query": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"time_span_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"start": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"store": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"power_sql_mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dashboard_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"role_arn": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"store_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"project": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"ui": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"region": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"end": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"chart_title": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"sink_alerthub": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"tags": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"send_resolved": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"FixedRate", "Cron"}, true),
						},
						"time_zone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"run_immdiately": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"cron_expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"delay": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"ENABLED", "DISABLED"}, true),
			},
		},
	}
}

func resourceAliCloudSlsAlertCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := fmt.Sprintf("/alerts")
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(d.Get("project_name").(string))
	request["name"] = d.Get("alert_name")

	request["displayName"] = d.Get("display_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("configuration"); v != nil {
		nodeNative, _ := jsonpath.Get("$[0].mute_until", d.Get("configuration"))
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["muteUntil"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].version", d.Get("configuration"))
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["version"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].type", d.Get("configuration"))
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["type"] = nodeNative2
		}
		templateConfiguration := make(map[string]interface{})
		nodeNative3, _ := jsonpath.Get("$[0].template_configuration[0].type", d.Get("configuration"))
		if nodeNative3 != nil && nodeNative3 != "" {
			templateConfiguration["type"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].template_configuration[0].version", d.Get("configuration"))
		if nodeNative4 != nil && nodeNative4 != "" {
			templateConfiguration["version"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].template_configuration[0].lang", d.Get("configuration"))
		if nodeNative5 != nil && nodeNative5 != "" {
			templateConfiguration["lang"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].template_configuration[0].template_id", d.Get("configuration"))
		if nodeNative6 != nil && nodeNative6 != "" {
			templateConfiguration["id"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].template_configuration[0].tokens", d.Get("configuration"))
		if nodeNative7 != nil && nodeNative7 != "" {
			templateConfiguration["tokens"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].template_configuration[0].annotations", d.Get("configuration"))
		if nodeNative8 != nil && nodeNative8 != "" {
			templateConfiguration["aonotations"] = nodeNative8
		}

		objectDataLocalMap["templateConfiguration"] = templateConfiguration
		nodeNative9, _ := jsonpath.Get("$[0].dashboard", d.Get("configuration"))
		if nodeNative9 != nil && nodeNative9 != "" {
			objectDataLocalMap["dashboard"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].threshold", d.Get("configuration"))
		if nodeNative10 != nil && nodeNative10 != "" {
			objectDataLocalMap["threshold"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].no_data_fire", d.Get("configuration"))
		if nodeNative11 != nil && nodeNative11 != "" {
			objectDataLocalMap["noDataFire"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].no_data_severity", d.Get("configuration"))
		if nodeNative12 != nil && nodeNative12 != "" {
			objectDataLocalMap["noDataSeverity"] = nodeNative12
		}
		nodeNative13, _ := jsonpath.Get("$[0].send_resolved", d.Get("configuration"))
		if nodeNative13 != nil && nodeNative13 != "" {
			objectDataLocalMap["sendResolved"] = nodeNative13
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData, err := jsonpath.Get("$[0].query_list", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["chartTitle"] = dataLoopTmp["chart_title"]
				dataLoopMap["query"] = dataLoopTmp["query"]
				dataLoopMap["timeSpanType"] = dataLoopTmp["time_span_type"]
				dataLoopMap["start"] = dataLoopTmp["start"]
				dataLoopMap["end"] = dataLoopTmp["end"]
				dataLoopMap["storeType"] = dataLoopTmp["store_type"]
				dataLoopMap["project"] = dataLoopTmp["project"]
				dataLoopMap["store"] = dataLoopTmp["store"]
				dataLoopMap["region"] = dataLoopTmp["region"]
				dataLoopMap["roleArn"] = dataLoopTmp["role_arn"]
				dataLoopMap["dashboardId"] = dataLoopTmp["dashboard_id"]
				dataLoopMap["powerSqlMode"] = dataLoopTmp["power_sql_mode"]
				dataLoopMap["ui"] = dataLoopTmp["ui"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["queryList"] = localMaps
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData1, err := jsonpath.Get("$[0].annotations", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["key"] = dataLoop1Tmp["key"]
				dataLoop1Map["value"] = dataLoop1Tmp["value"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			objectDataLocalMap["annotations"] = localMaps1
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData2, err := jsonpath.Get("$[0].labels", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps2 := make([]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := dataLoop2.(map[string]interface{})
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["key"] = dataLoop2Tmp["key"]
				dataLoop2Map["value"] = dataLoop2Tmp["value"]
				localMaps2 = append(localMaps2, dataLoop2Map)
			}
			objectDataLocalMap["labels"] = localMaps2
		}
		conditionConfiguration := make(map[string]interface{})
		nodeNative31, _ := jsonpath.Get("$[0].condition_configuration[0].count_condition", d.Get("configuration"))
		if nodeNative31 != nil && nodeNative31 != "" {
			conditionConfiguration["countCondition"] = nodeNative31
		}
		nodeNative32, _ := jsonpath.Get("$[0].condition_configuration[0].condition", d.Get("configuration"))
		if nodeNative32 != nil && nodeNative32 != "" {
			conditionConfiguration["condition"] = nodeNative32
		}

		objectDataLocalMap["conditionConfiguration"] = conditionConfiguration
		if v, ok := d.GetOk("configuration"); ok {
			localData3, err := jsonpath.Get("$[0].severity_configurations", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps3 := make([]interface{}, 0)
			for _, dataLoop3 := range localData3.([]interface{}) {
				dataLoop3Tmp := dataLoop3.(map[string]interface{})
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["severity"] = dataLoop3Tmp["severity"]
				if !IsNil(dataLoop3Tmp["eval_condition"]) {
					localData4 := make(map[string]interface{})
					nodeNative34, _ := jsonpath.Get("$[0].condition", dataLoop3Tmp["eval_condition"])
					if nodeNative34 != nil && nodeNative34 != "" {
						localData4["condition"] = nodeNative34
					}
					nodeNative35, _ := jsonpath.Get("$[0].count_condition", dataLoop3Tmp["eval_condition"])
					if nodeNative35 != nil && nodeNative35 != "" {
						localData4["countCondition"] = nodeNative35
					}
					dataLoop3Map["evalCondition"] = localData4
				}
				localMaps3 = append(localMaps3, dataLoop3Map)
			}
			objectDataLocalMap["severityConfigurations"] = localMaps3
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData5, err := jsonpath.Get("$[0].join_configurations", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps5 := make([]interface{}, 0)
			for _, dataLoop5 := range localData5.([]interface{}) {
				dataLoop5Tmp := dataLoop5.(map[string]interface{})
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["type"] = dataLoop5Tmp["type"]
				dataLoop5Map["condition"] = dataLoop5Tmp["condition"]
				localMaps5 = append(localMaps5, dataLoop5Map)
			}
			objectDataLocalMap["joinConfigurations"] = localMaps5
		}
		groupConfiguration := make(map[string]interface{})
		nodeNative38, _ := jsonpath.Get("$[0].group_configuration[0].type", d.Get("configuration"))
		if nodeNative38 != nil && nodeNative38 != "" {
			groupConfiguration["type"] = nodeNative38
		}
		nodeNative39, _ := jsonpath.Get("$[0].group_configuration[0].fields", v)
		if nodeNative39 != nil && nodeNative39 != "" {
			groupConfiguration["fields"] = nodeNative39
		}

		objectDataLocalMap["groupConfiguration"] = groupConfiguration
		policyConfiguration := make(map[string]interface{})
		nodeNative40, _ := jsonpath.Get("$[0].policy_configuration[0].alert_policy_id", d.Get("configuration"))
		if nodeNative40 != nil && nodeNative40 != "" {
			policyConfiguration["alertPolicyId"] = nodeNative40
		}
		nodeNative41, _ := jsonpath.Get("$[0].policy_configuration[0].action_policy_id", d.Get("configuration"))
		if nodeNative41 != nil && nodeNative41 != "" {
			policyConfiguration["actionPolicyId"] = nodeNative41
		}
		nodeNative42, _ := jsonpath.Get("$[0].policy_configuration[0].repeat_interval", d.Get("configuration"))
		if nodeNative42 != nil && nodeNative42 != "" {
			policyConfiguration["repeatInterval"] = nodeNative42
		}

		objectDataLocalMap["policyConfiguration"] = policyConfiguration
		nodeNative43, _ := jsonpath.Get("$[0].auto_annotation", d.Get("configuration"))
		if nodeNative43 != nil && nodeNative43 != "" {
			objectDataLocalMap["autoAnnotation"] = nodeNative43
		}
		sinkEventStore := make(map[string]interface{})
		nodeNative44, _ := jsonpath.Get("$[0].sink_event_store[0].enabled", d.Get("configuration"))
		if nodeNative44 != nil && nodeNative44 != "" {
			sinkEventStore["enabled"] = nodeNative44
		}
		nodeNative45, _ := jsonpath.Get("$[0].sink_event_store[0].endpoint", d.Get("configuration"))
		if nodeNative45 != nil && nodeNative45 != "" {
			sinkEventStore["endpoint"] = nodeNative45
		}
		nodeNative46, _ := jsonpath.Get("$[0].sink_event_store[0].project", d.Get("configuration"))
		if nodeNative46 != nil && nodeNative46 != "" {
			sinkEventStore["project"] = nodeNative46
		}
		nodeNative47, _ := jsonpath.Get("$[0].sink_event_store[0].event_store", d.Get("configuration"))
		if nodeNative47 != nil && nodeNative47 != "" {
			sinkEventStore["eventStore"] = nodeNative47
		}
		nodeNative48, _ := jsonpath.Get("$[0].sink_event_store[0].role_arn", d.Get("configuration"))
		if nodeNative48 != nil && nodeNative48 != "" {
			sinkEventStore["roleArn"] = nodeNative48
		}

		objectDataLocalMap["sinkEventStore"] = sinkEventStore
		sinkCms := make(map[string]interface{})
		nodeNative49, _ := jsonpath.Get("$[0].sink_cms[0].enabled", d.Get("configuration"))
		if nodeNative49 != nil && nodeNative49 != "" {
			sinkCms["enabled"] = nodeNative49
		}

		objectDataLocalMap["sinkCms"] = sinkCms
		sinkAlerthub := make(map[string]interface{})
		nodeNative50, _ := jsonpath.Get("$[0].sink_alerthub[0].enabled", d.Get("configuration"))
		if nodeNative50 != nil && nodeNative50 != "" {
			sinkAlerthub["enabled"] = nodeNative50
		}

		objectDataLocalMap["sinkAlerthub"] = sinkAlerthub
		nodeNative51, _ := jsonpath.Get("$[0].tags", v)
		if nodeNative51 != nil && nodeNative51 != "" {
			objectDataLocalMap["tags"] = nodeNative51
		}

		request["configuration"] = objectDataLocalMap
	}

	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("schedule"); v != nil {
		nodeNative52, _ := jsonpath.Get("$[0].type", d.Get("schedule"))
		if nodeNative52 != nil && nodeNative52 != "" {
			objectDataLocalMap1["type"] = nodeNative52
		}
		nodeNative53, _ := jsonpath.Get("$[0].cron_expression", d.Get("schedule"))
		if nodeNative53 != nil && nodeNative53 != "" {
			objectDataLocalMap1["cronExpression"] = nodeNative53
		}
		nodeNative54, _ := jsonpath.Get("$[0].run_immdiately", d.Get("schedule"))
		if nodeNative54 != nil && nodeNative54 != "" {
			objectDataLocalMap1["runImmediately"] = nodeNative54
		}
		nodeNative55, _ := jsonpath.Get("$[0].time_zone", d.Get("schedule"))
		if nodeNative55 != nil && nodeNative55 != "" {
			objectDataLocalMap1["timeZone"] = nodeNative55
		}
		nodeNative56, _ := jsonpath.Get("$[0].delay", d.Get("schedule"))
		if nodeNative56 != nil && nodeNative56 != "" {
			objectDataLocalMap1["delay"] = nodeNative56
		}
		nodeNative57, _ := jsonpath.Get("$[0].interval", d.Get("schedule"))
		if nodeNative57 != nil && nodeNative57 != "" {
			objectDataLocalMap1["interval"] = nodeNative57
		}

		request["schedule"] = objectDataLocalMap1
	}

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("POST", "2020-12-30", "CreateAlert", action), query, body, nil, hostMap, false)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sls_alert", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", *hostMap["project"], request["name"]))

	return resourceAliCloudSlsAlertUpdate(d, meta)
}

func resourceAliCloudSlsAlertRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	slsServiceV2 := SlsServiceV2{client}

	objectRaw, err := slsServiceV2.DescribeSlsAlert(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sls_alert DescribeSlsAlert Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("display_name", objectRaw["displayName"])
	d.Set("status", objectRaw["status"])
	d.Set("alert_name", objectRaw["name"])

	configurationMaps := make([]map[string]interface{}, 0)
	configurationMap := make(map[string]interface{})
	configuration1Raw := make(map[string]interface{})
	if objectRaw["configuration"] != nil {
		configuration1Raw = objectRaw["configuration"].(map[string]interface{})
	}
	if len(configuration1Raw) > 0 {
		configurationMap["auto_annotation"] = configuration1Raw["autoAnnotation"]
		configurationMap["dashboard"] = configuration1Raw["dashboard"]
		configurationMap["mute_until"] = configuration1Raw["muteUntil"]
		configurationMap["no_data_fire"] = configuration1Raw["noDataFire"]
		configurationMap["no_data_severity"] = configuration1Raw["noDataSeverity"]
		configurationMap["send_resolved"] = configuration1Raw["sendResolved"]
		configurationMap["threshold"] = configuration1Raw["threshold"]
		configurationMap["type"] = configuration1Raw["type"]
		configurationMap["version"] = configuration1Raw["version"]

		annotations1Raw := configuration1Raw["annotations"]
		annotationsMaps := make([]map[string]interface{}, 0)
		if annotations1Raw != nil {
			for _, annotationsChild1Raw := range annotations1Raw.([]interface{}) {
				annotationsMap := make(map[string]interface{})
				annotationsChild1Raw := annotationsChild1Raw.(map[string]interface{})
				annotationsMap["key"] = annotationsChild1Raw["key"]
				annotationsMap["value"] = annotationsChild1Raw["value"]

				annotationsMaps = append(annotationsMaps, annotationsMap)
			}
		}
		configurationMap["annotations"] = annotationsMaps
		conditionConfigurationMaps := make([]map[string]interface{}, 0)
		conditionConfigurationMap := make(map[string]interface{})
		conditionConfiguration1Raw := make(map[string]interface{})
		if configuration1Raw["conditionConfiguration"] != nil {
			conditionConfiguration1Raw = configuration1Raw["conditionConfiguration"].(map[string]interface{})
		}
		if len(conditionConfiguration1Raw) > 0 {
			conditionConfigurationMap["condition"] = conditionConfiguration1Raw["condition"]
			conditionConfigurationMap["count_condition"] = conditionConfiguration1Raw["countCondition"]

			conditionConfigurationMaps = append(conditionConfigurationMaps, conditionConfigurationMap)
		}
		configurationMap["condition_configuration"] = conditionConfigurationMaps
		groupConfigurationMaps := make([]map[string]interface{}, 0)
		groupConfigurationMap := make(map[string]interface{})
		groupConfiguration1Raw := make(map[string]interface{})
		if configuration1Raw["groupConfiguration"] != nil {
			groupConfiguration1Raw = configuration1Raw["groupConfiguration"].(map[string]interface{})
		}
		if len(groupConfiguration1Raw) > 0 {
			groupConfigurationMap["type"] = groupConfiguration1Raw["type"]

			fields1Raw := make([]interface{}, 0)
			if groupConfiguration1Raw["fields"] != nil {
				fields1Raw = groupConfiguration1Raw["fields"].([]interface{})
			}

			groupConfigurationMap["fields"] = fields1Raw
			groupConfigurationMaps = append(groupConfigurationMaps, groupConfigurationMap)
		}
		configurationMap["group_configuration"] = groupConfigurationMaps
		joinConfigurations1Raw := configuration1Raw["joinConfigurations"]
		joinConfigurationsMaps := make([]map[string]interface{}, 0)
		if joinConfigurations1Raw != nil {
			for _, joinConfigurationsChild1Raw := range joinConfigurations1Raw.([]interface{}) {
				joinConfigurationsMap := make(map[string]interface{})
				joinConfigurationsChild1Raw := joinConfigurationsChild1Raw.(map[string]interface{})
				joinConfigurationsMap["condition"] = joinConfigurationsChild1Raw["condition"]
				joinConfigurationsMap["type"] = joinConfigurationsChild1Raw["type"]

				joinConfigurationsMaps = append(joinConfigurationsMaps, joinConfigurationsMap)
			}
		}
		configurationMap["join_configurations"] = joinConfigurationsMaps
		labels1Raw := configuration1Raw["labels"]
		labelsMaps := make([]map[string]interface{}, 0)
		if labels1Raw != nil {
			for _, labelsChild1Raw := range labels1Raw.([]interface{}) {
				labelsMap := make(map[string]interface{})
				labelsChild1Raw := labelsChild1Raw.(map[string]interface{})
				labelsMap["key"] = labelsChild1Raw["key"]
				labelsMap["value"] = labelsChild1Raw["value"]

				labelsMaps = append(labelsMaps, labelsMap)
			}
		}
		configurationMap["labels"] = labelsMaps
		policyConfigurationMaps := make([]map[string]interface{}, 0)
		policyConfigurationMap := make(map[string]interface{})
		policyConfiguration1Raw := make(map[string]interface{})
		if configuration1Raw["policyConfiguration"] != nil {
			policyConfiguration1Raw = configuration1Raw["policyConfiguration"].(map[string]interface{})
		}
		if len(policyConfiguration1Raw) > 0 {
			policyConfigurationMap["action_policy_id"] = policyConfiguration1Raw["actionPolicyId"]
			policyConfigurationMap["alert_policy_id"] = policyConfiguration1Raw["alertPolicyId"]
			policyConfigurationMap["repeat_interval"] = policyConfiguration1Raw["repeatInterval"]

			policyConfigurationMaps = append(policyConfigurationMaps, policyConfigurationMap)
		}
		configurationMap["policy_configuration"] = policyConfigurationMaps
		queryList1Raw := configuration1Raw["queryList"]
		queryListMaps := make([]map[string]interface{}, 0)
		if queryList1Raw != nil {
			for _, queryListChild1Raw := range queryList1Raw.([]interface{}) {
				queryListMap := make(map[string]interface{})
				queryListChild1Raw := queryListChild1Raw.(map[string]interface{})
				queryListMap["chart_title"] = queryListChild1Raw["chartTitle"]
				queryListMap["dashboard_id"] = queryListChild1Raw["dashboardId"]
				queryListMap["end"] = queryListChild1Raw["end"]
				queryListMap["power_sql_mode"] = queryListChild1Raw["powerSqlMode"]
				queryListMap["project"] = queryListChild1Raw["project"]
				queryListMap["query"] = queryListChild1Raw["query"]
				queryListMap["region"] = queryListChild1Raw["region"]
				queryListMap["role_arn"] = queryListChild1Raw["roleArn"]
				queryListMap["start"] = queryListChild1Raw["start"]
				queryListMap["store"] = queryListChild1Raw["store"]
				queryListMap["store_type"] = queryListChild1Raw["storeType"]
				queryListMap["time_span_type"] = queryListChild1Raw["timeSpanType"]
				queryListMap["ui"] = queryListChild1Raw["ui"]

				queryListMaps = append(queryListMaps, queryListMap)
			}
		}
		configurationMap["query_list"] = queryListMaps
		severityConfigurations1Raw := configuration1Raw["severityConfigurations"]
		severityConfigurationsMaps := make([]map[string]interface{}, 0)
		if severityConfigurations1Raw != nil {
			for _, severityConfigurationsChild1Raw := range severityConfigurations1Raw.([]interface{}) {
				severityConfigurationsMap := make(map[string]interface{})
				severityConfigurationsChild1Raw := severityConfigurationsChild1Raw.(map[string]interface{})
				severityConfigurationsMap["severity"] = severityConfigurationsChild1Raw["severity"]

				evalConditionMaps := make([]map[string]interface{}, 0)
				evalConditionMap := make(map[string]interface{})
				evalCondition1Raw := make(map[string]interface{})
				if severityConfigurationsChild1Raw["evalCondition"] != nil {
					evalCondition1Raw = severityConfigurationsChild1Raw["evalCondition"].(map[string]interface{})
				}
				if len(evalCondition1Raw) > 0 {
					evalConditionMap["condition"] = evalCondition1Raw["condition"]
					evalConditionMap["count_condition"] = evalCondition1Raw["countCondition"]

					evalConditionMaps = append(evalConditionMaps, evalConditionMap)
				}
				severityConfigurationsMap["eval_condition"] = evalConditionMaps
				severityConfigurationsMaps = append(severityConfigurationsMaps, severityConfigurationsMap)
			}
		}
		configurationMap["severity_configurations"] = severityConfigurationsMaps
		sinkAlerthubMaps := make([]map[string]interface{}, 0)
		sinkAlerthubMap := make(map[string]interface{})
		sinkAlerthub1Raw := make(map[string]interface{})
		if configuration1Raw["sinkAlerthub"] != nil {
			sinkAlerthub1Raw = configuration1Raw["sinkAlerthub"].(map[string]interface{})
		}
		if len(sinkAlerthub1Raw) > 0 {
			sinkAlerthubMap["enabled"] = sinkAlerthub1Raw["enabled"]

			sinkAlerthubMaps = append(sinkAlerthubMaps, sinkAlerthubMap)
		}
		configurationMap["sink_alerthub"] = sinkAlerthubMaps
		sinkCmsMaps := make([]map[string]interface{}, 0)
		sinkCmsMap := make(map[string]interface{})
		sinkCms1Raw := make(map[string]interface{})
		if configuration1Raw["sinkCms"] != nil {
			sinkCms1Raw = configuration1Raw["sinkCms"].(map[string]interface{})
		}
		if len(sinkCms1Raw) > 0 {
			sinkCmsMap["enabled"] = sinkCms1Raw["enabled"]

			sinkCmsMaps = append(sinkCmsMaps, sinkCmsMap)
		}
		configurationMap["sink_cms"] = sinkCmsMaps
		sinkEventStoreMaps := make([]map[string]interface{}, 0)
		sinkEventStoreMap := make(map[string]interface{})
		sinkEventStore1Raw := make(map[string]interface{})
		if configuration1Raw["sinkEventStore"] != nil {
			sinkEventStore1Raw = configuration1Raw["sinkEventStore"].(map[string]interface{})
		}
		if len(sinkEventStore1Raw) > 0 {
			sinkEventStoreMap["enabled"] = sinkEventStore1Raw["enabled"]
			sinkEventStoreMap["endpoint"] = sinkEventStore1Raw["endpoint"]
			sinkEventStoreMap["event_store"] = sinkEventStore1Raw["eventStore"]
			sinkEventStoreMap["project"] = sinkEventStore1Raw["project"]
			sinkEventStoreMap["role_arn"] = sinkEventStore1Raw["roleArn"]

			sinkEventStoreMaps = append(sinkEventStoreMaps, sinkEventStoreMap)
		}
		configurationMap["sink_event_store"] = sinkEventStoreMaps
		tags1Raw := make([]interface{}, 0)
		if configuration1Raw["tags"] != nil {
			tags1Raw = configuration1Raw["tags"].([]interface{})
		}

		configurationMap["tags"] = tags1Raw
		templateConfigurationMaps := make([]map[string]interface{}, 0)
		templateConfigurationMap := make(map[string]interface{})
		templateConfiguration1Raw := make(map[string]interface{})
		if configuration1Raw["templateConfiguration"] != nil {
			templateConfiguration1Raw = configuration1Raw["templateConfiguration"].(map[string]interface{})
		}
		if len(templateConfiguration1Raw) > 0 {
			templateConfigurationMap["annotations"] = templateConfiguration1Raw["aonotations"]
			templateConfigurationMap["lang"] = templateConfiguration1Raw["lang"]
			templateConfigurationMap["template_id"] = templateConfiguration1Raw["id"]
			templateConfigurationMap["tokens"] = templateConfiguration1Raw["tokens"]
			templateConfigurationMap["type"] = templateConfiguration1Raw["type"]
			templateConfigurationMap["version"] = templateConfiguration1Raw["version"]

			templateConfigurationMaps = append(templateConfigurationMaps, templateConfigurationMap)
		}
		configurationMap["template_configuration"] = templateConfigurationMaps
		configurationMaps = append(configurationMaps, configurationMap)
	}
	d.Set("configuration", configurationMaps)
	scheduleMaps := make([]map[string]interface{}, 0)
	scheduleMap := make(map[string]interface{})
	schedule1Raw := make(map[string]interface{})
	if objectRaw["schedule"] != nil {
		schedule1Raw = objectRaw["schedule"].(map[string]interface{})
	}
	if len(schedule1Raw) > 0 {
		scheduleMap["cron_expression"] = schedule1Raw["cronExpression"]
		scheduleMap["delay"] = schedule1Raw["delay"]
		scheduleMap["interval"] = schedule1Raw["interval"]
		scheduleMap["run_immdiately"] = schedule1Raw["runImmediately"]
		scheduleMap["time_zone"] = schedule1Raw["timeZone"]
		scheduleMap["type"] = schedule1Raw["type"]

		scheduleMaps = append(scheduleMaps, scheduleMap)
	}
	d.Set("schedule", scheduleMaps)

	parts := strings.Split(d.Id(), ":")
	d.Set("project_name", parts[0])
	d.Set("alert_name", parts[1])

	return nil
}

func resourceAliCloudSlsAlertUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	parts := strings.Split(d.Id(), ":")
	alertName := parts[1]
	action := fmt.Sprintf("/alerts/%s", alertName)
	var err error
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})
	hostMap := make(map[string]*string)
	hostMap["project"] = StringPointer(parts[0])
	if !d.IsNewResource() && d.HasChange("display_name") {
		update = true
	}
	request["displayName"] = d.Get("display_name")
	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	request["description"] = d.Get("description")
	if d.HasChange("configuration") {
		update = true
	}
	objectDataLocalMap := make(map[string]interface{})
	if v := d.Get("configuration"); v != nil {
		nodeNative, _ := jsonpath.Get("$[0].mute_until", v)
		if nodeNative != nil && nodeNative != "" {
			objectDataLocalMap["muteUntil"] = nodeNative
		}
		nodeNative1, _ := jsonpath.Get("$[0].version", v)
		if nodeNative1 != nil && nodeNative1 != "" {
			objectDataLocalMap["version"] = nodeNative1
		}
		nodeNative2, _ := jsonpath.Get("$[0].type", v)
		if nodeNative2 != nil && nodeNative2 != "" {
			objectDataLocalMap["type"] = nodeNative2
		}
		templateConfiguration := make(map[string]interface{})
		nodeNative3, _ := jsonpath.Get("$[0].template_configuration[0].type", v)
		if nodeNative3 != nil && nodeNative3 != "" {
			templateConfiguration["type"] = nodeNative3
		}
		nodeNative4, _ := jsonpath.Get("$[0].template_configuration[0].version", v)
		if nodeNative4 != nil && nodeNative4 != "" {
			templateConfiguration["version"] = nodeNative4
		}
		nodeNative5, _ := jsonpath.Get("$[0].template_configuration[0].lang", v)
		if nodeNative5 != nil && nodeNative5 != "" {
			templateConfiguration["lang"] = nodeNative5
		}
		nodeNative6, _ := jsonpath.Get("$[0].template_configuration[0].template_id", v)
		if nodeNative6 != nil && nodeNative6 != "" {
			templateConfiguration["id"] = nodeNative6
		}
		nodeNative7, _ := jsonpath.Get("$[0].template_configuration[0].tokens", v)
		if nodeNative7 != nil && nodeNative7 != "" {
			templateConfiguration["tokens"] = nodeNative7
		}
		nodeNative8, _ := jsonpath.Get("$[0].template_configuration[0].annotations", v)
		if nodeNative8 != nil && nodeNative8 != "" {
			templateConfiguration["aonotations"] = nodeNative8
		}

		objectDataLocalMap["templateConfiguration"] = templateConfiguration
		nodeNative9, _ := jsonpath.Get("$[0].dashboard", v)
		if nodeNative9 != nil && nodeNative9 != "" {
			objectDataLocalMap["dashboard"] = nodeNative9
		}
		nodeNative10, _ := jsonpath.Get("$[0].threshold", v)
		if nodeNative10 != nil && nodeNative10 != "" {
			objectDataLocalMap["threshold"] = nodeNative10
		}
		nodeNative11, _ := jsonpath.Get("$[0].no_data_fire", v)
		if nodeNative11 != nil && nodeNative11 != "" {
			objectDataLocalMap["noDataFire"] = nodeNative11
		}
		nodeNative12, _ := jsonpath.Get("$[0].no_data_severity", v)
		if nodeNative12 != nil && nodeNative12 != "" {
			objectDataLocalMap["noDataSeverity"] = nodeNative12
		}
		nodeNative13, _ := jsonpath.Get("$[0].send_resolved", v)
		if nodeNative13 != nil && nodeNative13 != "" {
			objectDataLocalMap["sendResolved"] = nodeNative13
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData, err := jsonpath.Get("$[0].query_list", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps := make([]interface{}, 0)
			for _, dataLoop := range localData.([]interface{}) {
				dataLoopTmp := dataLoop.(map[string]interface{})
				dataLoopMap := make(map[string]interface{})
				dataLoopMap["chartTitle"] = dataLoopTmp["chart_title"]
				dataLoopMap["query"] = dataLoopTmp["query"]
				dataLoopMap["timeSpanType"] = dataLoopTmp["time_span_type"]
				dataLoopMap["start"] = dataLoopTmp["start"]
				dataLoopMap["end"] = dataLoopTmp["end"]
				dataLoopMap["storeType"] = dataLoopTmp["store_type"]
				dataLoopMap["project"] = dataLoopTmp["project"]
				dataLoopMap["store"] = dataLoopTmp["store"]
				dataLoopMap["region"] = dataLoopTmp["region"]
				dataLoopMap["roleArn"] = dataLoopTmp["role_arn"]
				dataLoopMap["powerSqlMode"] = dataLoopTmp["power_sql_mode"]
				dataLoopMap["dashboardId"] = dataLoopTmp["dashboard_id"]
				dataLoopMap["ui"] = dataLoopTmp["ui"]
				localMaps = append(localMaps, dataLoopMap)
			}
			objectDataLocalMap["queryList"] = localMaps
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData1, err := jsonpath.Get("$[0].annotations", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps1 := make([]interface{}, 0)
			for _, dataLoop1 := range localData1.([]interface{}) {
				dataLoop1Tmp := dataLoop1.(map[string]interface{})
				dataLoop1Map := make(map[string]interface{})
				dataLoop1Map["key"] = dataLoop1Tmp["key"]
				dataLoop1Map["value"] = dataLoop1Tmp["value"]
				localMaps1 = append(localMaps1, dataLoop1Map)
			}
			objectDataLocalMap["annotations"] = localMaps1
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData2, err := jsonpath.Get("$[0].labels", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps2 := make([]interface{}, 0)
			for _, dataLoop2 := range localData2.([]interface{}) {
				dataLoop2Tmp := dataLoop2.(map[string]interface{})
				dataLoop2Map := make(map[string]interface{})
				dataLoop2Map["key"] = dataLoop2Tmp["key"]
				dataLoop2Map["value"] = dataLoop2Tmp["value"]
				localMaps2 = append(localMaps2, dataLoop2Map)
			}
			objectDataLocalMap["labels"] = localMaps2
		}
		conditionConfiguration := make(map[string]interface{})
		nodeNative31, _ := jsonpath.Get("$[0].condition_configuration[0].condition", v)
		if nodeNative31 != nil && nodeNative31 != "" {
			conditionConfiguration["condition"] = nodeNative31
		}
		nodeNative32, _ := jsonpath.Get("$[0].condition_configuration[0].count_condition", v)
		if nodeNative32 != nil && nodeNative32 != "" {
			conditionConfiguration["countCondition"] = nodeNative32
		}

		objectDataLocalMap["conditionConfiguration"] = conditionConfiguration
		if v, ok := d.GetOk("configuration"); ok {
			localData3, err := jsonpath.Get("$[0].severity_configurations", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps3 := make([]interface{}, 0)
			for _, dataLoop3 := range localData3.([]interface{}) {
				dataLoop3Tmp := dataLoop3.(map[string]interface{})
				dataLoop3Map := make(map[string]interface{})
				dataLoop3Map["severity"] = dataLoop3Tmp["severity"]
				if !IsNil(dataLoop3Tmp["eval_condition"]) {
					localData4 := make(map[string]interface{})
					nodeNative34, _ := jsonpath.Get("$[0].condition", dataLoop3Tmp["eval_condition"])
					if nodeNative34 != nil && nodeNative34 != "" {
						localData4["condition"] = nodeNative34
					}
					nodeNative35, _ := jsonpath.Get("$[0].count_condition", dataLoop3Tmp["eval_condition"])
					if nodeNative35 != nil && nodeNative35 != "" {
						localData4["countCondition"] = nodeNative35
					}
					dataLoop3Map["evalCondition"] = localData4
				}
				localMaps3 = append(localMaps3, dataLoop3Map)
			}
			objectDataLocalMap["severityConfigurations"] = localMaps3
		}
		if v, ok := d.GetOk("configuration"); ok {
			localData5, err := jsonpath.Get("$[0].join_configurations", v)
			if err != nil {
				return WrapError(err)
			}
			localMaps5 := make([]interface{}, 0)
			for _, dataLoop5 := range localData5.([]interface{}) {
				dataLoop5Tmp := dataLoop5.(map[string]interface{})
				dataLoop5Map := make(map[string]interface{})
				dataLoop5Map["type"] = dataLoop5Tmp["type"]
				dataLoop5Map["condition"] = dataLoop5Tmp["condition"]
				localMaps5 = append(localMaps5, dataLoop5Map)
			}
			objectDataLocalMap["joinConfigurations"] = localMaps5
		}
		groupConfiguration := make(map[string]interface{})
		nodeNative38, _ := jsonpath.Get("$[0].group_configuration[0].type", v)
		if nodeNative38 != nil && nodeNative38 != "" {
			groupConfiguration["type"] = nodeNative38
		}
		nodeNative39, _ := jsonpath.Get("$[0].group_configuration[0].fields", d.Get("configuration"))
		if nodeNative39 != nil && nodeNative39 != "" {
			groupConfiguration["fields"] = nodeNative39
		}

		objectDataLocalMap["groupConfiguration"] = groupConfiguration
		policyConfiguration := make(map[string]interface{})
		nodeNative40, _ := jsonpath.Get("$[0].policy_configuration[0].alert_policy_id", v)
		if nodeNative40 != nil && nodeNative40 != "" {
			policyConfiguration["alertPolicyId"] = nodeNative40
		}
		nodeNative41, _ := jsonpath.Get("$[0].policy_configuration[0].action_policy_id", v)
		if nodeNative41 != nil && nodeNative41 != "" {
			policyConfiguration["actionPolicyId"] = nodeNative41
		}
		nodeNative42, _ := jsonpath.Get("$[0].policy_configuration[0].repeat_interval", v)
		if nodeNative42 != nil && nodeNative42 != "" {
			policyConfiguration["repeatInterval"] = nodeNative42
		}

		objectDataLocalMap["policyConfiguration"] = policyConfiguration
		nodeNative43, _ := jsonpath.Get("$[0].auto_annotation", v)
		if nodeNative43 != nil && nodeNative43 != "" {
			objectDataLocalMap["autoAnnotation"] = nodeNative43
		}
		sinkEventStore := make(map[string]interface{})
		nodeNative44, _ := jsonpath.Get("$[0].sink_event_store[0].enabled", v)
		if nodeNative44 != nil && nodeNative44 != "" {
			sinkEventStore["enabled"] = nodeNative44
		}
		nodeNative45, _ := jsonpath.Get("$[0].sink_event_store[0].endpoint", v)
		if nodeNative45 != nil && nodeNative45 != "" {
			sinkEventStore["endpoint"] = nodeNative45
		}
		nodeNative46, _ := jsonpath.Get("$[0].sink_event_store[0].project", v)
		if nodeNative46 != nil && nodeNative46 != "" {
			sinkEventStore["project"] = nodeNative46
		}
		nodeNative47, _ := jsonpath.Get("$[0].sink_event_store[0].event_store", v)
		if nodeNative47 != nil && nodeNative47 != "" {
			sinkEventStore["eventStore"] = nodeNative47
		}
		nodeNative48, _ := jsonpath.Get("$[0].sink_event_store[0].role_arn", v)
		if nodeNative48 != nil && nodeNative48 != "" {
			sinkEventStore["roleArn"] = nodeNative48
		}

		objectDataLocalMap["sinkEventStore"] = sinkEventStore
		sinkCms := make(map[string]interface{})
		nodeNative49, _ := jsonpath.Get("$[0].sink_cms[0].enabled", v)
		if nodeNative49 != nil && nodeNative49 != "" {
			sinkCms["enabled"] = nodeNative49
		}

		objectDataLocalMap["sinkCms"] = sinkCms
		sinkAlerthub := make(map[string]interface{})
		nodeNative50, _ := jsonpath.Get("$[0].sink_alerthub[0].enabled", v)
		if nodeNative50 != nil && nodeNative50 != "" {
			sinkAlerthub["enabled"] = nodeNative50
		}

		objectDataLocalMap["sinkAlerthub"] = sinkAlerthub
		nodeNative51, _ := jsonpath.Get("$[0].tags", d.Get("configuration"))
		if nodeNative51 != nil && nodeNative51 != "" {
			objectDataLocalMap["tags"] = nodeNative51
		}

		request["configuration"] = objectDataLocalMap
	}

	if d.HasChange("schedule") {
		update = true
	}
	objectDataLocalMap1 := make(map[string]interface{})
	if v := d.Get("schedule"); v != nil {
		nodeNative52, _ := jsonpath.Get("$[0].type", v)
		if nodeNative52 != nil && nodeNative52 != "" {
			objectDataLocalMap1["type"] = nodeNative52
		}
		nodeNative53, _ := jsonpath.Get("$[0].cron_expression", v)
		if nodeNative53 != nil && nodeNative53 != "" {
			objectDataLocalMap1["cronExpression"] = nodeNative53
		}
		nodeNative54, _ := jsonpath.Get("$[0].run_immdiately", v)
		if nodeNative54 != nil && nodeNative54 != "" {
			objectDataLocalMap1["runImmediately"] = nodeNative54
		}
		nodeNative55, _ := jsonpath.Get("$[0].time_zone", v)
		if nodeNative55 != nil && nodeNative55 != "" {
			objectDataLocalMap1["timeZone"] = nodeNative55
		}
		nodeNative56, _ := jsonpath.Get("$[0].delay", v)
		if nodeNative56 != nil && nodeNative56 != "" {
			objectDataLocalMap1["delay"] = nodeNative56
		}
		nodeNative57, _ := jsonpath.Get("$[0].interval", v)
		if nodeNative57 != nil && nodeNative57 != "" {
			objectDataLocalMap1["interval"] = nodeNative57
		}

		request["schedule"] = objectDataLocalMap1
	}

	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "UpdateAlert", action), query, body, nil, hostMap, false)
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
	}

	if d.HasChange("status") {
		client := meta.(*connectivity.AliyunClient)
		slsServiceV2 := SlsServiceV2{client}
		object, err := slsServiceV2.DescribeSlsAlert(d.Id())
		if err != nil {
			return WrapError(err)
		}

		target := d.Get("status").(string)
		if object["status"].(string) != target {
			if target == "ENABLED" {
				parts = strings.Split(d.Id(), ":")
				alertName = parts[1]
				action = fmt.Sprintf("/alerts/%s?action=enable", alertName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				hostMap := make(map[string]*string)
				hostMap["project"] = StringPointer(parts[0])
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "EnableAlert", action), query, body, nil, hostMap, false)
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

			}
			if target == "DISABLED" {
				parts = strings.Split(d.Id(), ":")
				alertName = parts[1]
				action = fmt.Sprintf("/alerts/%s?action=disable", alertName)
				request = make(map[string]interface{})
				query = make(map[string]*string)
				body = make(map[string]interface{})
				hostMap := make(map[string]*string)
				hostMap["project"] = StringPointer(parts[0])
				body = request
				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.Do("Sls", roaParam("PUT", "2020-12-30", "DisableAlert", action), query, body, nil, hostMap, false)
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

			}
		}
	}

	return resourceAliCloudSlsAlertRead(d, meta)
}

func resourceAliCloudSlsAlertDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	alertName := parts[1]
	action := fmt.Sprintf("/alerts/%s", alertName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	hostMap := make(map[string]*string)
	var err error
	request = make(map[string]interface{})
	hostMap["project"] = StringPointer(parts[0])

	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.Do("Sls", roaParam("DELETE", "2020-12-30", "DeleteAlert", action), query, body, nil, hostMap, false)
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

	return nil
}
