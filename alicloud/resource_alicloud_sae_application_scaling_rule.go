package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudSaeApplicationScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudSaeApplicationScalingRuleCreate,
		Read:   resourceAlicloudSaeApplicationScalingRuleRead,
		Update: resourceAlicloudSaeApplicationScalingRuleUpdate,
		Delete: resourceAlicloudSaeApplicationScalingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_rule_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9-]{1,31}`), "The name must be `2` to `32` characters in length and can contain lowercase letters, digits, and hyphens (-). The name must start with a lowercase letter."),
			},
			"scaling_rule_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"mix", "timing", "metric"}, false),
			},
			"scaling_rule_metric": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"scaling_rule_metric", "scaling_rule_timer"},
				MaxItems:     1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_replicas": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"min_replicas": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"metrics": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_target_average_utilization": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"metric_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"scale_up_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"step": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"stabilization_window_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"scale_down_rules": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"step": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"stabilization_window_seconds": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"disabled": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"min_ready_instances": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_ready_instance_ratio": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scaling_rule_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"scaling_rule_timer": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"begin_date": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"end_date": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"schedules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"at_time": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"target_replicas": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 50),
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return fmt.Sprint(d.Get("scaling_rule_type")) != "timing"
										},
									},
									"max_replicas": {
										Type:     schema.TypeInt,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return fmt.Sprint(d.Get("scaling_rule_type")) != "mix"
										},
									},
									"min_replicas": {
										Type:     schema.TypeInt,
										Optional: true,
										DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
											return fmt.Sprint(d.Get("scaling_rule_type")) != "mix"
										},
									},
								},
							},
						},
						"period": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudSaeApplicationScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "/pop/v1/sam/scale/applicationScalingRule"
	request := make(map[string]*string)
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	request["AppId"] = StringPointer(d.Get("app_id").(string))
	request["ScalingRuleName"] = StringPointer(d.Get("scaling_rule_name").(string))
	request["ScalingRuleType"] = StringPointer(d.Get("scaling_rule_type").(string))

	scalingRuleTimerMap := map[string]interface{}{}
	if scalingRuleTimer := d.Get("scaling_rule_timer").([]interface{}); len(scalingRuleTimer) > 0 {
		scalingRuleTimerArg := scalingRuleTimer[0].(map[string]interface{})
		scalingRuleTimerMap["beginDate"] = scalingRuleTimerArg["begin_date"]
		scalingRuleTimerMap["endDate"] = scalingRuleTimerArg["end_date"]
		scalingRuleTimerMap["period"] = scalingRuleTimerArg["period"]

		schedulesMaps := make([]map[string]interface{}, 0)
		for _, schedules := range scalingRuleTimerArg["schedules"].([]interface{}) {
			schedulesMap := map[string]interface{}{}
			schedulesArg := schedules.(map[string]interface{})
			if v, ok := schedulesArg["max_replicas"]; ok && v.(int) > 0 {
				schedulesMap["maxReplicas"] = v
			}
			if v, ok := schedulesArg["min_replicas"]; ok && v.(int) > 0 {
				schedulesMap["minReplicas"] = v
			}
			schedulesMap["atTime"] = schedulesArg["at_time"]
			if v, ok := schedulesArg["target_replicas"]; ok && v.(int) > 0 {
				schedulesMap["targetReplicas"] = v
			}
			schedulesMaps = append(schedulesMaps, schedulesMap)
		}

		scalingRuleTimerMap["schedules"] = schedulesMaps
		if v, err := convertArrayObjectToJsonString(scalingRuleTimerMap); err != nil {
			return WrapError(err)
		} else {
			request["ScalingRuleTimer"] = StringPointer(v)
		}
	}

	scalingRuleMetricMap := map[string]interface{}{}
	if scalingRuleMetric := d.Get("scaling_rule_metric").(*schema.Set).List(); len(scalingRuleMetric) > 0 {
		scalingRuleMetricArg := scalingRuleMetric[0].(map[string]interface{})
		scalingRuleMetricMap["maxReplicas"] = scalingRuleMetricArg["max_replicas"]
		scalingRuleMetricMap["minReplicas"] = scalingRuleMetricArg["min_replicas"]

		scaleUpRulesMap := map[string]interface{}{}
		for _, scaleUpRules := range scalingRuleMetricArg["scale_up_rules"].(*schema.Set).List() {
			scaleUpRulesArg := scaleUpRules.(map[string]interface{})
			scaleUpRulesMap["step"] = fmt.Sprint(scaleUpRulesArg["step"])
			scaleUpRulesMap["disabled"] = scaleUpRulesArg["disabled"]
			scaleUpRulesMap["stabilizationWindowSeconds"] = scaleUpRulesArg["stabilization_window_seconds"]
		}
		scalingRuleMetricMap["scaleUpRules"] = scaleUpRulesMap

		scaleDownRulesMap := map[string]interface{}{}
		for _, scaleDownRules := range scalingRuleMetricArg["scale_down_rules"].(*schema.Set).List() {
			scaleDownRulesArg := scaleDownRules.(map[string]interface{})
			scaleDownRulesMap["step"] = fmt.Sprint(scaleDownRulesArg["step"])
			scaleDownRulesMap["disabled"] = scaleDownRulesArg["disabled"]
			scaleDownRulesMap["stabilizationWindowSeconds"] = scaleDownRulesArg["stabilization_window_seconds"]
		}
		scalingRuleMetricMap["scaleDownRules"] = scaleDownRulesMap

		metricsMaps := make([]map[string]interface{}, 0)
		if v, ok := scalingRuleMetricArg["metrics"]; ok {
			for _, metrics := range v.(*schema.Set).List() {
				metricsArg := metrics.(map[string]interface{})
				metricsMap := map[string]interface{}{}
				metricsMap["metricType"] = metricsArg["metric_type"]
				metricsMap["metricTargetAverageUtilization"] = metricsArg["metric_target_average_utilization"]
				metricsMaps = append(metricsMaps, metricsMap)
			}
		}
		scalingRuleMetricMap["metrics"] = metricsMaps

		if v, err := convertArrayObjectToJsonString(scalingRuleMetricMap); err != nil {
			return WrapError(err)
		} else {
			request["ScalingRuleMetric"] = StringPointer(v)
		}
	}

	if v, ok := d.GetOk("min_ready_instances"); ok {
		request["MinReadyInstances"] = StringPointer(strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("min_ready_instance_ratio"); ok {
		request["MinReadyInstanceRatio"] = StringPointer(strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOkExists("scaling_rule_enable"); ok {
		request["ScalingRuleEnable"] = StringPointer(strconv.FormatBool(v.(bool)))
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("POST"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_sae_application_scaling_rule", "POST "+action, AlibabaCloudSdkGoERROR)
	}
	if respBody, isExist := response["body"]; isExist {
		response = respBody.(map[string]interface{})
	} else {
		return WrapError(fmt.Errorf("%s failed, response: %v", "POST "+action, response))
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(*request["AppId"], ":", responseData["ScaleRuleName"]))
	return resourceAlicloudSaeApplicationScalingRuleRead(d, meta)
}
func resourceAlicloudSaeApplicationScalingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	saeService := SaeService{client}
	object, err := saeService.DescribeSaeApplicationScalingRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_sae_application_scaling_rule saeService.DescribeSaeApplicationScalingRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("scaling_rule_enable", object["ScaleRuleEnabled"])
	d.Set("app_id", object["AppId"])
	d.Set("scaling_rule_type", object["ScaleRuleType"])
	d.Set("scaling_rule_name", object["ScaleRuleName"])

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
		if schedules, ok := scalingRuleMetricArg["Metrics"]; ok {
			schedulesArg := schedules.([]interface{})
			for _, v := range schedulesArg {
				metricsArg := v.(map[string]interface{})
				metricsMap := map[string]interface{}{}
				metricsMap["metric_type"] = metricsArg["MetricType"]
				metricsMap["metric_target_average_utilization"] = metricsArg["MetricTargetAverageUtilization"]
				metricsMaps = append(metricsMaps, metricsMap)
			}
			scalingRuleMetricMap["metrics"] = metricsMaps
		}

		scaleUpRulesMaps := make([]map[string]interface{}, 0)
		if scaleUpRules, ok := scalingRuleMetricArg["ScaleUpRules"]; ok {
			if scaleUpRulesArg, ok := scaleUpRules.(map[string]interface{}); ok && len(scaleUpRulesArg) > 0 {
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
			if scaleDownRulesArg, ok := scaleDownRules.(map[string]interface{}); ok && len(scaleDownRulesArg) > 0 {
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
			d.Set("scaling_rule_metric", scalingRuleMetricMaps)
		}
	}

	scalingRuleTimerMaps := make([]map[string]interface{}, 0)
	if scalingRuleTimer, ok := object["Timer"]; ok {
		if scalingRuleTimerArg, ok := scalingRuleTimer.(map[string]interface{}); ok && len(scalingRuleTimerArg) > 0 {
			scalingRuleTimerMap := make(map[string]interface{}, 0)
			scalingRuleTimerMap["end_date"] = scalingRuleTimerArg["EndDate"]
			scalingRuleTimerMap["begin_date"] = scalingRuleTimerArg["BeginDate"]
			scalingRuleTimerMap["period"] = scalingRuleTimerArg["Period"]
			schedulesMaps := make([]map[string]interface{}, 0)
			if schedules, ok := scalingRuleTimerArg["Schedules"]; ok {
				schedulesArg := schedules.([]interface{})
				for _, v := range schedulesArg {
					serverGroupTuplesArg := v.(map[string]interface{})
					schedulesMap := map[string]interface{}{}
					schedulesMap["at_time"] = serverGroupTuplesArg["AtTime"]
					if v, ok := serverGroupTuplesArg["TargetReplicas"]; ok {
						schedulesMap["target_replicas"] = v
					}
					if v, ok := serverGroupTuplesArg["MinReplicas"]; ok {
						schedulesMap["min_replicas"] = v
					}
					if v, ok := serverGroupTuplesArg["MaxReplicas"]; ok {
						schedulesMap["max_replicas"] = v
					}
					schedulesMaps = append(schedulesMaps, schedulesMap)
				}
				scalingRuleTimerMap["schedules"] = schedulesMaps
			}
			scalingRuleTimerMaps = append(scalingRuleTimerMaps, scalingRuleTimerMap)
			d.Set("scaling_rule_timer", scalingRuleTimerMaps)
		}
	}

	return nil
}
func resourceAlicloudSaeApplicationScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	saeService := SaeService{client}
	var response map[string]interface{}
	update := false
	d.Partial(true)
	request := map[string]*string{
		"AppId":           StringPointer(parts[0]),
		"ScalingRuleName": StringPointer(parts[1]),
	}

	if d.HasChange("scaling_rule_timer") {
		update = true
	}
	scalingRuleTimerMap := map[string]interface{}{}

	if scalingRuleTimer := d.Get("scaling_rule_timer").([]interface{}); len(scalingRuleTimer) > 0 {
		scalingRuleTimerArg := scalingRuleTimer[0].(map[string]interface{})
		scalingRuleTimerMap["beginDate"] = scalingRuleTimerArg["begin_date"]
		scalingRuleTimerMap["endDate"] = scalingRuleTimerArg["end_date"]
		scalingRuleTimerMap["period"] = scalingRuleTimerArg["period"]

		schedulesMaps := make([]map[string]interface{}, 0)
		for _, schedules := range scalingRuleTimerArg["schedules"].([]interface{}) {
			schedulesMap := map[string]interface{}{}
			schedulesArg := schedules.(map[string]interface{})
			if v, ok := schedulesArg["max_replicas"]; ok && v.(int) > 0 {
				schedulesMap["maxReplicas"] = v
			}
			if v, ok := schedulesArg["min_replicas"]; ok && v.(int) > 0 {
				schedulesMap["minReplicas"] = v
			}
			schedulesMap["atTime"] = schedulesArg["at_time"]
			if v, ok := schedulesArg["target_replicas"]; ok && v.(int) > 0 {
				schedulesMap["targetReplicas"] = v
			}
			schedulesMaps = append(schedulesMaps, schedulesMap)
		}
		scalingRuleTimerMap["schedules"] = schedulesMaps
	}
	if v, err := convertArrayObjectToJsonString(scalingRuleTimerMap); err != nil {
		return WrapError(err)
	} else {
		request["ScalingRuleTimer"] = StringPointer(v)
	}

	if d.HasChange("scaling_rule_metric") {
		update = true
	}
	scalingRuleMetricMap := map[string]interface{}{}

	if scalingRuleMetric := d.Get("scaling_rule_metric").(*schema.Set).List(); len(scalingRuleMetric) > 0 {
		scalingRuleMetricArg := scalingRuleMetric[0].(map[string]interface{})
		scalingRuleMetricMap["maxReplicas"] = scalingRuleMetricArg["max_replicas"]
		scalingRuleMetricMap["minReplicas"] = scalingRuleMetricArg["min_replicas"]
		scaleUpRulesMap := map[string]interface{}{}
		for _, scaleUpRules := range scalingRuleMetricArg["scale_up_rules"].(*schema.Set).List() {
			scaleUpRulesArg := scaleUpRules.(map[string]interface{})
			scaleUpRulesMap["step"] = fmt.Sprint(scaleUpRulesArg["step"])
			scaleUpRulesMap["disabled"] = scaleUpRulesArg["disabled"]
			scaleUpRulesMap["stabilizationWindowSeconds"] = scaleUpRulesArg["stabilization_window_seconds"]
		}
		scalingRuleMetricMap["scaleUpRules"] = scaleUpRulesMap
		scaleDownRulesMap := map[string]interface{}{}
		for _, scaleDownRules := range scalingRuleMetricArg["scale_down_rules"].(*schema.Set).List() {
			scaleDownRulesArg := scaleDownRules.(map[string]interface{})
			scaleDownRulesMap["step"] = fmt.Sprint(scaleDownRulesArg["step"])
			scaleDownRulesMap["disabled"] = scaleDownRulesArg["disabled"]
			scaleDownRulesMap["stabilizationWindowSeconds"] = scaleDownRulesArg["stabilization_window_seconds"]
		}
		scalingRuleMetricMap["scaleDownRules"] = scaleDownRulesMap
		metricsMaps := make([]map[string]interface{}, 0)
		if v, ok := scalingRuleMetricArg["metrics"]; ok {
			for _, metrics := range v.(*schema.Set).List() {
				metricsArg := metrics.(map[string]interface{})
				metricsMap := map[string]interface{}{}
				metricsMap["metricType"] = metricsArg["metric_type"]
				metricsMap["metricTargetAverageUtilization"] = metricsArg["metric_target_average_utilization"]
				metricsMaps = append(metricsMaps, metricsMap)
			}
		}
		scalingRuleMetricMap["metrics"] = metricsMaps
	}

	if v, err := convertArrayObjectToJsonString(scalingRuleMetricMap); err != nil {
		return WrapError(err)
	} else {
		request["ScalingRuleMetric"] = StringPointer(v)
	}

	if d.HasChange("min_ready_instances") {
		update = true
		if v, ok := d.GetOk("min_ready_instances"); ok {
			request["MinReadyInstances"] = StringPointer(strconv.Itoa(v.(int)))
		}
	}

	if d.HasChange("min_ready_instance_ratio") {
		update = true
		if v, ok := d.GetOk("min_ready_instance_ratio"); ok {
			request["MinReadyInstanceRatio"] = StringPointer(strconv.Itoa(v.(int)))
		}
	}

	if update {
		action := "/pop/v1/sam/scale/applicationScalingRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
			if err != nil {
				if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("min_ready_instances")
		d.SetPartial("min_ready_instance_ratio")
		d.SetPartial("scaling_rule_timer")
		d.SetPartial("scaling_rule_metric")
	}

	if d.HasChange("scaling_rule_enable") {
		object, err := saeService.DescribeSaeApplicationScalingRule(d.Id())
		if err != nil {
			return WrapError(err)
		}
		target := fmt.Sprint(d.Get("scaling_rule_enable"))
		if fmt.Sprint(object["ScaleRuleEnabled"]) != target {

			if target == "true" {
				request := map[string]*string{
					"AppId":           StringPointer(parts[0]),
					"ScalingRuleName": StringPointer(parts[1]),
				}
				action := "/pop/v1/sam/scale/enableApplicationScalingRule"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
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
			if target == "false" {
				request := map[string]*string{
					"AppId":           StringPointer(parts[0]),
					"ScalingRuleName": StringPointer(parts[1]),
				}
				action := "/pop/v1/sam/scale/disableApplicationScalingRule"
				wait := incrementalWait(3*time.Second, 3*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("PUT"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
					if err != nil {
						if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
							wait()
							return resource.RetryableError(err)
						}
						return resource.NonRetryableError(err)
					}
					return nil
				})
				addDebug(action, response, request)
				if err != nil {
					return WrapErrorf(err, DefaultErrorMsg, d.Id(), "PUT "+action, AlibabaCloudSdkGoERROR)
				}
			}
			d.SetPartial("scaling_rule_enable")
		}
	}
	d.Partial(false)
	return resourceAlicloudSaeApplicationScalingRuleRead(d, meta)
}
func resourceAlicloudSaeApplicationScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "/pop/v1/sam/scale/applicationScalingRule"
	var response map[string]interface{}
	conn, err := client.NewServerlessClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]*string{
		"AppId":           StringPointer(parts[0]),
		"ScalingRuleName": StringPointer(parts[1]),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer("2019-05-06"), nil, StringPointer("DELETE"), StringPointer("AK"), StringPointer(action), request, nil, nil, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Application.ChangerOrderRunning"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DELETE "+action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
