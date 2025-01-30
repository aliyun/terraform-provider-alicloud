package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsMetricRuleTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsMetricRuleTemplateCreate,
		Read:   resourceAliCloudCmsMetricRuleTemplateRead,
		Update: resourceAliCloudCmsMetricRuleTemplateUpdate,
		Delete: resourceAliCloudCmsMetricRuleTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"metric_rule_template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"apply_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"silence_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: IntBetween(0, 86400),
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alert_templates": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"webhook": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"escalations": {
							Type:     schema.TypeSet,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"critical": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"info": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"warn": {
										Type:     schema.TypeSet,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"comparison_operator": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
												},
												"statistics": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"threshold": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"times": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"rest_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliCloudCmsMetricRuleTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMetricRuleTemplate"
	request := make(map[string]interface{})
	var err error

	request["Name"] = d.Get("metric_rule_template_name")

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("alert_templates"); ok {
		alertTemplatesMaps := make([]map[string]interface{}, 0)
		for _, alertTemplates := range v.(*schema.Set).List() {
			alertTemplatesMap := map[string]interface{}{}
			alertTemplatesArg := alertTemplates.(map[string]interface{})

			alertTemplatesMap["RuleName"] = alertTemplatesArg["rule_name"]
			alertTemplatesMap["MetricName"] = alertTemplatesArg["metric_name"]
			alertTemplatesMap["Namespace"] = alertTemplatesArg["namespace"]
			alertTemplatesMap["Category"] = alertTemplatesArg["category"]
			alertTemplatesMap["Webhook"] = alertTemplatesArg["webhook"]

			if escalations, ok := alertTemplatesArg["escalations"]; ok {
				escalationsMap := map[string]interface{}{}
				for _, escalationsList := range escalations.(*schema.Set).List() {
					escalationsArg := escalationsList.(map[string]interface{})

					if critical, ok := escalationsArg["critical"]; ok {
						criticalMap := map[string]interface{}{}
						for _, criticalList := range critical.(*schema.Set).List() {
							criticalArg := criticalList.(map[string]interface{})

							if comparisonOperator, ok := criticalArg["comparison_operator"]; ok {
								criticalMap["ComparisonOperator"] = comparisonOperator
							}

							if statistics, ok := criticalArg["statistics"]; ok {
								criticalMap["Statistics"] = statistics
							}

							if threshold, ok := criticalArg["threshold"]; ok {
								criticalMap["Threshold"] = threshold
							}

							if times, ok := criticalArg["times"]; ok {
								criticalMap["Times"] = times
							}
						}

						escalationsMap["Critical"] = criticalMap
					}

					if info, ok := escalationsArg["info"]; ok {
						infoMap := map[string]interface{}{}
						for _, infoList := range info.(*schema.Set).List() {
							infoArg := infoList.(map[string]interface{})

							if comparisonOperator, ok := infoArg["comparison_operator"]; ok {
								infoMap["ComparisonOperator"] = comparisonOperator
							}

							if statistics, ok := infoArg["statistics"]; ok {
								infoMap["Statistics"] = statistics
							}

							if threshold, ok := infoArg["threshold"]; ok {
								infoMap["Threshold"] = threshold
							}

							if times, ok := infoArg["times"]; ok {
								infoMap["Times"] = times
							}
						}

						escalationsMap["Info"] = infoMap
					}

					if warn, ok := escalationsArg["warn"]; ok {
						warnMap := map[string]interface{}{}
						for _, warnList := range warn.(*schema.Set).List() {
							warnArg := warnList.(map[string]interface{})

							if comparisonOperator, ok := warnArg["comparison_operator"]; ok {
								warnMap["ComparisonOperator"] = comparisonOperator
							}

							if statistics, ok := warnArg["statistics"]; ok {
								warnMap["Statistics"] = statistics
							}

							if threshold, ok := warnArg["threshold"]; ok {
								warnMap["Threshold"] = threshold
							}

							if times, ok := warnArg["times"]; ok {
								warnMap["Times"] = times
							}
						}

						escalationsMap["Warn"] = warnMap
					}
				}

				alertTemplatesMap["Escalations"] = escalationsMap
			}

			alertTemplatesMaps = append(alertTemplatesMaps, alertTemplatesMap)
		}

		request["AlertTemplates"] = alertTemplatesMaps
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_metric_rule_template", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(response["Id"]))

	return resourceAliCloudCmsMetricRuleTemplateUpdate(d, meta)
}

func resourceAliCloudCmsMetricRuleTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	object, err := cmsService.DescribeCmsMetricRuleTemplate(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_metric_rule_template cmsService.DescribeCmsMetricRuleTemplate Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("metric_rule_template_name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("rest_version", object["RestVersion"])

	if alertTemplates, ok := object["AlertTemplates"]; ok {
		if alertTemplateList, ok := alertTemplates.(map[string]interface{})["AlertTemplate"]; ok {
			alertTemplateMaps := make([]map[string]interface{}, 0)
			for _, alertTemplate := range alertTemplateList.([]interface{}) {
				alertTemplateArg := alertTemplate.(map[string]interface{})
				alertTemplateMap := map[string]interface{}{}

				if ruleName, ok := alertTemplateArg["RuleName"]; ok {
					alertTemplateMap["rule_name"] = ruleName
				}

				if metricName, ok := alertTemplateArg["MetricName"]; ok {
					alertTemplateMap["metric_name"] = metricName
				}

				if namespace, ok := alertTemplateArg["Namespace"]; ok {
					alertTemplateMap["namespace"] = namespace
				}

				if category, ok := alertTemplateArg["Category"]; ok {
					alertTemplateMap["category"] = category
				}

				if webhook, ok := alertTemplateArg["Webhook"]; ok {
					alertTemplateMap["webhook"] = webhook
				}

				if escalations, ok := alertTemplateArg["Escalations"]; ok {
					escalationsMaps := make([]map[string]interface{}, 0)
					escalationsArg := escalations.(map[string]interface{})
					escalationsMap := make(map[string]interface{})

					if critical, ok := escalationsArg["Critical"]; ok && len(critical.(map[string]interface{})) > 0 {
						criticalMaps := make([]map[string]interface{}, 0)
						criticalArg := critical.(map[string]interface{})
						criticalMap := map[string]interface{}{}

						if comparisonOperator, ok := criticalArg["ComparisonOperator"]; ok {
							criticalMap["comparison_operator"] = comparisonOperator
						}

						if statistics, ok := criticalArg["Statistics"]; ok {
							criticalMap["statistics"] = statistics
						}

						if threshold, ok := criticalArg["Threshold"]; ok {
							criticalMap["threshold"] = threshold
						}

						if times, ok := criticalArg["Times"]; ok {
							criticalMap["times"] = times
						}

						criticalMaps = append(criticalMaps, criticalMap)

						escalationsMap["critical"] = criticalMaps
					}

					if info, ok := escalationsArg["Info"]; ok && len(info.(map[string]interface{})) > 0 {
						infoMaps := make([]map[string]interface{}, 0)
						infoArg := info.(map[string]interface{})
						infoMap := map[string]interface{}{}

						if comparisonOperator, ok := infoArg["ComparisonOperator"]; ok {
							infoMap["comparison_operator"] = comparisonOperator
						}

						if statistics, ok := infoArg["Statistics"]; ok {
							infoMap["statistics"] = statistics
						}

						if threshold, ok := infoArg["Threshold"]; ok {
							infoMap["threshold"] = threshold
						}

						if times, ok := infoArg["Times"]; ok {
							infoMap["times"] = times
						}

						infoMaps = append(infoMaps, infoMap)

						escalationsMap["info"] = infoMaps
					}

					if warn, ok := escalationsArg["Warn"]; ok && len(warn.(map[string]interface{})) > 0 {
						warnMaps := make([]map[string]interface{}, 0)
						warnArg := warn.(map[string]interface{})
						warnMap := map[string]interface{}{}

						if comparisonOperator, ok := warnArg["ComparisonOperator"]; ok {
							warnMap["comparison_operator"] = comparisonOperator
						}

						if statistics, ok := warnArg["Statistics"]; ok {
							warnMap["statistics"] = statistics
						}

						if threshold, ok := warnArg["Threshold"]; ok {
							warnMap["threshold"] = threshold
						}

						if times, ok := warnArg["Times"]; ok {
							warnMap["times"] = times
						}

						warnMaps = append(warnMaps, warnMap)

						escalationsMap["warn"] = warnMaps
					}

					escalationsMaps = append(escalationsMaps, escalationsMap)

					alertTemplateMap["escalations"] = escalationsMaps
				}

				alertTemplateMaps = append(alertTemplateMaps, alertTemplateMap)
			}

			d.Set("alert_templates", alertTemplateMaps)
		}
	}

	return nil
}

func resourceAliCloudCmsMetricRuleTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	modifyMetricRuleTemplateReq := map[string]interface{}{
		"TemplateId":  d.Id(),
		"RestVersion": d.Get("rest_version"),
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		modifyMetricRuleTemplateReq["Description"] = v
	}

	if !d.IsNewResource() && d.HasChange("alert_templates") {
		update = true
	}
	if v, ok := d.GetOk("alert_templates"); ok {
		alertTemplatesMaps := make([]map[string]interface{}, 0)
		for _, alertTemplates := range v.(*schema.Set).List() {
			alertTemplatesMap := map[string]interface{}{}
			alertTemplatesArg := alertTemplates.(map[string]interface{})

			alertTemplatesMap["RuleName"] = alertTemplatesArg["rule_name"]
			alertTemplatesMap["MetricName"] = alertTemplatesArg["metric_name"]
			alertTemplatesMap["Namespace"] = alertTemplatesArg["namespace"]
			alertTemplatesMap["Category"] = alertTemplatesArg["category"]
			alertTemplatesMap["Webhook"] = alertTemplatesArg["webhook"]

			if escalations, ok := alertTemplatesArg["escalations"]; ok {
				escalationsMap := map[string]interface{}{}
				for _, escalationsList := range escalations.(*schema.Set).List() {
					escalationsArg := escalationsList.(map[string]interface{})

					if critical, ok := escalationsArg["critical"]; ok {
						criticalMap := map[string]interface{}{}
						for _, criticalList := range critical.(*schema.Set).List() {
							criticalArg := criticalList.(map[string]interface{})

							if comparisonOperator, ok := criticalArg["comparison_operator"]; ok {
								criticalMap["ComparisonOperator"] = comparisonOperator
							}

							if statistics, ok := criticalArg["statistics"]; ok {
								criticalMap["Statistics"] = statistics
							}

							if threshold, ok := criticalArg["threshold"]; ok {
								criticalMap["Threshold"] = threshold
							}

							if times, ok := criticalArg["times"]; ok {
								criticalMap["Times"] = times
							}
						}

						escalationsMap["Critical"] = criticalMap
					}

					if info, ok := escalationsArg["info"]; ok {
						infoMap := map[string]interface{}{}
						for _, infoList := range info.(*schema.Set).List() {
							infoArg := infoList.(map[string]interface{})

							if comparisonOperator, ok := infoArg["comparison_operator"]; ok {
								infoMap["ComparisonOperator"] = comparisonOperator
							}

							if statistics, ok := infoArg["statistics"]; ok {
								infoMap["Statistics"] = statistics
							}

							if threshold, ok := infoArg["threshold"]; ok {
								infoMap["Threshold"] = threshold
							}

							if times, ok := infoArg["times"]; ok {
								infoMap["Times"] = times
							}
						}

						escalationsMap["Info"] = infoMap
					}

					if warn, ok := escalationsArg["warn"]; ok {
						warnMap := map[string]interface{}{}
						for _, warnList := range warn.(*schema.Set).List() {
							warnArg := warnList.(map[string]interface{})

							if comparisonOperator, ok := warnArg["comparison_operator"]; ok {
								warnMap["ComparisonOperator"] = comparisonOperator
							}

							if statistics, ok := warnArg["statistics"]; ok {
								warnMap["Statistics"] = statistics
							}

							if threshold, ok := warnArg["threshold"]; ok {
								warnMap["Threshold"] = threshold
							}

							if times, ok := warnArg["times"]; ok {
								warnMap["Times"] = times
							}
						}

						escalationsMap["Warn"] = warnMap
					}
				}

				alertTemplatesMap["Escalations"] = escalationsMap
			}

			alertTemplatesMaps = append(alertTemplatesMaps, alertTemplatesMap)
		}

		modifyMetricRuleTemplateReq["AlertTemplates"] = alertTemplatesMaps
	}

	if update {
		action := "ModifyMetricRuleTemplate"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, modifyMetricRuleTemplateReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, modifyMetricRuleTemplateReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("metric_rule_template_name")
		d.SetPartial("description")
		d.SetPartial("alert_templates")
	}

	update = false
	applyMetricRuleTemplateReq := map[string]interface{}{
		"TemplateIds": d.Id(),
	}

	if d.HasChange("group_id") {
		update = true
	}
	if v, ok := d.GetOk("group_id"); ok {
		applyMetricRuleTemplateReq["GroupId"] = v
	}

	if d.HasChange("apply_mode") {
		update = true
	}
	if v, ok := d.GetOk("apply_mode"); ok {
		applyMetricRuleTemplateReq["ApplyMode"] = v
	}

	if d.HasChange("notify_level") {
		update = true
	}
	if v, ok := d.GetOk("notify_level"); ok {
		applyMetricRuleTemplateReq["NotifyLevel"] = v
	}

	if d.HasChange("silence_time") {
		update = true
	}
	if v, ok := d.GetOkExists("silence_time"); ok {
		applyMetricRuleTemplateReq["SilenceTime"] = v
	}

	if d.HasChange("webhook") {
		update = true
	}
	if v, ok := d.GetOk("webhook"); ok {
		applyMetricRuleTemplateReq["Webhook"] = v
	}

	if d.HasChange("enable_start_time") {
		update = true
	}
	if v, ok := d.GetOk("enable_start_time"); ok {
		applyMetricRuleTemplateReq["EnableStartTime"] = v
	}

	if d.HasChange("enable_end_time") {
		update = true
	}
	if v, ok := d.GetOk("enable_end_time"); ok {
		applyMetricRuleTemplateReq["EnableEndTime"] = v
	}

	if update {
		action := "ApplyMetricRuleTemplate"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, applyMetricRuleTemplateReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, applyMetricRuleTemplateReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("group_id")
		d.SetPartial("apply_mode")
		d.SetPartial("notify_level")
		d.SetPartial("silence_time")
		d.SetPartial("webhook")
		d.SetPartial("enable_start_time")
		d.SetPartial("enable_end_time")
	}

	d.Partial(false)

	return resourceAliCloudCmsMetricRuleTemplateRead(d, meta)
}

func resourceAliCloudCmsMetricRuleTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMetricRuleTemplate"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"TemplateId": d.Id(),
	}

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
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		if IsExpectedErrorCodes(fmt.Sprint(response["Code"]), []string{"ResourceNotFound"}) || NotFoundError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
