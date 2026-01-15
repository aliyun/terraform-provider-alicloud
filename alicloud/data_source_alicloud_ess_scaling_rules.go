package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAliCloudEssScalingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudEssScalingRulesRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cooldown": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"predictive_scaling_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"initial_max_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"predictive_value_behavior": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"predictive_value_buffer": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"predictive_task_buffer_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metric_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hybrid_monitor_namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hybrid_metrics": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metric_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"statistic": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expression": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"dimensions": {
										Type:     schema.TypeSet,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"dimension_key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"dimension_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"step_adjustment": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_interval_lower_bound": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metric_interval_upper_bound": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"scaling_adjustment": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"estimated_instance_warmup": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scale_in_evaluation_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scale_out_evaluation_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"disable_scale_in": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"adjustment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"adjustment_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_adjustment_magnitude": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scaling_rule_ari": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAliCloudEssScalingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"PageSize":   requests.NewInteger(PageSizeLarge),
		"PageNumber": requests.NewInteger(1),
		"RegionId":   client.RegionId,
	}

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request["ScalingGroupId"] = scalingGroupId.(string)
	}

	if ruleType, ok := d.GetOk("type"); ok && ruleType.(string) != "" {
		request["ScalingRuleType"] = ruleType.(string)
	}

	var allScalingRules []interface{}
	for {
		response, err = client.RpcPost("Ess", "2014-08-28", "DescribeScalingRules", nil, request, true)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_scaling_rules", "DescribeScalingRules", AlibabaCloudSdkGoERROR)
		}
		addDebug("DescribeScalingRules", response, request, request)
		w, errInfo := jsonpath.Get("$.TotalCount", response)
		if errInfo != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "$.TotalCount", response)
		}
		i, errConvert := w.(json.Number).Int64()
		if errConvert != nil {
			return WrapErrorf(err, "Convert resource %s attribute failed!!! Response: %v.", "TotalCount", response)
		}
		if int(i) < 1 {
			break
		}
		v, err := jsonpath.Get("$.ScalingRules.ScalingRule", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, "$.ScalingRules.ScalingRule", response)
		}
		if len(v.([]interface{})) < 1 {
			break
		}
		allScalingRules = append(allScalingRules, v.([]interface{})...)

		if len(v.([]interface{})) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(requests.Integer(fmt.Sprint(request["PageNumber"]))); err != nil {
			return err
		} else {
			request["PageNumber"] = page
		}
	}

	var filteredScalingRulesTemp []interface{}

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, rule := range allScalingRules {
			var object map[string]interface{}
			object = rule.(map[string]interface{})
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if r != nil && !r.MatchString(object["ScalingRuleName"].(string)) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[object["ScalingRuleId"].(string)]; !ok {
					continue
				}
			}
			filteredScalingRulesTemp = append(filteredScalingRulesTemp, rule)
		}
	} else {
		filteredScalingRulesTemp = allScalingRules
	}
	return scalingRulesDescriptionAttribute(d, filteredScalingRulesTemp, meta)
}

func scalingRulesDescriptionAttribute(d *schema.ResourceData, scalingRules []interface{}, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, scalingRule := range scalingRules {
		var object map[string]interface{}
		object = scalingRule.(map[string]interface{})
		mapping := map[string]interface{}{
			"id":               object["ScalingRuleId"],
			"scaling_group_id": object["ScalingGroupId"],
			"name":             object["ScalingRuleName"],
			"type":             object["ScalingRuleType"],
			"adjustment_type":  object["AdjustmentType"],
			"adjustment_value": object["AdjustmentValue"],
			"scaling_rule_ari": object["ScalingRuleAri"],
		}
		if object["MetricName"] != nil {
			mapping["metric_name"] = object["MetricName"]
		}
		if object["MetricType"] != nil {
			mapping["metric_type"] = object["MetricType"]
		}
		if object["StepAdjustments"] != nil {
			steps, _ := flattenStepAdjustment(object["StepAdjustments"])
			mapping["step_adjustment"] = steps
		}
		if object["EstimatedInstanceWarmup"] != nil {
			mapping["estimated_instance_warmup"] = object["EstimatedInstanceWarmup"]
		}
		if object["ScaleInEvaluationCount"] != nil {
			mapping["scale_in_evaluation_count"] = object["ScaleInEvaluationCount"]
		}
		if object["ScaleOutEvaluationCount"] != nil {
			mapping["scale_out_evaluation_count"] = object["ScaleOutEvaluationCount"]
		}
		if object["DisableScaleIn"] != nil {
			mapping["disable_scale_in"] = object["DisableScaleIn"]
		}
		if object["HybridMonitorNamespace"] != nil {
			mapping["hybrid_monitor_namespace"] = object["HybridMonitorNamespace"]
		}
		if object["HybridMetrics"] != nil {
			hybridMetrics, _ := flattenHybridMetricsMappings(object["HybridMetrics"])
			mapping["hybrid_metrics"] = hybridMetrics
		}
		if object["TargetValue"] != nil {
			mapping["target_value"] = object["TargetValue"]
		}
		if object["PredictiveTaskBufferTime"] != nil {
			mapping["predictive_task_buffer_time"] = object["PredictiveTaskBufferTime"]
		}
		if object["PredictiveValueBuffer"] != nil {
			mapping["predictive_value_buffer"] = object["PredictiveValueBuffer"]
		}
		if object["PredictiveValueBehavior"] != nil {
			mapping["predictive_value_behavior"] = object["PredictiveValueBehavior"]
		}
		if object["InitialMaxSize"] != nil {
			mapping["initial_max_size"] = object["InitialMaxSize"]
		}
		if object["PredictiveScalingMode"] != nil {
			mapping["predictive_scaling_mode"] = object["PredictiveScalingMode"]
		}
		if object["Cooldown"] != nil {
			mapping["cooldown"] = object["Cooldown"]
		}
		if object["MinAdjustmentMagnitude"] != nil {
			mapping["min_adjustment_magnitude"] = object["MinAdjustmentMagnitude"]
		}
		ids = append(ids, object["ScalingRuleId"].(string))
		names = append(names, object["ScalingRuleName"].(string))
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}

func flattenStepAdjustment(list interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, i := range list.(map[string]interface{})["StepAdjustment"].([]interface{}) {
		stepAdjustment := i.(map[string]interface{})
		l := map[string]interface{}{}
		if stepAdjustment["MetricIntervalLowerBound"] != nil {
			numberLower := stepAdjustment["MetricIntervalLowerBound"].(json.Number)
			l["metric_interval_lower_bound"] = numberLower.String()
		}
		if stepAdjustment["MetricIntervalUpperBound"] != nil {
			numberUpper := stepAdjustment["MetricIntervalUpperBound"].(json.Number)
			l["metric_interval_upper_bound"] = numberUpper.String()
		}
		l["scaling_adjustment"] = stepAdjustment["ScalingAdjustment"]
		result = append(result, l)
	}
	return result, nil
}
