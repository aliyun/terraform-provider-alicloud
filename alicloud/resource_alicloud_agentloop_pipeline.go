// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudAgentloopPipeline() *schema.Resource {
	datasetSinkSchema := func() *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_space": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"dataset": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
			},
		}
	}
	sinkTargetSchema := func() *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"dataset": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					MaxItems: 1,
					Elem:     datasetSinkSchema(),
				},
			},
		}
	}
	return &schema.Resource{
		Create: resourceAliCloudAgentloopPipelineCreate,
		Read:   resourceAliCloudAgentloopPipelineRead,
		Update: resourceAliCloudAgentloopPipelineUpdate,
		Delete: resourceAliCloudAgentloopPipelineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"agent_space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execute_policy": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"run_once": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"from_time": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"to_time": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"scheduled": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"interval": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"from_time": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"pipeline": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"nodes": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"parameters": {
										Type:     schema.TypeMap,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"pipeline_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sink": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"dataset": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem:     datasetSinkSchema(),
						},
						"condition": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"match_mode": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"routes": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"expression": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"sink": {
													Type:     schema.TypeList,
													Optional: true,
													ForceNew: true,
													MaxItems: 1,
													Elem:     sinkTargetSchema(),
												},
											},
										},
									},
									"default_sink": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										MaxItems: 1,
										Elem:     sinkTargetSchema(),
									},
								},
							},
						},
					},
				},
			},
			"source": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"input_fields": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"logstore": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"logstore": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"query": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
						"dataset": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dataset": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"filter": {
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
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func expandAgentloopPipelineDatasetSink(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawList, ok := raw.([]interface{})
	if !ok || len(rawList) == 0 || rawList[0] == nil {
		return result
	}
	rawMap, ok := rawList[0].(map[string]interface{})
	if !ok {
		return result
	}
	if v, ok := rawMap["agent_space"].(string); ok && v != "" {
		result["agentSpace"] = v
	}
	if v, ok := rawMap["dataset"].(string); ok && v != "" {
		result["dataset"] = v
	}
	return result
}

func expandAgentloopPipelineSinkTarget(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawList, ok := raw.([]interface{})
	if !ok || len(rawList) == 0 || rawList[0] == nil {
		return result
	}
	rawMap, ok := rawList[0].(map[string]interface{})
	if !ok {
		return result
	}
	if v, ok := rawMap["type"].(string); ok && v != "" {
		result["type"] = v
	}
	if dataset := expandAgentloopPipelineDatasetSink(rawMap["dataset"]); len(dataset) > 0 {
		result["dataset"] = dataset
	}
	return result
}

func expandAgentloopPipelineSink(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawList, ok := raw.([]interface{})
	if !ok || len(rawList) == 0 || rawList[0] == nil {
		return result
	}
	rawMap, ok := rawList[0].(map[string]interface{})
	if !ok {
		return result
	}
	if v, ok := rawMap["type"].(string); ok && v != "" {
		result["type"] = v
	}
	if dataset := expandAgentloopPipelineDatasetSink(rawMap["dataset"]); len(dataset) > 0 {
		result["dataset"] = dataset
	}
	conditionList, ok := rawMap["condition"].([]interface{})
	if ok && len(conditionList) > 0 && conditionList[0] != nil {
		if conditionMap, ok := conditionList[0].(map[string]interface{}); ok {
			condition := make(map[string]interface{})
			if v, ok := conditionMap["match_mode"].(string); ok && v != "" {
				condition["matchMode"] = v
			}
			if routesRaw, ok := conditionMap["routes"].([]interface{}); ok {
				routes := make([]interface{}, 0)
				for _, routeRaw := range routesRaw {
					routeMap, ok := routeRaw.(map[string]interface{})
					if !ok {
						continue
					}
					route := make(map[string]interface{})
					if v, ok := routeMap["id"].(string); ok && v != "" {
						route["id"] = v
					}
					if v, ok := routeMap["expression"].(string); ok && v != "" {
						route["expression"] = v
					}
					if sink := expandAgentloopPipelineSinkTarget(routeMap["sink"]); len(sink) > 0 {
						route["sink"] = sink
					}
					routes = append(routes, route)
				}
				if len(routes) > 0 {
					condition["routes"] = routes
				}
			}
			if defaultSink := expandAgentloopPipelineSinkTarget(conditionMap["default_sink"]); len(defaultSink) > 0 {
				condition["defaultSink"] = defaultSink
			}
			if len(condition) > 0 {
				result["condition"] = condition
			}
		}
	}
	return result
}

func expandAgentloopPipelineSource(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawList, ok := raw.([]interface{})
	if !ok || len(rawList) == 0 || rawList[0] == nil {
		return result
	}
	rawMap, ok := rawList[0].(map[string]interface{})
	if !ok {
		return result
	}
	if v, ok := rawMap["type"].(string); ok && v != "" {
		result["type"] = v
	}
	if inputFieldsRaw, ok := rawMap["input_fields"].([]interface{}); ok {
		inputFields := make([]interface{}, 0)
		for _, fieldRaw := range inputFieldsRaw {
			fieldMap, ok := fieldRaw.(map[string]interface{})
			if !ok {
				continue
			}
			field := make(map[string]interface{})
			if v, ok := fieldMap["name"].(string); ok && v != "" {
				field["name"] = v
			}
			if v, ok := fieldMap["type"].(string); ok && v != "" {
				field["type"] = v
			}
			inputFields = append(inputFields, field)
		}
		if len(inputFields) > 0 {
			result["inputFields"] = inputFields
		}
	}
	if logstoreRaw, ok := rawMap["logstore"].([]interface{}); ok && len(logstoreRaw) > 0 && logstoreRaw[0] != nil {
		if logstoreMap, ok := logstoreRaw[0].(map[string]interface{}); ok {
			logstore := make(map[string]interface{})
			if v, ok := logstoreMap["project"].(string); ok && v != "" {
				logstore["project"] = v
			}
			if v, ok := logstoreMap["logstore"].(string); ok && v != "" {
				logstore["logstore"] = v
			}
			if v, ok := logstoreMap["query"].(string); ok && v != "" {
				logstore["query"] = v
			}
			if len(logstore) > 0 {
				result["logstore"] = logstore
			}
		}
	}
	if datasetRaw, ok := rawMap["dataset"].([]interface{}); ok && len(datasetRaw) > 0 && datasetRaw[0] != nil {
		if datasetMap, ok := datasetRaw[0].(map[string]interface{}); ok {
			dataset := make(map[string]interface{})
			if v, ok := datasetMap["dataset"].(string); ok && v != "" {
				dataset["dataset"] = v
			}
			if v, ok := datasetMap["filter"].(string); ok && v != "" {
				dataset["filter"] = v
			}
			if len(dataset) > 0 {
				result["dataset"] = dataset
			}
		}
	}
	return result
}

func expandAgentloopPipelineExecutePolicy(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawList, ok := raw.([]interface{})
	if !ok || len(rawList) == 0 || rawList[0] == nil {
		return result
	}
	rawMap, ok := rawList[0].(map[string]interface{})
	if !ok {
		return result
	}
	if v, ok := rawMap["mode"].(string); ok && v != "" {
		result["mode"] = v
	}
	if runOnceRaw, ok := rawMap["run_once"].([]interface{}); ok && len(runOnceRaw) > 0 && runOnceRaw[0] != nil {
		if runOnceMap, ok := runOnceRaw[0].(map[string]interface{}); ok {
			runOnce := make(map[string]interface{})
			if v, ok := runOnceMap["from_time"]; ok && v != nil {
				runOnce["fromTime"] = v
			}
			if v, ok := runOnceMap["to_time"]; ok && v != nil {
				runOnce["toTime"] = v
			}
			if len(runOnce) > 0 {
				result["runOnce"] = runOnce
			}
		}
	}
	if scheduledRaw, ok := rawMap["scheduled"].([]interface{}); ok && len(scheduledRaw) > 0 && scheduledRaw[0] != nil {
		if scheduledMap, ok := scheduledRaw[0].(map[string]interface{}); ok {
			scheduled := make(map[string]interface{})
			if v, ok := scheduledMap["interval"].(string); ok && v != "" {
				scheduled["interval"] = v
			}
			if v, ok := scheduledMap["from_time"]; ok && v != nil {
				scheduled["fromTime"] = v
			}
			if len(scheduled) > 0 {
				result["scheduled"] = scheduled
			}
		}
	}
	return result
}

func expandAgentloopPipelinePipeline(raw interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	rawList, ok := raw.([]interface{})
	if !ok || len(rawList) == 0 || rawList[0] == nil {
		return result
	}
	rawMap, ok := rawList[0].(map[string]interface{})
	if !ok {
		return result
	}
	nodes := make([]interface{}, 0)
	if nodesRaw, ok := rawMap["nodes"].([]interface{}); ok {
		for _, nodeRaw := range nodesRaw {
			nodeMap, ok := nodeRaw.(map[string]interface{})
			if !ok {
				continue
			}
			node := make(map[string]interface{})
			if v, ok := nodeMap["id"].(string); ok && v != "" {
				node["id"] = v
			}
			if v, ok := nodeMap["type"].(string); ok && v != "" {
				node["type"] = v
			}
			if v, ok := nodeMap["parameters"].(map[string]interface{}); ok && len(v) > 0 {
				node["parameters"] = v
			}
			nodes = append(nodes, node)
		}
	}
	result["nodes"] = nodes
	return result
}

func flattenAgentloopPipelineDatasetSink(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	rawMap, ok := raw.(map[string]interface{})
	if !ok || len(rawMap) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["agent_space"] = rawMap["agentSpace"]
	item["dataset"] = rawMap["dataset"]
	return append(result, item)
}

func flattenAgentloopPipelineSinkTarget(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	rawMap, ok := raw.(map[string]interface{})
	if !ok || len(rawMap) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["type"] = rawMap["type"]
	item["dataset"] = flattenAgentloopPipelineDatasetSink(rawMap["dataset"])
	return append(result, item)
}

func flattenAgentloopPipelineSink(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	rawMap, ok := raw.(map[string]interface{})
	if !ok || len(rawMap) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["type"] = rawMap["type"]
	item["dataset"] = flattenAgentloopPipelineDatasetSink(rawMap["dataset"])
	conditionMaps := make([]map[string]interface{}, 0)
	if conditionRaw, ok := rawMap["condition"].(map[string]interface{}); ok && len(conditionRaw) > 0 {
		conditionMap := make(map[string]interface{})
		conditionMap["match_mode"] = conditionRaw["matchMode"]
		routesMaps := make([]map[string]interface{}, 0)
		if routesRaw, ok := conditionRaw["routes"].([]interface{}); ok {
			for _, routeRaw := range routesRaw {
				routeMapRaw, ok := routeRaw.(map[string]interface{})
				if !ok {
					continue
				}
				routeMap := make(map[string]interface{})
				routeMap["id"] = routeMapRaw["id"]
				routeMap["expression"] = routeMapRaw["expression"]
				routeMap["sink"] = flattenAgentloopPipelineSinkTarget(routeMapRaw["sink"])
				routesMaps = append(routesMaps, routeMap)
			}
		}
		conditionMap["routes"] = routesMaps
		conditionMap["default_sink"] = flattenAgentloopPipelineSinkTarget(conditionRaw["defaultSink"])
		conditionMaps = append(conditionMaps, conditionMap)
	}
	item["condition"] = conditionMaps
	return append(result, item)
}

func flattenAgentloopPipelineSource(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	rawMap, ok := raw.(map[string]interface{})
	if !ok || len(rawMap) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["type"] = rawMap["type"]
	inputFieldsMaps := make([]map[string]interface{}, 0)
	if inputFieldsRaw, ok := rawMap["inputFields"].([]interface{}); ok {
		for _, fieldRaw := range inputFieldsRaw {
			fieldMapRaw, ok := fieldRaw.(map[string]interface{})
			if !ok {
				continue
			}
			fieldMap := make(map[string]interface{})
			fieldMap["name"] = fieldMapRaw["name"]
			fieldMap["type"] = fieldMapRaw["type"]
			inputFieldsMaps = append(inputFieldsMaps, fieldMap)
		}
	}
	item["input_fields"] = inputFieldsMaps
	logstoreMaps := make([]map[string]interface{}, 0)
	if logstoreRaw, ok := rawMap["logstore"].(map[string]interface{}); ok && len(logstoreRaw) > 0 {
		logstoreMap := make(map[string]interface{})
		logstoreMap["project"] = logstoreRaw["project"]
		logstoreMap["logstore"] = logstoreRaw["logstore"]
		logstoreMap["query"] = logstoreRaw["query"]
		logstoreMaps = append(logstoreMaps, logstoreMap)
	}
	item["logstore"] = logstoreMaps
	datasetMaps := make([]map[string]interface{}, 0)
	if datasetRaw, ok := rawMap["dataset"].(map[string]interface{}); ok && len(datasetRaw) > 0 {
		datasetMap := make(map[string]interface{})
		datasetMap["dataset"] = datasetRaw["dataset"]
		datasetMap["filter"] = datasetRaw["filter"]
		datasetMaps = append(datasetMaps, datasetMap)
	}
	item["dataset"] = datasetMaps
	return append(result, item)
}

func flattenAgentloopPipelineExecutePolicy(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	rawMap, ok := raw.(map[string]interface{})
	if !ok || len(rawMap) == 0 {
		return result
	}
	item := make(map[string]interface{})
	item["mode"] = rawMap["mode"]
	runOnceMaps := make([]map[string]interface{}, 0)
	if runOnceRaw, ok := rawMap["runOnce"].(map[string]interface{}); ok && len(runOnceRaw) > 0 {
		runOnceMap := make(map[string]interface{})
		runOnceMap["from_time"] = runOnceRaw["fromTime"]
		runOnceMap["to_time"] = runOnceRaw["toTime"]
		runOnceMaps = append(runOnceMaps, runOnceMap)
	}
	item["run_once"] = runOnceMaps
	scheduledMaps := make([]map[string]interface{}, 0)
	if scheduledRaw, ok := rawMap["scheduled"].(map[string]interface{}); ok && len(scheduledRaw) > 0 {
		scheduledMap := make(map[string]interface{})
		scheduledMap["interval"] = scheduledRaw["interval"]
		scheduledMap["from_time"] = scheduledRaw["fromTime"]
		scheduledMaps = append(scheduledMaps, scheduledMap)
	}
	item["scheduled"] = scheduledMaps
	return append(result, item)
}

func flattenAgentloopPipelinePipeline(raw interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	rawMap, ok := raw.(map[string]interface{})
	if !ok || len(rawMap) == 0 {
		return result
	}
	item := make(map[string]interface{})
	nodesMaps := make([]map[string]interface{}, 0)
	if nodesRaw, ok := rawMap["nodes"].([]interface{}); ok {
		for _, nodeRaw := range nodesRaw {
			nodeMapRaw, ok := nodeRaw.(map[string]interface{})
			if !ok {
				continue
			}
			nodeMap := make(map[string]interface{})
			nodeMap["id"] = nodeMapRaw["id"]
			nodeMap["type"] = nodeMapRaw["type"]
			nodeMap["parameters"] = nodeMapRaw["parameters"]
			nodesMaps = append(nodesMaps, nodeMap)
		}
	}
	item["nodes"] = nodesMaps
	return append(result, item)
}

func resourceAliCloudAgentloopPipelineCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	agentSpace := d.Get("agent_space")
	action := fmt.Sprintf("/agentspace/%s/pipeline", agentSpace)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	body := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["pipelineName"] = d.Get("pipeline_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("pipeline"); ok {
		request["pipeline"] = expandAgentloopPipelinePipeline(v)
	}
	if v, ok := d.GetOk("source"); ok {
		if source := expandAgentloopPipelineSource(v); len(source) > 0 {
			request["source"] = source
		}
	}
	if v, ok := d.GetOk("sink"); ok {
		if sink := expandAgentloopPipelineSink(v); len(sink) > 0 {
			request["sink"] = sink
		}
	}
	if v, ok := d.GetOk("execute_policy"); ok {
		if executePolicy := expandAgentloopPipelineExecutePolicy(v); len(executePolicy) > 0 {
			request["executePolicy"] = executePolicy
		}
	}
	body = request
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RoaPost("AgentLoop", "2026-05-20", action, query, nil, body, true)
		if err != nil {
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_agentloop_pipeline", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", agentSpace, request["pipelineName"]))

	return resourceAliCloudAgentloopPipelineRead(d, meta)
}

func resourceAliCloudAgentloopPipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	agentloopServiceV2 := AgentloopServiceV2{client}

	objectRaw, err := agentloopServiceV2.DescribeAgentloopPipeline(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_agentloop_pipeline DescribeAgentloopPipeline Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("create_time", objectRaw["createTime"])
	d.Set("description", objectRaw["description"])
	d.Set("pipeline_name", objectRaw["pipelineName"])
	d.Set("region_id", objectRaw["regionId"])
	d.Set("update_time", objectRaw["updateTime"])

	if err := d.Set("pipeline", flattenAgentloopPipelinePipeline(objectRaw["pipeline"])); err != nil {
		return err
	}
	if err := d.Set("source", flattenAgentloopPipelineSource(objectRaw["source"])); err != nil {
		return err
	}
	if err := d.Set("sink", flattenAgentloopPipelineSink(objectRaw["sink"])); err != nil {
		return err
	}
	if err := d.Set("execute_policy", flattenAgentloopPipelineExecutePolicy(objectRaw["executePolicy"])); err != nil {
		return err
	}

	parts := strings.Split(d.Id(), ":")
	d.Set("agent_space", parts[0])

	return nil
}

func resourceAliCloudAgentloopPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]*string
	var body map[string]interface{}
	update := false
	d.Partial(true)

	var err error
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	pipelineName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/pipeline/%s", agentSpace, pipelineName)
	request = make(map[string]interface{})
	query = make(map[string]*string)
	body = make(map[string]interface{})

	if d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok || d.HasChange("description") {
		request["description"] = v
	}
	// NOTE: the backend only allows updating description. Sending pipeline,
	// source, sink or executePolicy in the update request is rejected with
	// "Only 'description' is mutable" / "executePolicy.mode is immutable", so
	// those attributes are marked ForceNew and never sent here.
	body = request
	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RoaPut("AgentLoop", "2026-05-20", action, query, nil, body, true)
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
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

	d.Partial(false)
	return resourceAliCloudAgentloopPipelineRead(d, meta)
}

func resourceAliCloudAgentloopPipelineDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	agentSpace := parts[0]
	pipelineName := parts[1]
	action := fmt.Sprintf("/agentspace/%s/pipeline/%s", agentSpace, pipelineName)
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]*string)
	var err error
	request = make(map[string]interface{})

	wait := incrementalWait(3*time.Second, 5*time.Second)
	// A pipeline whose scheduleStatus is Active cannot be deleted; the backend
	// answers 409 InvalidOperation "Cannot delete an Active pipeline. Pause it
	// first.". When that happens, pause the pipeline through the dedicated
	// PausePipeline API (POST .../pause, idempotent) and keep retrying the
	// delete until the schedule change takes effect.
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RoaDelete("AgentLoop", "2026-05-20", action, query, nil, nil, true)
		if err != nil {
			if IsExpectedErrors(err, []string{"InvalidOperation"}) && strings.Contains(err.Error(), "Pause it first") {
				pauseAction := fmt.Sprintf("%s/pause", action)
				pauseBody := map[string]interface{}{"reason": "terraform-provider-alicloud: delete pipeline"}
				if _, pauseErr := client.RoaPost("AgentLoop", "2026-05-20", pauseAction, query, nil, pauseBody, true); pauseErr != nil {
					log.Printf("[DEBUG] Pause pipeline %s before delete failed: %v", d.Id(), pauseErr)
				}
				wait()
				return resource.RetryableError(err)
			}
			if NeedRetry(err) || IsExpectedErrors(err, []string{"Conflict."}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrors(err, []string{"PipelineNotExist"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}
