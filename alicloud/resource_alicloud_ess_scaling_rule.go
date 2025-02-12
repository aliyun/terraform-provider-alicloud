package alicloud

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudEssScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingRuleCreate,
		Read:   resourceAliyunEssScalingRuleRead,
		Update: resourceAliyunEssScalingRuleUpdate,
		Delete: resourceAliyunEssScalingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"adjustment_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: StringInSlice([]string{"QuantityChangeInCapacity", "PercentChangeInCapacity", "TotalCapacity"}, false),
			},
			"adjustment_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scaling_rule_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"ari": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cooldown": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
			"scaling_rule_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SimpleScalingRule",
				ForceNew: true,
				ValidateFunc: StringInSlice([]string{
					string(SimpleScalingRule),
					string(TargetTrackingScalingRule),
					string(StepScalingRule),
					string(PredictiveScalingRule),
				}, false),
			},
			"estimated_instance_warmup": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"scale_in_evaluation_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"scale_out_evaluation_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"min_adjustment_magnitude": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if v, ok := d.GetOk("scaling_rule_type"); ok && (v.(string) == "" || v.(string) == "SimpleScalingRule" || v.(string) == "StepScalingRule") {
						if w, ok := d.GetOk("adjustment_type"); ok && w.(string) == "PercentChangeInCapacity" {
							return false
						}
					}
					return true
				},
			},
			"metric_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"predictive_scaling_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: StringInSlice([]string{
					"PredictOnly",
					"PredictAndScale",
				}, false),
			},
			"initial_max_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"predictive_value_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: StringInSlice([]string{
					"MaxOverridePredictiveValue",
					"PredictiveValueOverrideMax",
					"PredictiveValueOverrideMaxWithBuffer",
				}, false),
			},
			"predictive_value_buffer": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 100),
			},
			"predictive_task_buffer_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: IntBetween(0, 60),
			},
			"target_value": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"hybrid_monitor_namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hybrid_metrics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"statistic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expression": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dimensions": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dimension_key": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"dimension_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"disable_scale_in": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"step_adjustment": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_interval_lower_bound": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metric_interval_upper_bound": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"scaling_adjustment": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"alarm_dimension": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimension_key": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: StringInSlice([]string{
								"rulePool",
							}, false),
						},
						"dimension_value": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliyunEssScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	request, err := buildAlicloudEssScalingRuleArgs(d, meta)
	if err != nil {
		return WrapError(err)
	}
	client := meta.(*connectivity.AliyunClient)

	raw, err := client.RpcPost("Ess", "2014-08-28", "CreateScalingRule", nil, request, true)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scalingrule", "CreateScalingRule", AlibabaCloudSdkGoERROR)
	}
	addDebug("CreateScalingRule", raw, request, request)
	d.SetId(fmt.Sprint(raw["ScalingRuleId"]))
	return resourceAliyunEssScalingRuleRead(d, meta)
}

func resourceAliyunEssScalingRuleRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	_, ok := d.GetOk("alarm_dimension")
	var object map[string]interface{}
	var err error
	if !ok {
		object, err = essService.DescribeEssScalingRule(d.Id())
	} else {
		object, err = essService.DescribeEssScalingRuleWithAlarm(d.Id())
	}
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("scaling_group_id", object["ScalingGroupId"])
	d.Set("ari", object["ScalingRuleAri"])
	d.Set("adjustment_type", object["AdjustmentType"])
	if object["AdjustmentValue"] != nil {
		d.Set("adjustment_value", object["AdjustmentValue"])
	}
	d.Set("scaling_rule_name", object["ScalingRuleName"])
	if object["Cooldown"] != nil {
		d.Set("cooldown", object["Cooldown"])
	}
	d.Set("scaling_rule_type", object["ScalingRuleType"])
	if object["EstimatedInstanceWarmup"] != nil {
		d.Set("estimated_instance_warmup", object["EstimatedInstanceWarmup"])
	}
	if object["ScaleInEvaluationCount"] != nil {
		d.Set("scale_in_evaluation_count", object["ScaleInEvaluationCount"])
	}
	if object["ScaleOutEvaluationCount"] != nil {
		d.Set("scale_out_evaluation_count", object["ScaleOutEvaluationCount"])
	}
	if object["MinAdjustmentMagnitude"] != nil {
		d.Set("min_adjustment_magnitude", object["MinAdjustmentMagnitude"])
	}

	d.Set("predictive_scaling_mode", object["PredictiveScalingMode"])
	if object["InitialMaxSize"] != nil {
		d.Set("initial_max_size", object["InitialMaxSize"])
	}
	d.Set("predictive_value_behavior", object["PredictiveValueBehavior"])
	d.Set("predictive_value_buffer", object["PredictiveValueBuffer"])
	if object["PredictiveTaskBufferTime"] != nil {
		d.Set("predictive_task_buffer_time", object["PredictiveTaskBufferTime"])
	}
	if object["TargetValue"] != nil {
		f, _ := object["TargetValue"].(json.Number).Float64()
		targetValue, err := strconv.ParseFloat(strconv.FormatFloat(f, 'f', 3, 64), 64)
		if err != nil {
			return WrapError(err)
		}
		d.Set("target_value", targetValue)
	}

	d.Set("disable_scale_in", object["DisableScaleIn"])
	if object["StepAdjustment"] != nil {
		steps, _ := flattenStepAdjustmentMappings(object["StepAdjustments"])
		d.Set("step_adjustment", steps)
	}
	if object["AlarmDimensions"] != nil {
		alarmDimensions := flattenAlarmDimensionMappings(object["AlarmDimensions"])
		d.Set("alarm_dimension", alarmDimensions)
	}
	if object["HybridMonitorNamespace"] != nil {
		d.Set("hybrid_monitor_namespace", object["HybridMonitorNamespace"])
	}

	if object["HybridMetrics"] != nil {
		dimensions, _ := flattenHybridMetricsMappings(object["HybridMetrics"])
		d.Set("hybrid_metrics", dimensions)
	}
	if object["MetricName"] != nil && object["HybridMetrics"] == nil {
		d.Set("metric_name", object["MetricName"])
	}
	if object["MetricType"] != nil {
		d.Set("metric_type", object["MetricType"])
	}
	return nil
}

func resourceAliyunEssScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	request := ess.CreateDeleteScalingRuleRequest()
	request.ScalingRuleId = d.Id()
	request.RegionId = client.RegionId
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.DeleteScalingRule(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"InvalidScalingRuleId.NotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return WrapError(essService.WaitForEssScalingRule(d.Id(), Deleted, DefaultTimeout))
}

func resourceAliyunEssScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {

	//Compatible with older versions id
	if strings.Contains(d.Id(), COLON_SEPARATED) {
		parts, _ := ParseResourceId(d.Id(), 2)
		d.SetId(parts[1])
	}

	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"ScalingRuleId": d.Id(),
		"RegionId":      client.RegionId,
	}
	if d.HasChange("scaling_rule_name") {
		request["ScalingRuleName"] = d.Get("scaling_rule_name")
	}
	if d.HasChange("metric_type") {
		request["MetricType"] = d.Get("metric_type")
	}
	scalingRuleType := d.Get("scaling_rule_type")
	switch scalingRuleType {
	case string(SimpleScalingRule):
		if d.HasChange("adjustment_type") {
			request["AdjustmentType"] = d.Get("adjustment_type").(string)
		}
		if d.HasChange("min_adjustment_magnitude") {
			if v, ok := d.GetOkExists("min_adjustment_magnitude"); ok {
				request["MinAdjustmentMagnitude"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("adjustment_value") {
			if v, ok := d.GetOkExists("adjustment_value"); ok {
				request["AdjustmentValue"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("cooldown") {
			if v, ok := d.GetOkExists("cooldown"); ok {
				request["Cooldown"] = requests.NewInteger(v.(int))
			}
		}
	case string(TargetTrackingScalingRule):
		if d.HasChange("metric_name") {
			request["MetricName"] = d.Get("metric_name").(string)
		}
		if d.HasChange("scale_in_evaluation_count") {
			if v, ok := d.GetOkExists("scale_in_evaluation_count"); ok {
				request["ScaleInEvaluationCount"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("scale_out_evaluation_count") {
			if v, ok := d.GetOkExists("scale_out_evaluation_count"); ok {
				request["ScaleOutEvaluationCount"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("disable_scale_in") {
			if v, ok := d.GetOkExists("disable_scale_in"); ok {
				request["DisableScaleIn"] = requests.NewBoolean(v.(bool))
			}
		}
		if d.HasChange("estimated_instance_warmup") {
			if v, ok := d.GetOkExists("estimated_instance_warmup"); ok {
				request["EstimatedInstanceWarmup"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("target_value") {
			if v, ok := d.GetOkExists("target_value"); ok {
				targetValue, err := strconv.ParseFloat(strconv.FormatFloat(v.(float64), 'f', 3, 64), 64)
				if err != nil {
					return WrapError(err)
				}
				request["TargetValue"] = requests.NewFloat(targetValue)
			}
		}
		if d.HasChange("hybrid_monitor_namespace") {
			request["HybridMonitorNamespace"] = d.Get("hybrid_monitor_namespace")
		}
		if d.HasChange("hybrid_metrics") {
			hybridMetrics := make([]map[string]interface{}, 0)
			for _, e := range d.Get("hybrid_metrics").([]interface{}) {
				pack := e.(map[string]interface{})
				hybridMetric := make(map[string]interface{})
				hybridMetric["Id"] = pack["id"].(string)
				if pack["metric_name"] != nil && pack["metric_name"] != "" {
					hybridMetric["MetricName"] = pack["metric_name"].(string)
				}
				if pack["statistic"] != nil && pack["statistic"] != "" {
					hybridMetric["Statistic"] = pack["statistic"].(string)
				}
				if pack["expression"] != nil && pack["expression"] != "" {
					hybridMetric["Expression"] = pack["expression"].(string)
				}
				if pack["dimensions"] != nil {
					dimensions := make([]map[string]interface{}, len(pack["dimensions"].(*schema.Set).List()))
					for i, dimension := range pack["dimensions"].(*schema.Set).List() {
						dimensionVar := dimension.(map[string]interface{})
						dimensions[i] = make(map[string]interface{})
						dimensions[i]["DimensionKey"] = dimensionVar["dimension_key"]
						dimensions[i]["DimensionValue"] = dimensionVar["dimension_value"]
					}
					hybridMetric["Dimensions"] = dimensions
				}
				hybridMetrics = append(hybridMetrics, hybridMetric)
			}
			request["HybridMetrics"] = hybridMetrics
		}
	case string(StepScalingRule):
		if d.HasChange("min_adjustment_magnitude") {
			if v, ok := d.GetOkExists("min_adjustment_magnitude"); ok {
				request["MinAdjustmentMagnitude"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("step_adjustment") {
			steps := make([]map[string]interface{}, 0)
			for _, e := range d.Get("step_adjustment").([]interface{}) {
				pack := e.(map[string]interface{})
				step := make(map[string]interface{})
				step["ScalingAdjustment"] = strconv.Itoa(pack["scaling_adjustment"].(int))
				if pack["metric_interval_lower_bound"] != "" {
					lowerBound, err := strconv.ParseFloat(pack["metric_interval_lower_bound"].(string), 64)
					if err != nil {
						return WrapError(err)
					}
					step["MetricIntervalLowerBound"] = strconv.FormatFloat(lowerBound, 'f', 3, 64)
				}
				if pack["metric_interval_upper_bound"] != "" {
					upperBound, err := strconv.ParseFloat(pack["metric_interval_upper_bound"].(string), 64)
					if err != nil {
						return WrapError(err)
					}
					step["MetricIntervalUpperBound"] = strconv.FormatFloat(upperBound, 'f', 3, 64)
				}
				steps = append(steps, step)
			}
			request["StepAdjustment"] = &steps
		}

	case string(PredictiveScalingRule):
		if d.HasChange("metric_name") {
			request["MetricName"] = d.Get("metric_name").(string)
		}
		if d.HasChange("target_value") {
			if v, ok := d.GetOkExists("target_value"); ok {
				targetValue, err := strconv.ParseFloat(strconv.FormatFloat(v.(float64), 'f', 3, 64), 64)
				if err != nil {
					return WrapError(err)
				}
				request["TargetValue"] = requests.NewFloat(targetValue)
			}
		}
		if d.HasChange("initial_max_size") {
			if v, ok := d.GetOkExists("initial_max_size"); ok {
				request["InitialMaxSize"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("predictive_value_buffer") {
			if v, ok := d.GetOkExists("predictive_value_buffer"); ok {
				request["PredictiveValueBuffer"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("predictive_task_buffer_time") {
			if v, ok := d.GetOkExists("predictive_task_buffer_time"); ok {
				request["PredictiveTaskBufferTime"] = requests.NewInteger(v.(int))
			}
		}
		if d.HasChange("predictive_scaling_mode") {
			request["PredictiveScalingMode"] = d.Get("predictive_scaling_mode").(string)
		}
		if d.HasChange("predictive_value_behavior") {
			request["PredictiveValueBehavior"] = d.Get("predictive_value_behavior").(string)
		}
	}

	raw, err := client.RpcPost("Ess", "2014-08-28", "ModifyScalingRule", nil, request, true)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ModifyScalingRule", AlibabaCloudSdkGoERROR)
	}
	addDebug("ModifyScalingRule", raw, request, request)
	return resourceAliyunEssScalingRuleRead(d, meta)
}

func buildAlicloudEssScalingRuleArgs(d *schema.ResourceData, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*connectivity.AliyunClient)
	request := map[string]interface{}{
		"ScalingGroupId":  d.Get("scaling_group_id"),
		"ScalingRuleType": d.Get("scaling_rule_type"),
		"RegionId":        client.RegionId,
	}
	// common params
	scalingRuleType := d.Get("scaling_rule_type").(string)
	if v, ok := d.GetOk("scaling_rule_name"); ok && v.(string) != "" {
		request["ScalingRuleName"] = v.(string)
	}
	if v, ok := d.GetOk("metric_type"); ok && v.(string) != "" {
		request["MetricType"] = v.(string)
	}
	switch scalingRuleType {
	case string(SimpleScalingRule):
		if v, ok := d.GetOk("adjustment_type"); ok && v.(string) != "" {
			request["AdjustmentType"] = v.(string)
		}
		if v, ok := d.GetOkExists("min_adjustment_magnitude"); ok {
			request["MinAdjustmentMagnitude"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("adjustment_value"); ok {
			request["AdjustmentValue"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("cooldown"); ok {
			request["Cooldown"] = requests.NewInteger(v.(int))
		}
	case string(TargetTrackingScalingRule):
		if v, ok := d.GetOkExists("estimated_instance_warmup"); ok {
			request["EstimatedInstanceWarmup"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("scale_in_evaluation_count"); ok {
			request["ScaleInEvaluationCount"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("scale_out_evaluation_count"); ok {
			request["ScaleOutEvaluationCount"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOk("metric_name"); ok && v.(string) != "" {
			request["MetricName"] = v.(string)
		}
		if v, ok := d.GetOkExists("target_value"); ok {
			targetValue, err := strconv.ParseFloat(strconv.FormatFloat(v.(float64), 'f', 3, 64), 64)
			if err != nil {
				return nil, WrapError(err)
			}
			request["TargetValue"] = requests.NewFloat(targetValue)
		}
		if v, ok := d.GetOkExists("disable_scale_in"); ok {
			request["DisableScaleIn"] = requests.NewBoolean(v.(bool))
		}
		v, ok := d.GetOk("alarm_dimension")
		if ok {
			steps := make([]ess.CreateScalingRuleAlarmDimension, 0)
			for _, e := range v.([]interface{}) {
				pack := e.(map[string]interface{})
				step := ess.CreateScalingRuleAlarmDimension{
					DimensionKey:   pack["dimension_key"].(string),
					DimensionValue: pack["dimension_value"].(string),
				}
				steps = append(steps, step)
			}
			request["AlarmDimension"] = &steps
		}
		if hybridMonitorNamespace, exist := d.GetOk("hybrid_monitor_namespace"); exist && hybridMonitorNamespace.(string) != "" {
			request["HybridMonitorNamespace"] = hybridMonitorNamespace.(string)
		}
		hM, ok := d.GetOk("hybrid_metrics")
		if ok {
			hybridMetrics := make([]map[string]interface{}, 0)
			for _, e := range hM.([]interface{}) {
				pack := e.(map[string]interface{})
				hybridMetric := make(map[string]interface{})
				hybridMetric["Id"] = pack["id"].(string)
				if pack["metric_name"] != nil && pack["metric_name"] != "" {
					hybridMetric["MetricName"] = pack["metric_name"].(string)
				}
				if pack["statistic"] != nil && pack["statistic"] != "" {
					hybridMetric["Statistic"] = pack["statistic"].(string)
				}
				if pack["expression"] != nil && pack["expression"] != "" {
					hybridMetric["Expression"] = pack["expression"].(string)
				}
				if pack["dimensions"] != nil {
					dimensions := make([]map[string]interface{}, len(pack["dimensions"].(*schema.Set).List()))
					for i, dimension := range pack["dimensions"].(*schema.Set).List() {
						dimensionVar := dimension.(map[string]interface{})
						dimensions[i] = make(map[string]interface{})
						dimensions[i]["DimensionKey"] = dimensionVar["dimension_key"]
						dimensions[i]["DimensionValue"] = dimensionVar["dimension_value"]
					}
					hybridMetric["Dimensions"] = dimensions
				}
				hybridMetrics = append(hybridMetrics, hybridMetric)
			}
			request["HybridMetrics"] = hybridMetrics
		}
	case string(StepScalingRule):
		v, ok := d.GetOk("step_adjustment")
		if v, ok := d.GetOk("adjustment_type"); ok && v.(string) != "" {
			request["AdjustmentType"] = v.(string)
		}
		if v, ok := d.GetOkExists("min_adjustment_magnitude"); ok {
			request["MinAdjustmentMagnitude"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("estimated_instance_warmup"); ok {
			request["EstimatedInstanceWarmup"] = requests.NewInteger(v.(int))
		}
		if ok {
			steps := make([]map[string]interface{}, 0)
			for _, e := range v.([]interface{}) {
				pack := e.(map[string]interface{})
				step := make(map[string]interface{})
				step["ScalingAdjustment"] = strconv.Itoa(pack["scaling_adjustment"].(int))
				if pack["metric_interval_lower_bound"] != "" {
					lowerBound, err := strconv.ParseFloat(pack["metric_interval_lower_bound"].(string), 64)
					if err != nil {
						return nil, WrapError(err)
					}
					step["MetricIntervalLowerBound"] = strconv.FormatFloat(lowerBound, 'f', 3, 64)
				}
				if pack["metric_interval_upper_bound"] != "" {
					upperBound, err := strconv.ParseFloat(pack["metric_interval_upper_bound"].(string), 64)
					if err != nil {
						return nil, WrapError(err)
					}
					step["MetricIntervalUpperBound"] = strconv.FormatFloat(upperBound, 'f', 3, 64)
				}
				steps = append(steps, step)
			}
			request["StepAdjustment"] = &steps
		}
	case string(PredictiveScalingRule):
		if v, ok := d.GetOk("metric_name"); ok && v.(string) != "" {
			request["MetricName"] = v.(string)
		}
		if v, ok := d.GetOkExists("target_value"); ok {
			targetValue, err := strconv.ParseFloat(strconv.FormatFloat(v.(float64), 'f', 3, 64), 64)
			if err != nil {
				return nil, WrapError(err)
			}
			request["TargetValue"] = requests.NewFloat(targetValue)
		}
		if v, ok := d.GetOk("predictive_scaling_mode"); ok && v.(string) != "" {
			request["PredictiveScalingMode"] = v.(string)
		}
		if v, ok := d.GetOk("predictive_value_behavior"); ok && v.(string) != "" {
			request["PredictiveValueBehavior"] = v.(string)
		}
		if v, ok := d.GetOkExists("predictive_value_buffer"); ok {
			request["PredictiveValueBuffer"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("predictive_task_buffer_time"); ok {
			request["PredictiveTaskBufferTime"] = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("initial_max_size"); ok {
			request["InitialMaxSize"] = requests.NewInteger(v.(int))
		}
	}
	return request, nil
}

func flattenStepAdjustmentMappings(list interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, i := range list.(map[string]interface{})["StepAdjustment"].([]interface{}) {
		stepAdjustment := i.(map[string]interface{})
		l := map[string]interface{}{}
		if stepAdjustment["MetricIntervalLowerBound"] != nil {
			lowerBound, err := strconv.ParseFloat(strconv.FormatFloat(stepAdjustment["MetricIntervalLowerBound"].(float64), 'f', 3, 64), 64)
			if err != nil {
				return nil, WrapError(err)
			}
			l["metric_interval_lower_bound"] = lowerBound
		}

		if stepAdjustment["MetricIntervalUpperBound"] != nil {
			upperBound, err := strconv.ParseFloat(strconv.FormatFloat(stepAdjustment["MetricIntervalUpperBound"].(float64), 'f', 3, 64), 64)
			if err != nil {
				return nil, WrapError(err)
			}
			l["metric_interval_upper_bound"] = upperBound
		}
		l["scaling_adjustment"] = stepAdjustment["ScalingAdjustment"]
		result = append(result, l)
	}
	return result, nil
}

func flattenHybridMetricsMappings(list interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)
	for _, i := range list.(map[string]interface{})["HybridMetric"].([]interface{}) {
		hybridMetric := i.(map[string]interface{})
		l := map[string]interface{}{}
		l["id"] = hybridMetric["Id"].(string)
		if hybridMetric["MetricName"] != nil {
			l["metric_name"] = hybridMetric["MetricName"].(string)
		}
		if hybridMetric["Statistic"] != nil {
			l["statistic"] = hybridMetric["Statistic"].(string)
		}
		if hybridMetric["Expression"] != nil {
			l["expression"] = hybridMetric["Expression"].(string)
		}

		if hybridMetric["Dimensions"] != nil {
			dimensionsVarsMaps := make([]map[string]interface{}, 0)
			m := hybridMetric["Dimensions"].(map[string]interface{})
			for _, dimensionsVarsValue := range m["Dimension"].([]interface{}) {
				dimensionsVars := dimensionsVarsValue.(map[string]interface{})
				dimensionsVarsMap := map[string]interface{}{
					"dimension_key":   dimensionsVars["DimensionKey"],
					"dimension_value": dimensionsVars["DimensionValue"],
				}
				dimensionsVarsMaps = append(dimensionsVarsMaps, dimensionsVarsMap)
			}
			l["dimensions"] = dimensionsVarsMaps
		}
		result = append(result, l)
	}
	return result, nil
}

func flattenAlarmDimensionMappings(list interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	for _, i := range list.(map[string]interface{})["AlarmDimension"].([]interface{}) {
		alarmDimension := i.(map[string]interface{})
		l := map[string]interface{}{
			"dimension_key":   alarmDimension["DimensionKey"],
			"dimension_value": alarmDimension["DimensionValue"],
		}
		result = append(result, l)
	}
	return result
}
