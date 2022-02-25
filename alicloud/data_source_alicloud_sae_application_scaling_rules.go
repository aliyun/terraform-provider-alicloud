package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudSaeApplicationScalingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudSaeApplicationScalingRulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_rule_enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"scaling_rule_metric": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"max_replicas": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"min_replicas": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"metrics": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_target_average_utilization": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"metric_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"metrics_status": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"desired_replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"next_scale_time_period": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"current_replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"last_scale_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"max_replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min_replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"current_metrics": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"current_value": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
												"next_scale_metrics": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"next_scale_out_average_utilization": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"next_scale_in_average_utilization": {
																Type:     schema.TypeInt,
																Computed: true,
															},
															"name": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"scale_down_rules": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"stabilization_window_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"step": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
									"scale_up_rules": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"disabled": {
													Type:     schema.TypeBool,
													Computed: true,
												},
												"stabilization_window_seconds": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"step": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"scaling_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_rule_timer": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_date": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"period": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"schedules": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"at_time": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"target_replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"max_replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"min_replicas": {
													Type:     schema.TypeInt,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"scaling_rule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudSaeApplicationScalingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "/pop/v1/sam/scale/applicationScalingRules"
	request := make(map[string]*string)
	var objects []map[string]interface{}
	request["AppId"] = StringPointer(d.Get("app_id").(string))
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("GET"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_sae_application_scaling_rules", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "GET "+action, response))
	}
	resp, err := jsonpath.Get("$.Data.ApplicationScalingRules", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.ApplicationScalingRules", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["AppId"], ":", item["ScaleRuleName"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                  fmt.Sprint(object["AppId"], ":", object["ScaleRuleName"]),
			"app_id":              fmt.Sprint(object["AppId"]),
			"create_time":         fmt.Sprint(object["CreateTime"]),
			"scaling_rule_enable": object["ScaleRuleEnabled"],
			"scaling_rule_name":   object["ScaleRuleName"],
			"scaling_rule_type":   object["ScaleRuleType"],
		}
		scalingRuleTimerMaps := make([]map[string]interface{}, 0)
		if ScalingRuleTimer, ok := object["Timer"]; ok {
			scalingRuleTimerArg := ScalingRuleTimer.(map[string]interface{})
			if len(scalingRuleTimerArg) > 0 {
				scalingRuleTimerMap := make(map[string]interface{}, 0)
				if v, ok := scalingRuleTimerArg["EndDate"]; ok {
					scalingRuleTimerMap["end_date"] = v
				}
				if v, ok := scalingRuleTimerArg["BeginDate"]; ok {
					scalingRuleTimerMap["begin_date"] = v
				}
				scalingRuleTimerMap["period"] = scalingRuleTimerArg["Period"]
				schedulesMaps := make([]map[string]interface{}, 0)
				if schedules, ok := scalingRuleTimerArg["Schedules"]; ok {
					schedulesArg := schedules.([]interface{})
					for _, v := range schedulesArg {
						serverGroupTuplesArg := v.(map[string]interface{})
						schedulesMap := map[string]interface{}{}
						schedulesMap["at_time"] = serverGroupTuplesArg["AtTime"]
						schedulesMap["target_replicas"] = serverGroupTuplesArg["TargetReplicas"]
						schedulesMap["max_replicas"] = serverGroupTuplesArg["MaxReplicas"]
						schedulesMap["min_replicas"] = serverGroupTuplesArg["MinReplicas"]
						schedulesMaps = append(schedulesMaps, schedulesMap)
					}
				}
				scalingRuleTimerMap["schedules"] = schedulesMaps
				scalingRuleTimerMaps = append(scalingRuleTimerMaps, scalingRuleTimerMap)
				mapping["scaling_rule_timer"] = scalingRuleTimerMaps
			}
		}
		scalingRuleMetricMaps := make([]map[string]interface{}, 0)
		if scalingRuleMetric, ok := object["Metric"]; ok {
			scalingRuleMetricArg := scalingRuleMetric.(map[string]interface{})
			scalingRuleMetricMap := make(map[string]interface{}, 0)
			if v, ok := scalingRuleMetricArg["MaxReplicas"]; ok {
				scalingRuleMetricMap["max_replicas"] = v
			}
			if v, ok := scalingRuleMetricArg["MinReplicas"]; ok {
				scalingRuleMetricMap["min_replicas"] = v
			}
			metricsMaps := make([]map[string]interface{}, 0)
			if v, ok := scalingRuleMetricArg["Metrics"]; ok {
				metrics := v.([]interface{})
				for _, v := range metrics {
					metricsArg := v.(map[string]interface{})
					metricsMap := map[string]interface{}{}
					metricsMap["metric_type"] = metricsArg["MetricType"]
					metricsMap["metric_target_average_utilization"] = metricsArg["MetricTargetAverageUtilization"]
					metricsMaps = append(metricsMaps, metricsMap)
				}
				scalingRuleMetricMap["metrics"] = metricsMaps
			}

			metricsStatusMaps := make([]map[string]interface{}, 0)
			if v, ok := scalingRuleMetricArg["MetricsStatus"]; ok {
				metricsStatusArg := v.(map[string]interface{})
				metricsStatusMap := map[string]interface{}{}
				metricsStatusMap["desired_replicas"] = metricsStatusArg["DesiredReplicas"]
				metricsStatusMap["next_scale_time_period"] = metricsStatusArg["NextScaleTimePeriod"]
				metricsStatusMap["current_replicas"] = metricsStatusArg["CurrentReplicas"]
				metricsStatusMap["last_scale_time"] = metricsStatusArg["LastScaleTime"]
				metricsStatusMap["max_replicas"] = metricsStatusArg["MaxReplicas"]
				metricsStatusMap["min_replicas"] = metricsStatusArg["MinReplicas"]

				currentMetricsMaps := make([]map[string]interface{}, 0)
				if v, ok := metricsStatusArg["CurrentMetrics"]; ok {
					currentMetrics := v.([]interface{})
					for _, v := range currentMetrics {
						currentMetricsArg := v.(map[string]interface{})
						currentMetricsMap := map[string]interface{}{}
						currentMetricsMap["type"] = currentMetricsArg["Type"]
						currentMetricsMap["current_value"] = currentMetricsArg["CurrentValue"]
						currentMetricsMap["name"] = currentMetricsArg["Name"]

						currentMetricsMaps = append(currentMetricsMaps, currentMetricsMap)
					}
					metricsStatusMap["current_metrics"] = currentMetricsMaps
				}

				nextScaleMetricsMaps := make([]map[string]interface{}, 0)
				if v, ok := metricsStatusArg["NextScaleMetrics"]; ok {
					nextScaleMetrics := v.([]interface{})
					for _, v := range nextScaleMetrics {
						nextScaleMetricsArg := v.(map[string]interface{})
						nextScaleMetricsMap := map[string]interface{}{}
						nextScaleMetricsMap["next_scale_out_average_utilization"] = nextScaleMetricsArg["NextScaleOutAverageUtilization"]
						nextScaleMetricsMap["next_scale_in_average_utilization"] = nextScaleMetricsArg["NextScaleInAverageUtilization"]
						nextScaleMetricsMap["name"] = nextScaleMetricsArg["Name"]

						nextScaleMetricsMaps = append(nextScaleMetricsMaps, nextScaleMetricsMap)
					}
					metricsStatusMap["next_scale_metrics"] = nextScaleMetricsMaps
				}

				metricsStatusMaps = append(metricsStatusMaps, metricsStatusMap)
				scalingRuleMetricMap["metrics_status"] = metricsStatusMaps
			}

			scaleUpRulesMaps := make([]map[string]interface{}, 0)
			if scaleUpRules, ok := scalingRuleMetricArg["ScaleUpRules"]; ok {
				scaleUpRulesArg := scaleUpRules.(map[string]interface{})
				if len(scaleUpRulesArg) > 0 {
					scaleUpRulesMap := map[string]interface{}{}
					scaleUpRulesMap["step"] = formatInt(scaleUpRulesArg["Step"])
					scaleUpRulesMap["disabled"] = scaleUpRulesArg["Disabled"]
					scaleUpRulesMap["stabilization_window_seconds"] = scaleUpRulesArg["StabilizationWindowSeconds"]
					scaleUpRulesMaps = append(scaleUpRulesMaps, scaleUpRulesMap)
					scalingRuleMetricMap["scale_up_rules"] = scaleUpRulesMaps
				}
			}

			scaleDownRulesMaps := make([]map[string]interface{}, 0)
			if scaleDownRules, ok := scalingRuleMetricArg["ScaleDownRules"]; ok {
				scaleDownRulesArg := scaleDownRules.(map[string]interface{})
				if len(scaleDownRulesArg) > 0 {
					scaleDownRulesMap := map[string]interface{}{}
					scaleDownRulesMap["step"] = formatInt(scaleDownRulesArg["Step"])
					scaleDownRulesMap["disabled"] = scaleDownRulesArg["Disabled"]
					scaleDownRulesMap["stabilization_window_seconds"] = scaleDownRulesArg["StabilizationWindowSeconds"]
					scaleDownRulesMaps = append(scaleDownRulesMaps, scaleDownRulesMap)
					scalingRuleMetricMap["scale_down_rules"] = scaleDownRulesMaps
				}
			}
			if len(scalingRuleMetricMap) > 0 {
				scalingRuleMetricMaps = append(scalingRuleMetricMaps, scalingRuleMetricMap)
				mapping["scaling_rule_metric"] = scalingRuleMetricMaps
			}
		}

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
