package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCmsGroupMetricRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCmsGroupMetricRuleCreate,
		Read:   resourceAliCloudCmsGroupMetricRuleRead,
		Update: resourceAliCloudCmsGroupMetricRuleUpdate,
		Delete: resourceAliCloudCmsGroupMetricRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_metric_rule_name": {
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
				ForceNew: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_groups": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dimensions": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"email_subject": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"no_effective_interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"interval": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"silence_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"webhook": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"targets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
			"escalations": {
				Type:     schema.TypeSet,
				Required: true,
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
										Type:     schema.TypeString,
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
									"times": {
										Type:     schema.TypeInt,
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
										Type:     schema.TypeString,
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
									"times": {
										Type:     schema.TypeInt,
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
										Type:     schema.TypeString,
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
									"times": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudCmsGroupMetricRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutGroupMetricRule"
	request := make(map[string]interface{})
	var err error

	request["RuleId"] = d.Get("rule_id")
	request["GroupId"] = d.Get("group_id")
	request["RuleName"] = d.Get("group_metric_rule_name")
	request["MetricName"] = d.Get("metric_name")
	request["Namespace"] = d.Get("namespace")

	if v, ok := d.GetOk("category"); ok {
		request["Category"] = v
	}

	if v, ok := d.GetOk("contact_groups"); ok {
		request["ContactGroups"] = v
	}

	if v, ok := d.GetOk("dimensions"); ok {
		request["ExtraDimensionJson"] = v
	}

	if v, ok := d.GetOk("email_subject"); ok {
		request["EmailSubject"] = v
	}

	if v, ok := d.GetOk("effective_interval"); ok {
		request["EffectiveInterval"] = v
	}

	if v, ok := d.GetOk("no_effective_interval"); ok {
		request["NoEffectiveInterval"] = v
	}

	if v, ok := d.GetOk("interval"); ok {
		request["Interval"] = v
	}

	if v, ok := d.GetOk("period"); ok {
		request["Period"] = v
	}

	if v, ok := d.GetOk("silence_time"); ok {
		request["SilenceTime"] = v
	}

	if v, ok := d.GetOk("webhook"); ok {
		request["Webhook"] = v
	}

	escalationsMap := map[string]interface{}{}
	for _, escalationsList := range d.Get("escalations").(*schema.Set).List() {
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

	request["Escalations"] = escalationsMap

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ExceedingQuota"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_group_metric_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["RuleId"]))

	return resourceAliCloudCmsGroupMetricRuleUpdate(d, meta)
}

func resourceAliCloudCmsGroupMetricRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}

	object, err := cmsService.DescribeCmsGroupMetricRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_group_metric_rule cmsService.DescribeCmsGroupMetricRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule_id", object["RuleId"])
	d.Set("group_id", object["GroupId"])
	d.Set("group_metric_rule_name", object["RuleName"])
	d.Set("metric_name", object["MetricName"])
	d.Set("namespace", object["Namespace"])
	d.Set("contact_groups", object["ContactGroups"])
	d.Set("dimensions", removeSquareBracketsFromDimensions(fmt.Sprint(object["Dimensions"])))
	d.Set("email_subject", object["MailSubject"])
	d.Set("effective_interval", object["EffectiveInterval"])
	d.Set("no_effective_interval", object["NoEffectiveInterval"])
	d.Set("period", formatInt(object["Period"]))
	d.Set("silence_time", formatInt(object["SilenceTime"]))
	d.Set("webhook", object["Webhook"])
	d.Set("status", object["AlertState"])

	if escalations, ok := object["Escalations"]; ok {
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

		d.Set("escalations", escalationsMaps)
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
			targetsMap["id"] = id
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

func resourceAliCloudCmsGroupMetricRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RuleId":    d.Id(),
		"Namespace": d.Get("namespace"),
	}

	if !d.IsNewResource() && d.HasChange("group_id") {
		update = true
	}
	request["GroupId"] = d.Get("group_id")

	if !d.IsNewResource() && d.HasChange("group_metric_rule_name") {
		update = true
	}
	request["RuleName"] = d.Get("group_metric_rule_name")

	if !d.IsNewResource() && d.HasChange("metric_name") {
		update = true
	}
	request["MetricName"] = d.Get("metric_name")

	if !d.IsNewResource() && d.HasChange("contact_groups") {
		update = true
	}
	if v, ok := d.GetOk("contact_groups"); ok {
		request["ContactGroups"] = v
	}

	if !d.IsNewResource() && d.HasChange("dimensions") {
		update = true
	}
	if v, ok := d.GetOk("dimensions"); ok {
		request["ExtraDimensionJson"] = v
	}

	if !d.IsNewResource() && d.HasChange("email_subject") {
		update = true
	}
	if v, ok := d.GetOk("email_subject"); ok {
		request["EmailSubject"] = v
	}

	if !d.IsNewResource() && d.HasChange("effective_interval") {
		update = true
	}
	if v, ok := d.GetOk("effective_interval"); ok {
		request["EffectiveInterval"] = v
	}

	if !d.IsNewResource() && d.HasChange("no_effective_interval") {
		update = true
	}
	if v, ok := d.GetOk("no_effective_interval"); ok {
		request["NoEffectiveInterval"] = v
	}

	if !d.IsNewResource() && d.HasChange("period") {
		update = true

		if v, ok := d.GetOkExists("period"); ok {
			request["Period"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("silence_time") {
		update = true

		if v, ok := d.GetOkExists("silence_time"); ok {
			request["SilenceTime"] = v
		}
	}

	if !d.IsNewResource() && d.HasChange("webhook") {
		update = true
	}
	if v, ok := d.GetOk("webhook"); ok {
		request["Webhook"] = v
	}

	if !d.IsNewResource() && d.HasChange("escalations") {
		update = true
	}
	escalationsMap := map[string]interface{}{}
	for _, escalationsList := range d.Get("escalations").(*schema.Set).List() {
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

	request["Escalations"] = escalationsMap

	if update {
		action := "PutGroupMetricRule"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if IsExpectedErrors(err, []string{"ExceedingQuota"}) || NeedRetry(err) {
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

		d.SetPartial("group_id")
		d.SetPartial("group_metric_rule_name")
		d.SetPartial("metric_name")
		d.SetPartial("contact_groups")
		d.SetPartial("dimensions")
		d.SetPartial("email_subject")
		d.SetPartial("effective_interval")
		d.SetPartial("no_effective_interval")
		d.SetPartial("period")
		d.SetPartial("silence_time")
		d.SetPartial("webhook")
		d.SetPartial("escalations")
	}

	update = false
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

			if id, ok := targetsArg["id"]; ok {
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

		d.SetPartial("targets")
	}

	d.Partial(false)

	return resourceAliCloudCmsGroupMetricRuleRead(d, meta)
}

func resourceAliCloudCmsGroupMetricRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMetricRules"
	var response map[string]interface{}
	var err error

	request := map[string]interface{}{
		"Id": []string{d.Id()},
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *resource.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if IsExpectedErrors(err, []string{"ExceedingQuota"}) || NeedRetry(err) {
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

	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"400", "403", "404", "ResourceNotFound"}) {
		return nil
	}

	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("DeleteMetricRules failed for " + response["Message"].(string)))
	}

	return nil
}

func removeSquareBracketsFromDimensions(src string) string {
	if len(src) < 1 && src[0] != '[' && src[len(src)-1] != ']' {
		return src
	}

	return src[1 : len(src)-1]
}
