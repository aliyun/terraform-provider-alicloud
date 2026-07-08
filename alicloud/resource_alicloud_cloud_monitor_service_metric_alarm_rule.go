// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCloudMonitorServiceMetricAlarmRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceMetricAlarmRuleCreate,
		Read:   resourceAliCloudCloudMonitorServiceMetricAlarmRuleRead,
		Update: resourceAliCloudCloudMonitorServiceMetricAlarmRuleUpdate,
		Delete: resourceAliCloudCloudMonitorServiceMetricAlarmRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"composite_expression": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"expression_raw": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expression_list_join": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"level": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expression_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"comparison_operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"period": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"statistics": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"contact_groups": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email_subject": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"escalations": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"critical": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"comparison_operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"pre_condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"statistics": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"info": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"comparison_operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"pre_condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"statistics": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"warn": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"comparison_operator": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"pre_condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"statistics": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"metric_alarm_rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"no_data_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"KEEP_LAST_STATE", "INSUFFICIENT_DATA", "OK"}, false),
			},
			"no_effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"prometheus": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"prom_ql": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"resources": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"send_ok": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"silence_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceMetricAlarmRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "PutResourceMetricRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("metric_alarm_rule_id"); ok {
		request["RuleId"] = v
	}

	request["Namespace"] = d.Get("namespace")
	if v, ok := d.GetOk("silence_time"); ok {
		request["SilenceTime"] = v
	}
	prometheus := make(map[string]interface{})

	if v := d.Get("prometheus"); !IsNil(v) {
		level1, _ := jsonpath.Get("$[0].level", v)
		if level1 != nil && level1 != "" {
			prometheus["Level"] = level1
		}
		localData, err := jsonpath.Get("$[0].annotations", v)
		if err != nil {
			localData = make([]interface{}, 0)
		}
		localMaps := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(localData) {
			dataLoopTmp := make(map[string]interface{})
			if dataLoop != nil {
				dataLoopTmp = dataLoop.(map[string]interface{})
			}
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Value"] = dataLoopTmp["value"]
			dataLoopMap["Key"] = dataLoopTmp["key"]
			localMaps = append(localMaps, dataLoopMap)
		}
		prometheus["Annotations"] = localMaps

		promQl, _ := jsonpath.Get("$[0].prom_ql", v)
		if promQl != nil && promQl != "" {
			prometheus["PromQL"] = promQl
		}
		times1, _ := jsonpath.Get("$[0].times", v)
		if times1 != nil && times1 != "" {
			prometheus["Times"] = times1
		}

		prometheusJson, err := json.Marshal(prometheus)
		if err != nil {
			return WrapError(err)
		}
		request["Prometheus"] = string(prometheusJson)
	}

	escalations := make(map[string]interface{})

	if v := d.Get("escalations"); !IsNil(v) {
		info := make(map[string]interface{})
		comparisonOperator1, _ := jsonpath.Get("$[0].info[0].comparison_operator", d.Get("escalations"))
		if comparisonOperator1 != nil && comparisonOperator1 != "" {
			info["ComparisonOperator"] = comparisonOperator1
		}
		statistics1, _ := jsonpath.Get("$[0].info[0].statistics", d.Get("escalations"))
		if statistics1 != nil && statistics1 != "" {
			info["Statistics"] = statistics1
		}
		threshold1, _ := jsonpath.Get("$[0].info[0].threshold", d.Get("escalations"))
		if threshold1 != nil && threshold1 != "" {
			info["Threshold"] = threshold1
		}
		times3, _ := jsonpath.Get("$[0].info[0].times", d.Get("escalations"))
		if times3 != nil && times3 != "" {
			info["Times"] = times3
		}

		if len(info) > 0 {
			escalations["Info"] = info
		}
		warn := make(map[string]interface{})
		statistics3, _ := jsonpath.Get("$[0].warn[0].statistics", d.Get("escalations"))
		if statistics3 != nil && statistics3 != "" {
			warn["Statistics"] = statistics3
		}
		threshold3, _ := jsonpath.Get("$[0].warn[0].threshold", d.Get("escalations"))
		if threshold3 != nil && threshold3 != "" {
			warn["Threshold"] = threshold3
		}
		comparisonOperator3, _ := jsonpath.Get("$[0].warn[0].comparison_operator", d.Get("escalations"))
		if comparisonOperator3 != nil && comparisonOperator3 != "" {
			warn["ComparisonOperator"] = comparisonOperator3
		}
		times5, _ := jsonpath.Get("$[0].warn[0].times", d.Get("escalations"))
		if times5 != nil && times5 != "" {
			warn["Times"] = times5
		}

		if len(warn) > 0 {
			escalations["Warn"] = warn
		}
		critical := make(map[string]interface{})
		statistics5, _ := jsonpath.Get("$[0].critical[0].statistics", d.Get("escalations"))
		if statistics5 != nil && statistics5 != "" {
			critical["Statistics"] = statistics5
		}
		threshold5, _ := jsonpath.Get("$[0].critical[0].threshold", d.Get("escalations"))
		if threshold5 != nil && threshold5 != "" {
			critical["Threshold"] = threshold5
		}
		comparisonOperator5, _ := jsonpath.Get("$[0].critical[0].comparison_operator", d.Get("escalations"))
		if comparisonOperator5 != nil && comparisonOperator5 != "" {
			critical["ComparisonOperator"] = comparisonOperator5
		}
		times7, _ := jsonpath.Get("$[0].critical[0].times", d.Get("escalations"))
		if times7 != nil && times7 != "" {
			critical["Times"] = times7
		}

		if len(critical) > 0 {
			escalations["Critical"] = critical
		}

		request["Escalations"] = escalations
	}

	compositeExpression := make(map[string]interface{})

	if v := d.Get("composite_expression"); !IsNil(v) {
		expressionListJoin1, _ := jsonpath.Get("$[0].expression_list_join", v)
		if expressionListJoin1 != nil && expressionListJoin1 != "" {
			compositeExpression["ExpressionListJoin"] = expressionListJoin1
		}
		localData1, err := jsonpath.Get("$[0].expression_list", v)
		if err != nil {
			localData1 = make([]interface{}, 0)
		}
		localMaps1 := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(localData1) {
			dataLoop1Tmp := make(map[string]interface{})
			if dataLoop1 != nil {
				dataLoop1Tmp = dataLoop1.(map[string]interface{})
			}
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Statistics"] = dataLoop1Tmp["statistics"]
			dataLoop1Map["ComparisonOperator"] = dataLoop1Tmp["comparison_operator"]
			dataLoop1Map["MetricName"] = dataLoop1Tmp["metric_name"]
			dataLoop1Map["Period"] = dataLoop1Tmp["period"]
			dataLoop1Map["Threshold"] = dataLoop1Tmp["threshold"]
			localMaps1 = append(localMaps1, dataLoop1Map)
		}
		compositeExpression["ExpressionList"] = localMaps1

		level3, _ := jsonpath.Get("$[0].level", v)
		if level3 != nil && level3 != "" {
			compositeExpression["Level"] = level3
		}
		expressionRaw1, _ := jsonpath.Get("$[0].expression_raw", v)
		if expressionRaw1 != nil && expressionRaw1 != "" {
			compositeExpression["ExpressionRaw"] = expressionRaw1
		}
		times9, _ := jsonpath.Get("$[0].times", v)
		if times9 != nil && times9 != "" {
			compositeExpression["Times"] = times9
		}

		compositeExpressionJson, err := json.Marshal(compositeExpression)
		if err != nil {
			return WrapError(err)
		}
		request["CompositeExpression"] = string(compositeExpressionJson)
	}

	if v, ok := d.GetOk("no_data_policy"); ok {
		request["NoDataPolicy"] = v
	}
	if v, ok := d.GetOk("effective_interval"); ok {
		request["EffectiveInterval"] = v
	}
	if v, ok := d.GetOk("webhook"); ok {
		request["Webhook"] = v
	}
	if v, ok := d.GetOk("interval"); ok {
		request["Interval"] = v
	}
	if v, ok := d.GetOk("email_subject"); ok {
		request["EmailSubject"] = v
	}
	if v, ok := d.GetOk("no_effective_interval"); ok {
		request["NoEffectiveInterval"] = v
	}
	request["ContactGroups"] = d.Get("contact_groups")
	request["Resources"] = d.Get("resources")
	if v, ok := d.GetOkExists("send_ok"); ok {
		request["SendOK"] = v
	}
	if v, ok := d.GetOk("labels"); ok {
		labelsMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range convertToInterfaceArray(v) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["Value"] = dataLoop2Tmp["value"]
			dataLoop2Map["Key"] = dataLoop2Tmp["key"]
			labelsMapsArray = append(labelsMapsArray, dataLoop2Map)
		}
		request["Labels"] = labelsMapsArray
	}

	request["RuleName"] = d.Get("rule_name")
	request["MetricName"] = d.Get("metric_name")
	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}
	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_monitor_service_metric_alarm_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["RuleId"]))

	return resourceAliCloudCloudMonitorServiceMetricAlarmRuleUpdate(d, meta)
}

func resourceAliCloudCloudMonitorServiceMetricAlarmRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	objectRaw, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceMetricAlarmRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_metric_alarm_rule DescribeCloudMonitorServiceMetricAlarmRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("contact_groups", objectRaw["ContactGroups"])
	d.Set("dimensions", objectRaw["Dimensions"])
	d.Set("effective_interval", objectRaw["EffectiveInterval"])
	d.Set("email_subject", objectRaw["MailSubject"])
	d.Set("metric_name", objectRaw["MetricName"])
	d.Set("namespace", objectRaw["Namespace"])
	d.Set("no_data_policy", objectRaw["NoDataPolicy"])
	d.Set("no_effective_interval", objectRaw["NoEffectiveInterval"])
	d.Set("period", objectRaw["Period"])
	d.Set("resources", objectRaw["Resources"])
	d.Set("rule_name", objectRaw["RuleName"])
	d.Set("send_ok", objectRaw["SendOK"])
	d.Set("silence_time", objectRaw["SilenceTime"])
	d.Set("source_type", objectRaw["SourceType"])
	d.Set("status", objectRaw["EnableState"])
	d.Set("webhook", objectRaw["Webhook"])
	d.Set("interval", objectRaw["Interval"])
	d.Set("metric_alarm_rule_id", objectRaw["RuleId"])

	compositeExpressionMaps := make([]map[string]interface{}, 0)
	compositeExpressionMap := make(map[string]interface{})
	compositeExpressionRaw := make(map[string]interface{})
	if objectRaw["CompositeExpression"] != nil {
		compositeExpressionRaw = objectRaw["CompositeExpression"].(map[string]interface{})
	}
	if len(compositeExpressionRaw) > 0 {
		compositeExpressionMap["expression_list_join"] = compositeExpressionRaw["ExpressionListJoin"]
		compositeExpressionMap["expression_raw"] = compositeExpressionRaw["ExpressionRaw"]
		compositeExpressionMap["level"] = compositeExpressionRaw["Level"]
		compositeExpressionMap["times"] = compositeExpressionRaw["Times"]

		expressionListRaw, _ := jsonpath.Get("$.CompositeExpression.ExpressionList.ExpressionList", objectRaw)
		expressionListMaps := make([]map[string]interface{}, 0)
		if expressionListRaw != nil {
			for _, expressionListChildRaw := range convertToInterfaceArray(expressionListRaw) {
				expressionListMap := make(map[string]interface{})
				expressionListChildRaw := expressionListChildRaw.(map[string]interface{})
				expressionListMap["comparison_operator"] = expressionListChildRaw["ComparisonOperator"]
				expressionListMap["metric_name"] = expressionListChildRaw["MetricName"]
				expressionListMap["period"] = expressionListChildRaw["Period"]
				expressionListMap["statistics"] = expressionListChildRaw["Statistics"]
				expressionListMap["threshold"] = expressionListChildRaw["Threshold"]

				expressionListMaps = append(expressionListMaps, expressionListMap)
			}
		}
		compositeExpressionMap["expression_list"] = expressionListMaps
		compositeExpressionMaps = append(compositeExpressionMaps, compositeExpressionMap)
	}
	if err := d.Set("composite_expression", compositeExpressionMaps); err != nil {
		return err
	}
	escalationsMaps := make([]map[string]interface{}, 0)
	escalationsMap := make(map[string]interface{})
	escalationsRaw := make(map[string]interface{})
	if objectRaw["Escalations"] != nil {
		escalationsRaw = objectRaw["Escalations"].(map[string]interface{})
	}
	if len(escalationsRaw) > 0 {

		criticalMaps := make([]map[string]interface{}, 0)
		criticalMap := make(map[string]interface{})
		criticalRaw := make(map[string]interface{})
		if escalationsRaw["Critical"] != nil {
			criticalRaw = escalationsRaw["Critical"].(map[string]interface{})
		}
		// PROMETHEUS-style rules cause the server to inject synthetic Critical
		// entries (Statistics="value", Times=prometheus.times) with empty
		// ComparisonOperator/Threshold. Treat those as "not user-configured" so
		// Read does not surface them as escalations on resources that never set
		// the block.
		hasRealEscalation := func(raw map[string]interface{}) bool {
			return fmt.Sprint(raw["ComparisonOperator"]) != "" || fmt.Sprint(raw["Threshold"]) != ""
		}
		if len(criticalRaw) > 0 && hasRealEscalation(criticalRaw) {
			criticalMap["comparison_operator"] = criticalRaw["ComparisonOperator"]
			criticalMap["pre_condition"] = criticalRaw["PreCondition"]
			criticalMap["statistics"] = criticalRaw["Statistics"]
			criticalMap["threshold"] = criticalRaw["Threshold"]
			criticalMap["times"] = criticalRaw["Times"]

			criticalMaps = append(criticalMaps, criticalMap)
		}
		escalationsMap["critical"] = criticalMaps
		infoMaps := make([]map[string]interface{}, 0)
		infoMap := make(map[string]interface{})
		infoRaw := make(map[string]interface{})
		if escalationsRaw["Info"] != nil {
			infoRaw = escalationsRaw["Info"].(map[string]interface{})
		}
		if len(infoRaw) > 0 && hasRealEscalation(infoRaw) {
			infoMap["comparison_operator"] = infoRaw["ComparisonOperator"]
			infoMap["pre_condition"] = infoRaw["PreCondition"]
			infoMap["statistics"] = infoRaw["Statistics"]
			infoMap["threshold"] = infoRaw["Threshold"]
			infoMap["times"] = infoRaw["Times"]

			infoMaps = append(infoMaps, infoMap)
		}
		escalationsMap["info"] = infoMaps
		warnMaps := make([]map[string]interface{}, 0)
		warnMap := make(map[string]interface{})
		warnRaw := make(map[string]interface{})
		if escalationsRaw["Warn"] != nil {
			warnRaw = escalationsRaw["Warn"].(map[string]interface{})
		}
		if len(warnRaw) > 0 && hasRealEscalation(warnRaw) {
			warnMap["comparison_operator"] = warnRaw["ComparisonOperator"]
			warnMap["pre_condition"] = warnRaw["PreCondition"]
			warnMap["statistics"] = warnRaw["Statistics"]
			warnMap["threshold"] = warnRaw["Threshold"]
			warnMap["times"] = warnRaw["Times"]

			warnMaps = append(warnMaps, warnMap)
		}
		escalationsMap["warn"] = warnMaps
		if len(criticalMaps) > 0 || len(infoMaps) > 0 || len(warnMaps) > 0 {
			escalationsMaps = append(escalationsMaps, escalationsMap)
		}
	}
	if err := d.Set("escalations", escalationsMaps); err != nil {
		return err
	}
	labelsRaw, _ := jsonpath.Get("$.Labels.Labels", objectRaw)
	labelsMaps := make([]map[string]interface{}, 0)
	if labelsRaw != nil {
		for _, labelsChildRaw := range convertToInterfaceArray(labelsRaw) {
			labelsMap := make(map[string]interface{})
			labelsChildRaw := labelsChildRaw.(map[string]interface{})
			labelsMap["key"] = labelsChildRaw["Key"]
			labelsMap["value"] = labelsChildRaw["Value"]

			labelsMaps = append(labelsMaps, labelsMap)
		}
	}
	if err := d.Set("labels", labelsMaps); err != nil {
		return err
	}
	prometheusMaps := make([]map[string]interface{}, 0)
	prometheusMap := make(map[string]interface{})
	prometheusRaw := make(map[string]interface{})
	if objectRaw["Prometheus"] != nil {
		prometheusRaw = objectRaw["Prometheus"].(map[string]interface{})
	}
	if len(prometheusRaw) > 0 {
		prometheusMap["level"] = prometheusRaw["Level"]
		prometheusMap["prom_ql"] = prometheusRaw["PromQL"]
		prometheusMap["times"] = prometheusRaw["Times"]

		annotationsRaw, _ := jsonpath.Get("$.Prometheus.Annotations.Annotations", objectRaw)
		annotationsMaps := make([]map[string]interface{}, 0)
		if annotationsRaw != nil {
			for _, annotationsChildRaw := range convertToInterfaceArray(annotationsRaw) {
				annotationsMap := make(map[string]interface{})
				annotationsChildRaw := annotationsChildRaw.(map[string]interface{})
				annotationsMap["key"] = annotationsChildRaw["Key"]
				annotationsMap["value"] = annotationsChildRaw["Value"]

				annotationsMaps = append(annotationsMaps, annotationsMap)
			}
		}
		prometheusMap["annotations"] = annotationsMaps
		prometheusMaps = append(prometheusMaps, prometheusMap)
	}
	if err := d.Set("prometheus", prometheusMaps); err != nil {
		return err
	}

	d.Set("metric_alarm_rule_id", d.Id())

	return nil
}

func resourceAliCloudCloudMonitorServiceMetricAlarmRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}
	objectRaw, _ := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceMetricAlarmRule(d.Id())

	initedStatus := false
	if _, ok := d.GetOkExists("status"); ok && d.IsNewResource() {
		initedStatus = true
	}
	if initedStatus || d.HasChange("status") {
		var err error
		target := d.Get("status").(bool)

		currentStatus, err := jsonpath.Get("EnableState", objectRaw)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, d.Id(), "EnableState", objectRaw)
		}
		if formatBool(currentStatus) != target {
			if target == true {
				action := "EnableMetricRules"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RuleId.1"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
			if target == false {
				action := "DisableMetricRules"
				request = make(map[string]interface{})
				query = make(map[string]interface{})
				request["RuleId.1"] = d.Id()

				wait := incrementalWait(3*time.Second, 5*time.Second)
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
		}
	}

	var err error
	update = false
	action := "PutResourceMetricRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["RuleId"] = d.Id()

	// PutResourceMetricRule is a full-replace PUT: any field omitted in the
	// request body is reset to the server's default value (confirmed by CMS team
	// on Aone 83358430). Send the complete resource definition on every
	// invocation — including the Create→Update follow-up call — so that fields
	// already written by Create aren't blanked out by the second PUT.
	request["Namespace"] = d.Get("namespace")
	request["SilenceTime"] = d.Get("silence_time")
	request["NoDataPolicy"] = d.Get("no_data_policy")
	request["EffectiveInterval"] = d.Get("effective_interval")
	request["Webhook"] = d.Get("webhook")
	request["EmailSubject"] = d.Get("email_subject")
	request["NoEffectiveInterval"] = d.Get("no_effective_interval")
	request["ContactGroups"] = d.Get("contact_groups")
	request["Resources"] = d.Get("resources")
	request["SendOK"] = d.Get("send_ok")
	if v, ok := d.GetOk("interval"); ok {
		request["Interval"] = v
	}
	request["RuleName"] = d.Get("rule_name")
	request["MetricName"] = d.Get("metric_name")
	request["Period"] = d.Get("period")

	if v := d.Get("prometheus"); v != nil {
		prometheus := make(map[string]interface{})

		level1, _ := jsonpath.Get("$[0].level", v)
		if level1 != nil && level1 != "" {
			prometheus["Level"] = level1
		}
		localData, err := jsonpath.Get("$[0].annotations", v)
		if err != nil {
			localData = make([]interface{}, 0)
		}
		localMaps := make([]interface{}, 0)
		for _, dataLoop := range convertToInterfaceArray(localData) {
			dataLoopTmp := make(map[string]interface{})
			if dataLoop != nil {
				dataLoopTmp = dataLoop.(map[string]interface{})
			}
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Value"] = dataLoopTmp["value"]
			dataLoopMap["Key"] = dataLoopTmp["key"]
			localMaps = append(localMaps, dataLoopMap)
		}
		prometheus["Annotations"] = localMaps

		promQl, _ := jsonpath.Get("$[0].prom_ql", v)
		if promQl != nil && promQl != "" {
			prometheus["PromQL"] = promQl
		}
		times1, _ := jsonpath.Get("$[0].times", v)
		if times1 != nil && times1 != "" {
			prometheus["Times"] = times1
		}

		if len(prometheus) > 0 {
			prometheusJson, err := json.Marshal(prometheus)
			if err != nil {
				return WrapError(err)
			}
			request["Prometheus"] = string(prometheusJson)
		}
	}

	if v := d.Get("escalations"); v != nil {
		escalations := make(map[string]interface{})

		info := make(map[string]interface{})
		comparisonOperator1, _ := jsonpath.Get("$[0].info[0].comparison_operator", v)
		if comparisonOperator1 != nil && comparisonOperator1 != "" {
			info["ComparisonOperator"] = comparisonOperator1
		}
		statistics1, _ := jsonpath.Get("$[0].info[0].statistics", v)
		if statistics1 != nil && statistics1 != "" {
			info["Statistics"] = statistics1
		}
		threshold1, _ := jsonpath.Get("$[0].info[0].threshold", v)
		if threshold1 != nil && threshold1 != "" {
			info["Threshold"] = threshold1
		}
		times3, _ := jsonpath.Get("$[0].info[0].times", v)
		if times3 != nil && times3 != "" {
			info["Times"] = times3
		}
		if len(info) > 0 {
			escalations["Info"] = info
		}

		warn := make(map[string]interface{})
		statistics3, _ := jsonpath.Get("$[0].warn[0].statistics", v)
		if statistics3 != nil && statistics3 != "" {
			warn["Statistics"] = statistics3
		}
		threshold3, _ := jsonpath.Get("$[0].warn[0].threshold", v)
		if threshold3 != nil && threshold3 != "" {
			warn["Threshold"] = threshold3
		}
		comparisonOperator3, _ := jsonpath.Get("$[0].warn[0].comparison_operator", v)
		if comparisonOperator3 != nil && comparisonOperator3 != "" {
			warn["ComparisonOperator"] = comparisonOperator3
		}
		times5, _ := jsonpath.Get("$[0].warn[0].times", v)
		if times5 != nil && times5 != "" {
			warn["Times"] = times5
		}
		if len(warn) > 0 {
			escalations["Warn"] = warn
		}

		critical := make(map[string]interface{})
		statistics5, _ := jsonpath.Get("$[0].critical[0].statistics", v)
		if statistics5 != nil && statistics5 != "" {
			critical["Statistics"] = statistics5
		}
		threshold5, _ := jsonpath.Get("$[0].critical[0].threshold", v)
		if threshold5 != nil && threshold5 != "" {
			critical["Threshold"] = threshold5
		}
		comparisonOperator5, _ := jsonpath.Get("$[0].critical[0].comparison_operator", v)
		if comparisonOperator5 != nil && comparisonOperator5 != "" {
			critical["ComparisonOperator"] = comparisonOperator5
		}
		times7, _ := jsonpath.Get("$[0].critical[0].times", v)
		if times7 != nil && times7 != "" {
			critical["Times"] = times7
		}
		if len(critical) > 0 {
			escalations["Critical"] = critical
		}

		if len(escalations) > 0 {
			request["Escalations"] = escalations
		}
	}

	if v := d.Get("composite_expression"); !IsNil(v) {
		compositeExpression := make(map[string]interface{})

		expressionListJoin1, _ := jsonpath.Get("$[0].expression_list_join", v)
		if expressionListJoin1 != nil && expressionListJoin1 != "" {
			compositeExpression["ExpressionListJoin"] = expressionListJoin1
		}
		localData1, err := jsonpath.Get("$[0].expression_list", v)
		if err != nil {
			localData1 = make([]interface{}, 0)
		}
		localMaps1 := make([]interface{}, 0)
		for _, dataLoop1 := range convertToInterfaceArray(localData1) {
			dataLoop1Tmp := make(map[string]interface{})
			if dataLoop1 != nil {
				dataLoop1Tmp = dataLoop1.(map[string]interface{})
			}
			dataLoop1Map := make(map[string]interface{})
			dataLoop1Map["Statistics"] = dataLoop1Tmp["statistics"]
			dataLoop1Map["ComparisonOperator"] = dataLoop1Tmp["comparison_operator"]
			dataLoop1Map["MetricName"] = dataLoop1Tmp["metric_name"]
			dataLoop1Map["Period"] = dataLoop1Tmp["period"]
			dataLoop1Map["Threshold"] = dataLoop1Tmp["threshold"]
			localMaps1 = append(localMaps1, dataLoop1Map)
		}
		if len(localMaps1) > 0 {
			compositeExpression["ExpressionList"] = localMaps1
		}

		level3, _ := jsonpath.Get("$[0].level", v)
		if level3 != nil && level3 != "" {
			compositeExpression["Level"] = level3
		}
		expressionRaw1, _ := jsonpath.Get("$[0].expression_raw", v)
		if expressionRaw1 != nil && expressionRaw1 != "" {
			compositeExpression["ExpressionRaw"] = expressionRaw1
		}
		times9, _ := jsonpath.Get("$[0].times", v)
		if times9 != nil && times9 != "" {
			compositeExpression["Times"] = times9
		}

		if len(compositeExpression) > 0 {
			compositeExpressionJson, err := json.Marshal(compositeExpression)
			if err != nil {
				return WrapError(err)
			}
			request["CompositeExpression"] = string(compositeExpressionJson)
		}
	}

	if v, ok := d.GetOk("labels"); ok {
		labelsMapsArray := make([]interface{}, 0)
		for _, dataLoop2 := range convertToInterfaceArray(v) {
			dataLoop2Tmp := dataLoop2.(map[string]interface{})
			dataLoop2Map := make(map[string]interface{})
			dataLoop2Map["Value"] = dataLoop2Tmp["value"]
			dataLoop2Map["Key"] = dataLoop2Tmp["key"]
			labelsMapsArray = append(labelsMapsArray, dataLoop2Map)
		}
		request["Labels"] = labelsMapsArray
	}

	update = true

	if update {
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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

	return resourceAliCloudCloudMonitorServiceMetricAlarmRuleRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceMetricAlarmRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMetricRules"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["Id.1"] = d.Id()

	wait := incrementalWait(3*time.Second, 5*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, query, request, true)
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
