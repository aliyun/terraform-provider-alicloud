package alicloud

import (
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
	request.RegionId = client.RegionId
	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.CreateScalingRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ess_scalingrule", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*ess.CreateScalingRuleResponse)

	d.SetId(response.ScalingRuleId)
	return resourceAliyunEssScalingRuleRead(d, meta)
}

func resourceAliyunEssScalingRuleRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	essService := EssService{client}
	_, ok := d.GetOk("alarm_dimension")
	var object ess.ScalingRule
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

	d.Set("scaling_group_id", object.ScalingGroupId)
	d.Set("ari", object.ScalingRuleAri)
	d.Set("adjustment_type", object.AdjustmentType)
	d.Set("adjustment_value", object.AdjustmentValue)
	d.Set("scaling_rule_name", object.ScalingRuleName)
	d.Set("cooldown", object.Cooldown)
	d.Set("scaling_rule_type", object.ScalingRuleType)
	d.Set("estimated_instance_warmup", object.EstimatedInstanceWarmup)
	d.Set("scale_in_evaluation_count", object.ScaleInEvaluationCount)
	d.Set("scale_out_evaluation_count", object.ScaleOutEvaluationCount)
	d.Set("min_adjustment_magnitude", object.MinAdjustmentMagnitude)
	d.Set("metric_name", object.MetricName)
	d.Set("predictive_scaling_mode", object.PredictiveScalingMode)
	d.Set("initial_max_size", object.InitialMaxSize)
	d.Set("predictive_value_behavior", object.PredictiveValueBehavior)
	d.Set("predictive_value_buffer", object.PredictiveValueBuffer)
	d.Set("predictive_task_buffer_time", object.PredictiveTaskBufferTime)
	targetValue, err := strconv.ParseFloat(strconv.FormatFloat(object.TargetValue, 'f', 3, 64), 64)
	if err != nil {
		return WrapError(err)
	}
	d.Set("target_value", targetValue)
	d.Set("disable_scale_in", object.DisableScaleIn)
	steps, err := flattenStepAdjustmentMappings(object.StepAdjustments.StepAdjustment)
	alarmDimensions := flattenAlarmDimensionMappings(object.AlarmDimensions.AlarmDimension)
	d.Set("alarm_dimension", alarmDimensions)
	if err != nil {
		return WrapError(err)
	}
	d.Set("step_adjustment", steps)

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
	request := ess.CreateModifyScalingRuleRequest()
	request.ScalingRuleId = d.Id()
	request.RegionId = client.RegionId
	if d.HasChange("scaling_rule_name") {
		request.ScalingRuleName = d.Get("scaling_rule_name").(string)
	}
	scalingRuleType := d.Get("scaling_rule_type")
	switch scalingRuleType {
	case string(SimpleScalingRule):
		if d.HasChange("adjustment_type") {
			request.AdjustmentType = d.Get("adjustment_type").(string)
		}
		if d.HasChange("min_adjustment_magnitude") {
			request.MinAdjustmentMagnitude = requests.NewInteger(d.Get("min_adjustment_magnitude").(int))
		}
		if d.HasChange("adjustment_value") {
			request.AdjustmentValue = requests.NewInteger(d.Get("adjustment_value").(int))
		}
		if d.HasChange("cooldown") {
			request.Cooldown = requests.NewInteger(d.Get("cooldown").(int))
		}
	case string(TargetTrackingScalingRule):
		if d.HasChange("metric_name") {
			request.MetricName = d.Get("metric_name").(string)
		}
		if d.HasChange("scale_in_evaluation_count") {
			request.ScaleInEvaluationCount = requests.NewInteger(d.Get("scale_in_evaluation_count").(int))
		}
		if d.HasChange("scale_out_evaluation_count") {
			request.ScaleOutEvaluationCount = requests.NewInteger(d.Get("scale_out_evaluation_count").(int))
		}
		if d.HasChange("disable_scale_in") {
			request.DisableScaleIn = requests.NewBoolean(d.Get("disable_scale_in").(bool))
		}
		if d.HasChange("estimated_instance_warmup") {
			request.EstimatedInstanceWarmup = requests.NewInteger(d.Get("estimated_instance_warmup").(int))
		}
		if d.HasChange("target_value") {
			targetValue, err := strconv.ParseFloat(strconv.FormatFloat(d.Get("target_value").(float64), 'f', 3, 64), 64)
			if err != nil {
				return WrapError(err)
			}
			request.TargetValue = requests.NewFloat(targetValue)
		}
	case string(StepScalingRule):
		if d.HasChange("min_adjustment_magnitude") {
			request.MinAdjustmentMagnitude = requests.NewInteger(d.Get("min_adjustment_magnitude").(int))
		}
		if d.HasChange("step_adjustment") {
			steps := make([]ess.ModifyScalingRuleStepAdjustment, 0)
			for _, e := range d.Get("step_adjustment").([]interface{}) {
				pack := e.(map[string]interface{})
				step := ess.ModifyScalingRuleStepAdjustment{
					ScalingAdjustment: strconv.Itoa(pack["scaling_adjustment"].(int)),
				}
				if pack["metric_interval_lower_bound"] != "" {
					lowerBound, err := strconv.ParseFloat(pack["metric_interval_lower_bound"].(string), 64)
					if err != nil {
						return WrapError(err)
					}
					step.MetricIntervalLowerBound = strconv.FormatFloat(lowerBound, 'f', 3, 64)
				}
				if pack["metric_interval_upper_bound"] != "" {
					upperBound, err := strconv.ParseFloat(pack["metric_interval_upper_bound"].(string), 64)
					if err != nil {
						return WrapError(err)
					}
					step.MetricIntervalUpperBound = strconv.FormatFloat(upperBound, 'f', 3, 64)
				}
				steps = append(steps, step)
			}
			request.StepAdjustment = &steps
		}
	case string(PredictiveScalingRule):
		if d.HasChange("metric_name") {
			request.MetricName = d.Get("metric_name").(string)
		}
		if d.HasChange("target_value") {
			targetValue, err := strconv.ParseFloat(strconv.FormatFloat(d.Get("target_value").(float64), 'f', 3, 64), 64)
			if err != nil {
				return WrapError(err)
			}
			request.TargetValue = requests.NewFloat(targetValue)
		}
		if d.HasChange("initial_max_size") {
			request.InitialMaxSize = requests.NewInteger(d.Get("initial_max_size").(int))
		}
		if d.HasChange("predictive_value_buffer") {
			request.PredictiveValueBuffer = requests.NewInteger(d.Get("predictive_value_buffer").(int))
		}
		if d.HasChange("predictive_task_buffer_time") {
			request.PredictiveTaskBufferTime = requests.NewInteger(d.Get("predictive_task_buffer_time").(int))
		}
		if d.HasChange("predictive_scaling_mode") {
			request.PredictiveScalingMode = d.Get("predictive_scaling_mode").(string)
		}
		if d.HasChange("predictive_value_behavior") {
			request.PredictiveValueBehavior = d.Get("predictive_value_behavior").(string)
		}
	}

	raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
		return essClient.ModifyScalingRule(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return resourceAliyunEssScalingRuleRead(d, meta)
}

func buildAlicloudEssScalingRuleArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingRuleRequest, error) {
	request := ess.CreateCreateScalingRuleRequest()
	// common params
	request.ScalingGroupId = d.Get("scaling_group_id").(string)
	scalingRuleType := d.Get("scaling_rule_type").(string)
	if v, ok := d.GetOk("scaling_rule_name"); ok && v.(string) != "" {
		request.ScalingRuleName = v.(string)
	}
	request.ScalingRuleType = d.Get("scaling_rule_type").(string)

	switch scalingRuleType {
	case string(SimpleScalingRule):
		if v, ok := d.GetOk("adjustment_type"); ok && v.(string) != "" {
			request.AdjustmentType = v.(string)
		}
		if v, ok := d.GetOkExists("min_adjustment_magnitude"); ok {
			request.MinAdjustmentMagnitude = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("adjustment_value"); ok {
			request.AdjustmentValue = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOk("cooldown"); ok {
			request.Cooldown = requests.NewInteger(v.(int))
		}
	case string(TargetTrackingScalingRule):
		if v, ok := d.GetOk("estimated_instance_warmup"); ok {
			request.EstimatedInstanceWarmup = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOk("scale_in_evaluation_count"); ok {
			request.ScaleInEvaluationCount = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOk("scale_out_evaluation_count"); ok {
			request.ScaleOutEvaluationCount = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOk("metric_name"); ok && v.(string) != "" {
			request.MetricName = v.(string)
		}
		if v, ok := d.GetOk("target_value"); ok {
			targetValue, err := strconv.ParseFloat(strconv.FormatFloat(v.(float64), 'f', 3, 64), 64)
			if err != nil {
				return nil, WrapError(err)
			}
			request.TargetValue = requests.NewFloat(targetValue)
		}
		if v, ok := d.GetOk("disable_scale_in"); ok {
			request.DisableScaleIn = requests.NewBoolean(v.(bool))
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
			request.AlarmDimension = &steps
		}
	case string(StepScalingRule):
		v, ok := d.GetOk("step_adjustment")
		if v, ok := d.GetOk("adjustment_type"); ok && v.(string) != "" {
			request.AdjustmentType = v.(string)
		}
		if v, ok := d.GetOkExists("min_adjustment_magnitude"); ok {
			request.MinAdjustmentMagnitude = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOk("estimated_instance_warmup"); ok {
			request.EstimatedInstanceWarmup = requests.NewInteger(v.(int))
		}
		if ok {
			steps := make([]ess.CreateScalingRuleStepAdjustment, 0)
			for _, e := range v.([]interface{}) {
				pack := e.(map[string]interface{})
				step := ess.CreateScalingRuleStepAdjustment{
					ScalingAdjustment: strconv.Itoa(pack["scaling_adjustment"].(int)),
				}
				if pack["metric_interval_lower_bound"] != "" {
					lowerBound, err := strconv.ParseFloat(pack["metric_interval_lower_bound"].(string), 64)
					if err != nil {
						return nil, WrapError(err)
					}
					step.MetricIntervalLowerBound = strconv.FormatFloat(lowerBound, 'f', 3, 64)
				}
				if pack["metric_interval_upper_bound"] != "" {
					upperBound, err := strconv.ParseFloat(pack["metric_interval_upper_bound"].(string), 64)
					if err != nil {
						return nil, WrapError(err)
					}
					step.MetricIntervalUpperBound = strconv.FormatFloat(upperBound, 'f', 3, 64)
				}
				steps = append(steps, step)
			}
			request.StepAdjustment = &steps
		}
	case string(PredictiveScalingRule):
		if v, ok := d.GetOk("metric_name"); ok && v.(string) != "" {
			request.MetricName = v.(string)
		}
		if v, ok := d.GetOk("target_value"); ok {
			targetValue, err := strconv.ParseFloat(strconv.FormatFloat(v.(float64), 'f', 3, 64), 64)
			if err != nil {
				return nil, WrapError(err)
			}
			request.TargetValue = requests.NewFloat(targetValue)
		}
		if v, ok := d.GetOk("predictive_scaling_mode"); ok && v.(string) != "" {
			request.PredictiveScalingMode = v.(string)
		}
		if v, ok := d.GetOk("predictive_value_behavior"); ok && v.(string) != "" {
			request.PredictiveValueBehavior = v.(string)
		}
		if v, ok := d.GetOkExists("predictive_value_buffer"); ok {
			request.PredictiveValueBuffer = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("predictive_task_buffer_time"); ok {
			request.PredictiveTaskBufferTime = requests.NewInteger(v.(int))
		}
		if v, ok := d.GetOkExists("initial_max_size"); ok {
			request.InitialMaxSize = requests.NewInteger(v.(int))
		}
	}
	return request, nil
}

func flattenStepAdjustmentMappings(list []ess.StepAdjustment) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		lowerBound, err := strconv.ParseFloat(strconv.FormatFloat(i.MetricIntervalLowerBound, 'f', 3, 64), 64)
		if err != nil {
			return nil, WrapError(err)
		}
		upperBound, err := strconv.ParseFloat(strconv.FormatFloat(i.MetricIntervalUpperBound, 'f', 3, 64), 64)
		if err != nil {
			return nil, WrapError(err)
		}
		l := map[string]interface{}{
			"metric_interval_lower_bound": lowerBound,
			"metric_interval_upper_bound": upperBound,
			"scaling_adjustment":          i.ScalingAdjustment,
		}
		result = append(result, l)
	}
	return result, nil
}

func flattenAlarmDimensionMappings(list []ess.AlarmDimension) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, i := range list {
		l := map[string]interface{}{
			"dimension_key":   i.DimensionKey,
			"dimension_value": i.DimensionValue,
		}
		result = append(result, l)
	}
	return result
}
