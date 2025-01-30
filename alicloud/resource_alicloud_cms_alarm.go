package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
)

func resourceAliCloudCmsAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsAlarmCreate,
		Read:   resourceAliCloudCmsAlarmRead,
		Update: resourceAliCloudCmsAlarmUpdate,
		Delete: resourceAliCloudCmsAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"contact_groups": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"metric_dimensions": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"dimensions"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareArrayJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"effective_interval": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "00:00-23:59",
				ConflictsWith: []string{"start_time", "end_time"},
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  300,
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      86400,
				ValidateFunc: IntBetween(300, 86400),
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"tags": tagsSchema(),
			"escalations_critical": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  MoreThan,
							ValidateFunc: StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual, Equal,
								"GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek",
								"LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod",
							}, false),
						},
						"statistics": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Average,
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
			},
			"escalations_info": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  MoreThan,
							ValidateFunc: StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual, Equal,
								"GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek",
								"LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod",
							}, false),
						},
						"statistics": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Average,
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
			},
			"escalations_warn": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"comparison_operator": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  MoreThan,
							ValidateFunc: StringInSlice([]string{
								MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual, Equal,
								"GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek",
								"LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod",
							}, false),
						},
						"statistics": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  Average,
						},
						"threshold": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
			},
			"prometheus": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prom_ql": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"level": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: StringInSlice([]string{"Critical", "Warn", "Info"}, false),
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem:     schema.TypeString,
						},
					},
				},
			},
			"targets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"json_params": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"level": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"Critical", "Warn", "Info"}, false),
						},
						"arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"composite_expression": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"level": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"CRITICAL", "WARN", "INFO"}, false),
						},
						"times": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"expression_raw": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expression_list_join": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: StringInSlice([]string{"&&", "||"}, false),
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
										ValidateFunc: StringInSlice([]string{
											MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual, Equal,
											"GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek",
											"LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod",
										}, false),
									},
									"statistics": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"period": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dimensions": {
				Type:          schema.TypeMap,
				Optional:      true,
				Computed:      true,
				Elem:          schema.TypeString,
				ConflictsWith: []string{"metric_dimensions"},
				Deprecated:    "Field `dimensions` has been deprecated from provider version 1.173.0. New field `metric_dimensions` instead.",
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:      0,
				//ValidateFunc: IntBetween(0, 24),
				ConflictsWith: []string{"effective_interval"},
				Deprecated:    "Field `start_time` has been deprecated from provider version 1.50.0. New field `effective_interval` instead.",
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
				//Default:      24,
				//ValidateFunc: IntBetween(0, 24),
				ConflictsWith: []string{"effective_interval"},
				Deprecated:    "Field `end_time` has been deprecated from provider version 1.50.0. New field `effective_interval` instead.",
			},
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: StringInSlice([]string{
					MoreThan, MoreThanOrEqual, LessThan, LessThanOrEqual, NotEqual,
				}, false),
				Removed: "Field `operator` has been removed from provider version 1.216.0. New field `escalations_critical.comparison_operator` instead.",
			},
			"statistics": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Field `statistics` has been removed from provider version 1.216.0. New field `escalations_critical.statistics` instead.",
			},
			"threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Removed:  "Field `threshold` has been removed from provider version 1.216.0. New field `escalations_critical.threshold` instead.",
			},
			"triggered_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Removed:  "Field `triggered_count` has been removed from provider version 1.216.0. New field `escalations_critical.times` instead.",
			},
			"notify_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntInSlice([]int{0, 1}),
				Removed:      "Field `notify_type` has been removed from provider version 1.50.0.",
			},
		},
	}
}

func resourceAliCloudCmsAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(resource.UniqueId())
	return resourceAliCloudCmsAlarmUpdate(d, meta)
}

func resourceAliCloudCmsAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	object, err := cmsService.DescribeAlarm(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object["RuleName"])
	d.Set("project", object["Namespace"])
	d.Set("metric", object["MetricName"])
	d.Set("contact_groups", strings.Split(object["ContactGroups"].(string), ","))
	d.Set("effective_interval", object["EffectiveInterval"])
	d.Set("period", object["Period"])
	d.Set("silence_time", object["SilenceTime"])
	d.Set("webhook", object["Webhook"])
	d.Set("enabled", object["EnableState"])
	d.Set("status", object["AlertState"])

	if err := d.Set("metric_dimensions", object["Resources"]); err != nil {
		return WrapError(err)
	}

	if tags, ok := object["Labels"]; ok {
		tagsArg := tags.(map[string]interface{})

		if labels, ok := tagsArg["Labels"]; ok {
			d.Set("tags", tagsToMap(labels))
		}
	}

	if escalations, ok := object["Escalations"].(map[string]interface{}); ok {
		if critical, ok := escalations["Critical"].(map[string]interface{}); ok {
			mapping := map[string]interface{}{
				"statistics":          critical["Statistics"],
				"comparison_operator": convertCmsAlarmComparisonOperator(fmt.Sprint(critical["ComparisonOperator"])),
				"threshold":           critical["Threshold"],
				"times":               formatInt(critical["Times"]),
			}
			d.Set("escalations_critical", []map[string]interface{}{mapping})
		}

		if info, ok := escalations["Info"].(map[string]interface{}); ok {
			mappingInfo := map[string]interface{}{
				"statistics":          info["Statistics"],
				"comparison_operator": convertCmsAlarmComparisonOperator(fmt.Sprint(info["ComparisonOperator"])),
				"threshold":           info["Threshold"],
				"times":               formatInt(info["Times"]),
			}
			d.Set("escalations_info", []map[string]interface{}{mappingInfo})
		}

		if warn, ok := escalations["Warn"].(map[string]interface{}); ok {
			mappingWarn := map[string]interface{}{
				"statistics":          warn["Statistics"],
				"comparison_operator": convertCmsAlarmComparisonOperator(fmt.Sprint(warn["ComparisonOperator"])),
				"threshold":           warn["Threshold"],
				"times":               formatInt(warn["Times"]),
			}
			d.Set("escalations_warn", []map[string]interface{}{mappingWarn})
		}
	}

	if v, ok := object["Prometheus"]; ok {
		if prometheus, ok := v.(map[string]interface{}); ok && len(prometheus) > 0 {

			prometheusList := make([]map[string]interface{}, 0)
			mapping := map[string]interface{}{
				"prom_ql": prometheus["PromQL"],
			}

			if v, ok := prometheus["Level"]; ok {
				mapping["level"] = convertCmsAlarmPrometheusLevelResponse(v.(string))
			}

			if v, ok := prometheus["Times"]; ok {
				mapping["times"] = v
			}

			annotationsMap := make(map[string]interface{}, 0)
			if v, ok := prometheus["Annotations"]; ok {
				annotations := v.(map[string]interface{})
				if v, ok := annotations["Annotations"]; ok && len(v.([]interface{})) > 0 {
					for _, item := range v.([]interface{}) {
						itemArg := item.(map[string]interface{})
						annotationsMap[itemArg["Key"].(string)] = itemArg["Value"]
					}
				}
			}

			mapping["annotations"] = annotationsMap
			prometheusList = append(prometheusList, mapping)

			if err := d.Set("prometheus", prometheusList); err != nil {
				return WrapError(err)
			}
		}
	}

	if compositeExpression, ok := object["CompositeExpression"]; ok {
		compositeExpressionMaps := make([]map[string]interface{}, 0)
		compositeExpressionArg := compositeExpression.(map[string]interface{})
		compositeExpressionMap := make(map[string]interface{})

		if level, ok := compositeExpressionArg["Level"]; ok {
			compositeExpressionMap["level"] = level
		}

		if times, ok := compositeExpressionArg["Times"]; ok {
			compositeExpressionMap["times"] = times
		}

		if expressionRaw, ok := compositeExpressionArg["ExpressionRaw"]; ok {
			compositeExpressionMap["expression_raw"] = expressionRaw
		}

		if expressionListJoin, ok := compositeExpressionArg["ExpressionListJoin"]; ok {
			compositeExpressionMap["expression_list_join"] = expressionListJoin
		}

		if expressionList, ok := compositeExpressionArg["ExpressionList"]; ok {
			if expressionLists, ok := expressionList.(map[string]interface{})["ExpressionList"]; ok {
				expressionListMaps := make([]map[string]interface{}, 0)
				for _, v := range expressionLists.([]interface{}) {
					expressionListArg := v.(map[string]interface{})
					expressionListMap := map[string]interface{}{}

					if metricName, ok := expressionListArg["MetricName"]; ok {
						expressionListMap["metric_name"] = metricName
					}

					if comparisonOperator, ok := expressionListArg["ComparisonOperator"]; ok {
						expressionListMap["comparison_operator"] = convertCmsAlarmComparisonOperator(fmt.Sprint(comparisonOperator))
					}

					if statistics, ok := expressionListArg["Statistics"]; ok {
						expressionListMap["statistics"] = statistics
					}

					if threshold, ok := expressionListArg["Threshold"]; ok {
						expressionListMap["threshold"] = threshold
					}

					if period, ok := expressionListArg["Period"]; ok {
						expressionListMap["period"] = period
					}

					expressionListMaps = append(expressionListMaps, expressionListMap)
				}

				compositeExpressionMap["expression_list"] = expressionListMaps
			}
		}

		if len(compositeExpressionMap) > 0 {
			compositeExpressionMaps = append(compositeExpressionMaps, compositeExpressionMap)
		}

		d.Set("composite_expression", compositeExpressionMaps)
	}

	dims := make([]map[string]interface{}, 0)
	if fmt.Sprint(object["Resources"]) != "" {
		if err := json.Unmarshal([]byte(object["Resources"].(string)), &dims); err != nil {
			return fmt.Errorf("Unmarshaling Dimensions got an error: %#v.", err)
		}
	}

	dimensionList := make(map[string]interface{}, 0)
	for _, raw := range dims {
		for k, v := range raw {
			if dimensionListValue, ok := dimensionList[k]; ok {
				dimensionList[k] = fmt.Sprint(dimensionListValue.(string), ",", v.(string))
			} else {
				dimensionList[k] = v
			}
		}
	}

	if err := d.Set("dimensions", dimensionList); err != nil {
		return WrapError(err)
	}

	targetsList, err := cmsService.DescribeMetricRuleTargets(d.Id())
	if err != nil {
		if NotFoundError(err) {
			return nil
		}
		return WrapError(err)
	}

	targetsMaps := make([]map[string]interface{}, 0)
	for _, targets := range targetsList {
		targetsArg := targets.(map[string]interface{})
		targetsMap := map[string]interface{}{}

		if id, ok := targetsArg["Id"]; ok {
			targetsMap["target_id"] = id
		}

		if jsonParams, ok := targetsArg["JsonParams"]; ok {
			targetsMap["json_params"] = jsonParams
		}

		if level, ok := targetsArg["Level"]; ok {
			targetsMap["level"] = level
		}

		if arn, ok := targetsArg["Arn"]; ok {
			targetsMap["arn"] = arn
		}

		targetsMaps = append(targetsMaps, targetsMap)
	}

	d.Set("targets", targetsMaps)

	return nil
}

func resourceAliCloudCmsAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	var response map[string]interface{}
	d.Partial(true)

	action := "PutResourceMetricRule"
	request := make(map[string]interface{})
	var err error

	request["RuleId"] = d.Id()
	request["RuleName"] = d.Get("name").(string)
	request["Namespace"] = d.Get("project").(string)
	request["MetricName"] = d.Get("metric").(string)
	request["ContactGroups"] = strings.Join(expandStringList(d.Get("contact_groups").([]interface{})), ",")

	if v, ok := d.GetOk("metric_dimensions"); ok && v.(string) != "" {
		request["Resources"] = v.(string)
	} else {
		var dimList []map[string]string
		if dimensions, ok := d.GetOk("dimensions"); ok {
			for k, v := range dimensions.(map[string]interface{}) {
				values := strings.Split(v.(string), COMMA_SEPARATED)
				if len(values) > 0 {
					for _, vv := range values {
						dimList = append(dimList, map[string]string{k: Trim(vv)})
					}
				} else {
					dimList = append(dimList, map[string]string{k: Trim(v.(string))})
				}

			}
		}

		if len(dimList) > 0 {
			if bytes, err := json.Marshal(dimList); err != nil {
				return fmt.Errorf("Marshaling dimensions to json string got an error: %#v.", err)
			} else {
				request["Resources"] = string(bytes[:])
			}
		}
	}

	if v, ok := d.GetOk("effective_interval"); ok && v.(string) != "" {
		request["EffectiveInterval"] = v.(string)
	} else {
		start, startOk := d.GetOk("start_time")
		end, endOk := d.GetOk("end_time")
		if startOk && endOk && end.(int) > 0 {
			// The EffectiveInterval valid value between 00:00 and 23:59
			request["EffectiveInterval"] = fmt.Sprintf("%d:00-%d:59", start.(int), end.(int)-1)
		}
	}

	if v, ok := d.GetOkExists("period"); ok {
		request["Period"] = strconv.Itoa(v.(int))
	}

	if v, ok := d.GetOkExists("silence_time"); ok {
		request["SilenceTime"] = requests.NewInteger(v.(int))
	}

	if v, ok := d.GetOk("webhook"); ok && v.(string) != "" {
		request["Webhook"] = v.(string)
	}

	// Critical
	if v, ok := d.GetOk("escalations_critical"); ok && len(v.([]interface{})) != 0 {

		for _, escalationsCriticalList := range v.([]interface{}) {
			escalationsCriticalArg := escalationsCriticalList.(map[string]interface{})

			if comparisonOperator, ok := escalationsCriticalArg["comparison_operator"]; ok {
				request["Escalations.Critical.ComparisonOperator"] = convertCmsAlarmComparisonOperator(fmt.Sprint(comparisonOperator))
			}

			if statistics, ok := escalationsCriticalArg["statistics"]; ok {
				request["Escalations.Critical.Statistics"] = statistics
			}

			if threshold, ok := escalationsCriticalArg["threshold"]; ok {
				request["Escalations.Critical.Threshold"] = threshold
			}

			if times, ok := escalationsCriticalArg["times"]; ok {
				request["Escalations.Critical.Times"] = times
			}
		}
	}

	// Info
	if v, ok := d.GetOk("escalations_info"); ok && len(v.([]interface{})) != 0 {
		for _, escalationsInfoList := range v.([]interface{}) {
			escalationsInfoArg := escalationsInfoList.(map[string]interface{})

			if comparisonOperator, ok := escalationsInfoArg["comparison_operator"]; ok {
				request["Escalations.Info.ComparisonOperator"] = convertCmsAlarmComparisonOperator(fmt.Sprint(comparisonOperator))
			}

			if statistics, ok := escalationsInfoArg["statistics"]; ok {
				request["Escalations.Info.Statistics"] = statistics
			}

			if threshold, ok := escalationsInfoArg["threshold"]; ok {
				request["Escalations.Info.Threshold"] = threshold
			}

			if times, ok := escalationsInfoArg["times"]; ok {
				request["Escalations.Info.Times"] = times
			}
		}
	}

	// Warn
	if v, ok := d.GetOk("escalations_warn"); ok && len(v.([]interface{})) != 0 {
		for _, escalationsWarnList := range v.([]interface{}) {
			escalationsWarnArg := escalationsWarnList.(map[string]interface{})

			if comparisonOperator, ok := escalationsWarnArg["comparison_operator"]; ok {
				request["Escalations.Warn.ComparisonOperator"] = convertCmsAlarmComparisonOperator(fmt.Sprint(comparisonOperator))
			}

			if statistics, ok := escalationsWarnArg["statistics"]; ok {
				request["Escalations.Warn.Statistics"] = statistics
			}

			if threshold, ok := escalationsWarnArg["threshold"]; ok {
				request["Escalations.Warn.Threshold"] = threshold
			}

			if times, ok := escalationsWarnArg["times"]; ok {
				request["Escalations.Warn.Times"] = times
			}
		}
	}

	if v, ok := d.GetOk("prometheus"); ok && len(v.([]interface{})) > 0 {
		prometheus := v.([]interface{})[0]
		prometheusMap := make(map[string]interface{}, 0)
		prometheusArg := prometheus.(map[string]interface{})
		prometheusMap["PromQL"] = prometheusArg["prom_ql"]
		prometheusMap["Level"] = prometheusArg["level"]
		prometheusMap["Times"] = prometheusArg["times"]
		if vv, ok := prometheusArg["annotations"]; ok {
			tags := make([]map[string]interface{}, 0)
			for key, value := range vv.(map[string]interface{}) {
				tags = append(tags, map[string]interface{}{
					"Key":   key,
					"Value": value,
				})
			}
			prometheusMap["Annotations"] = tags
		}

		request["Prometheus"], _ = convertMaptoJsonString(prometheusMap)
	}

	if v, ok := d.GetOk("composite_expression"); ok {
		compositeExpressionMap := map[string]interface{}{}
		for _, compositeExpression := range v.([]interface{}) {
			compositeExpressionArg := compositeExpression.(map[string]interface{})

			if level, ok := compositeExpressionArg["level"]; ok {
				compositeExpressionMap["Level"] = level
			}

			if times, ok := compositeExpressionArg["times"]; ok {
				compositeExpressionMap["Times"] = times
			}

			if expressionRaw, ok := compositeExpressionArg["expression_raw"]; ok {
				compositeExpressionMap["ExpressionRaw"] = expressionRaw
			}

			if expressionListJoin, ok := compositeExpressionArg["expression_list_join"]; ok {
				compositeExpressionMap["ExpressionListJoin"] = expressionListJoin
			}

			if expressionList, ok := compositeExpressionArg["expression_list"]; ok {
				expressionListMaps := make([]map[string]interface{}, 0)
				for _, expressionListArgList := range expressionList.([]interface{}) {
					expressionListMap := map[string]interface{}{}
					expressionListArg := expressionListArgList.(map[string]interface{})

					if metricName, ok := expressionListArg["metric_name"]; ok {
						expressionListMap["MetricName"] = metricName
					}

					if comparisonOperator, ok := expressionListArg["comparison_operator"]; ok {
						expressionListMap["ComparisonOperator"] = convertCmsAlarmComparisonOperator(fmt.Sprint(comparisonOperator))
					}

					if statistics, ok := expressionListArg["statistics"]; ok {
						expressionListMap["Statistics"] = statistics
					}

					if threshold, ok := expressionListArg["threshold"]; ok {
						expressionListMap["Threshold"] = threshold
					}

					if period, ok := expressionListArg["period"]; ok {
						expressionListMap["Period"] = period
					}

					expressionListMaps = append(expressionListMaps, expressionListMap)
				}

				compositeExpressionMap["ExpressionList"] = expressionListMaps
			}
		}

		compositeExpressionJson, err := convertMaptoJsonString(compositeExpressionMap)
		if err != nil {
			return WrapError(err)
		}

		request["CompositeExpression"] = compositeExpressionJson
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]map[string]interface{}, 0)
		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, map[string]interface{}{
				"Key":   key,
				"Value": value.(string),
			})
		}

		request["Labels"] = tags
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
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

	d.SetPartial("name")
	d.SetPartial("contact_groups")
	d.SetPartial("metric_dimensions")
	d.SetPartial("effective_interval")
	d.SetPartial("period")
	d.SetPartial("silence_time")
	d.SetPartial("webhook")
	d.SetPartial("escalations_critical")
	d.SetPartial("escalations_info")
	d.SetPartial("escalations_warn")
	d.SetPartial("prometheus")
	d.SetPartial("composite_expression")
	d.SetPartial("tags")
	d.SetPartial("dimensions")
	d.SetPartial("start_time")
	d.SetPartial("end_time")

	if d.Get("enabled").(bool) {
		action := "EnableMetricRules"
		enableMetricRequest := make(map[string]interface{})
		enableMetricRequest["RuleId"] = []string{d.Id()}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, enableMetricRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, enableMetricRequest)
	} else {
		action := "DisableMetricRules"
		disableMetricRequest := make(map[string]interface{})
		disableMetricRequest["RuleId"] = []string{d.Id()}

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, disableMetricRequest, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, disableMetricRequest)
	}

	if err := cmsService.WaitForCmsAlarm(d.Id(), d.Get("enabled").(bool), 102); err != nil {
		return WrapError(err)
	}

	d.SetPartial("enabled")

	update := false
	putMetricRuleTargetsReq := map[string]interface{}{
		"RuleId": d.Id(),
	}

	if d.HasChange("targets") {
		update = true
	}
	if v, ok := d.GetOk("targets"); ok {
		targetsMaps := make([]map[string]interface{}, 0)
		for _, targets := range v.(*schema.Set).List() {
			targetsMap := map[string]interface{}{}
			targetsArg := targets.(map[string]interface{})

			if id, ok := targetsArg["target_id"]; ok {
				targetsMap["Id"] = id
			}

			if jsonParams, ok := targetsArg["json_params"]; ok {
				targetsMap["JsonParams"] = jsonParams
			}

			if level, ok := targetsArg["level"]; ok {
				targetsMap["Level"] = level
			}

			if arn, ok := targetsArg["arn"]; ok {
				targetsMap["Arn"] = arn
			}

			targetsMaps = append(targetsMaps, targetsMap)
		}

		putMetricRuleTargetsReq["Targets"] = targetsMaps
	}

	if update {
		action := "PutMetricRuleTargets"

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, putMetricRuleTargetsReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, putMetricRuleTargetsReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("targets")
	}

	d.Partial(false)

	return resourceAliCloudCmsAlarmRead(d, meta)
}

func resourceAliCloudCmsAlarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	action := "DeleteMetricRules"
	var response map[string]interface{}

	var err error

	request := map[string]interface{}{
		"Id": []string{d.Id()},
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		_, err = cmsService.DescribeAlarm(d.Id())
		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe alarm rule got an error: %#v", err))
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	return nil
}

func convertCmsAlarmComparisonOperator(comparisonOperator string) string {
	switch comparisonOperator {
	case MoreThan:
		return "GreaterThanThreshold"
	case MoreThanOrEqual:
		return "GreaterThanOrEqualToThreshold"
	case LessThan:
		return "LessThanThreshold"
	case LessThanOrEqual:
		return "LessThanOrEqualToThreshold"
	case NotEqual:
		return "NotEqualToThreshold"
	case Equal:
		return "EqualToThreshold"
	case "GreaterThanThreshold":
		return MoreThan
	case "GreaterThanOrEqualToThreshold":
		return MoreThanOrEqual
	case "LessThanThreshold":
		return LessThan
	case "LessThanOrEqualToThreshold":
		return LessThanOrEqual
	case "NotEqualToThreshold":
		return NotEqual
	case "EqualToThreshold":
		return Equal
	default:
		return comparisonOperator
	}
}

func convertCmsAlarmPrometheusLevelResponse(source interface{}) interface{} {
	switch source {
	case "2":
		return "Critical"
	case "3":
		return "Warn"
	case "4":
		return "Info"
	}

	return source
}
