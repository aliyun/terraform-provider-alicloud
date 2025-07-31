// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"time"
)

func dataSourceAliCloudCloudMonitorServiceMetricAlarmRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAliCloudCloudMonitorServiceMetricAlarmRuleRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dimensions": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"composite_expression": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"expression_raw": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expression_list_join": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"level": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"expression_list": {
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
								},
							},
						},
						"contact_groups": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effective_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email_subject": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"escalations": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"critical": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
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
									"info": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
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
									"warn": {
										Type:     schema.TypeList,
										Computed: true,
										MaxItems: 1,
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
								},
							},
						},
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"no_data_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"no_effective_interval": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prometheus": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"annotations": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"prom_ql": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"times": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"level": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"resources": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"silence_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"webhook": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceAliCloudCloudMonitorServiceMetricAlarmRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	action := "DescribeMetricRuleList"
	var err error
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	if v, ok := d.GetOk("dimensions"); ok {
		request["Dimensions"] = v
	}

	if v, ok := d.GetOk("metric_name"); ok {
		request["MetricName"] = v
	}

	if v, ok := d.GetOk("namespace"); ok {
		request["Namespace"] = v
	}

	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}

	if v, ok := d.GetOkExists("status"); ok {
		request["EnableState"] = v
	}

	request["PageSize"] = PageSizeLarge
	request["Page"] = 1
	for {
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

		resp, err := jsonpath.Get("$.Alarms.Alarm", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Alarms.Alarm", response)
		}

		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RuleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}

		if len(result) < request["PageSize"].(int) {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}

	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, objectRaw := range objects {
		mapping := map[string]interface{}{}

		mapping["id"] = objectRaw["RuleId"]

		mapping["contact_groups"] = objectRaw["ContactGroups"]
		mapping["dimensions"] = objectRaw["Dimensions"]
		mapping["effective_interval"] = objectRaw["EffectiveInterval"]
		mapping["email_subject"] = objectRaw["MailSubject"]
		mapping["metric_name"] = objectRaw["MetricName"]
		mapping["namespace"] = objectRaw["Namespace"]
		mapping["no_data_policy"] = objectRaw["NoDataPolicy"]
		mapping["no_effective_interval"] = objectRaw["NoEffectiveInterval"]
		mapping["period"] = objectRaw["Period"]
		mapping["resources"] = objectRaw["Resources"]
		mapping["rule_name"] = objectRaw["RuleName"]
		mapping["silence_time"] = objectRaw["SilenceTime"]
		mapping["source_type"] = objectRaw["SourceType"]
		mapping["status"] = objectRaw["EnableState"]
		mapping["webhook"] = objectRaw["Webhook"]

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
				for _, expressionListChildRaw := range expressionListRaw.([]interface{}) {
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
		mapping["composite_expression"] = compositeExpressionMaps
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
			if len(criticalRaw) > 0 {
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
			if len(infoRaw) > 0 {
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
			if len(warnRaw) > 0 {
				warnMap["comparison_operator"] = warnRaw["ComparisonOperator"]
				warnMap["pre_condition"] = warnRaw["PreCondition"]
				warnMap["statistics"] = warnRaw["Statistics"]
				warnMap["threshold"] = warnRaw["Threshold"]
				warnMap["times"] = warnRaw["Times"]

				warnMaps = append(warnMaps, warnMap)
			}
			escalationsMap["warn"] = warnMaps
			escalationsMaps = append(escalationsMaps, escalationsMap)
		}
		mapping["escalations"] = escalationsMaps
		labelsRaw, _ := jsonpath.Get("$.Labels.Labels", objectRaw)
		labelsMaps := make([]map[string]interface{}, 0)
		if labelsRaw != nil {
			for _, labelsChildRaw := range labelsRaw.([]interface{}) {
				labelsMap := make(map[string]interface{})
				labelsChildRaw := labelsChildRaw.(map[string]interface{})
				labelsMap["key"] = labelsChildRaw["Key"]
				labelsMap["value"] = labelsChildRaw["Value"]

				labelsMaps = append(labelsMaps, labelsMap)
			}
		}
		mapping["labels"] = labelsMaps
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
				for _, annotationsChildRaw := range annotationsRaw.([]interface{}) {
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
		mapping["prometheus"] = prometheusMaps

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
