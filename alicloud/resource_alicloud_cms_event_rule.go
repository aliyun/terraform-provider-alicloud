package alicloud

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAliCloudCloudMonitorServiceEventRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceEventRuleCreate,
		Read:   resourceAliCloudCloudMonitorServiceEventRuleRead,
		Update: resourceAliCloudCloudMonitorServiceEventRuleUpdate,
		Delete: resourceAliCloudCloudMonitorServiceEventRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"silence_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: StringInSlice([]string{"DISABLED", "ENABLED"}, false),
			},
			"event_pattern": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"product": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sql_filter": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"level_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"event_type_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"contact_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact_parameters_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"contact_group_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"level": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"webhook_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"webhook_parameters_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"fc_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fc_parameters_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"function_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"sls_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sls_parameters_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"project": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"log_store": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"mns_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mns_parameters_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"queue": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"topic": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"open_api_parameters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"open_api_parameters_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"product": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"action": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"version": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"role": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceAliCloudCloudMonitorServiceEventRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "PutEventRule"
	request := make(map[string]interface{})
	var err error

	request["RuleName"] = d.Get("rule_name")

	if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	}

	if v, ok := d.GetOkExists("silence_time"); ok {
		request["SilenceTime"] = v
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
	}

	eventPatternMaps := make([]map[string]interface{}, 0)
	for _, eventPattern := range d.Get("event_pattern").(*schema.Set).List() {
		eventPatternMap := make(map[string]interface{})
		eventPatternArg := eventPattern.(map[string]interface{})

		eventPatternMap["Product"] = eventPatternArg["product"]

		if v, ok := eventPatternArg["sql_filter"].(string); ok && v != "" {
			eventPatternMap["SQLFilter"] = v
		}

		if v, ok := eventPatternArg["event_type_list"]; ok {
			eventPatternMap["EventTypeList"] = v
		}

		if v, ok := eventPatternArg["level_list"]; ok {
			eventPatternMap["LevelList"] = v
		}

		if v, ok := eventPatternArg["name_list"]; ok {
			eventPatternMap["NameList"] = v
		}

		eventPatternMaps = append(eventPatternMaps, eventPatternMap)
	}

	request["EventPattern"] = eventPatternMaps
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_event_rule", action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["RuleName"]))

	return resourceAliCloudCloudMonitorServiceEventRuleUpdate(d, meta)
}

func resourceAliCloudCloudMonitorServiceEventRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	object, err := cmsService.DescribeCmsEventRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("rule_name", object["Name"])
	d.Set("group_id", object["GroupId"])
	d.Set("silence_time", formatInt(object["SilenceTime"]))
	d.Set("description", object["Description"])
	d.Set("status", object["State"])

	if eventPattern, ok := object["EventPattern"]; ok {
		eventPatternMaps := make([]map[string]interface{}, 0)
		eventPatternArg := eventPattern.(map[string]interface{})
		eventPatternMap := make(map[string]interface{})

		if v, ok := eventPatternArg["Product"]; ok {
			eventPatternMap["product"] = v
		}

		if v, ok := eventPatternArg["SQLFilter"]; ok {
			eventPatternMap["sql_filter"] = v
		}

		if v, ok := eventPatternArg["NameList"]; ok {
			eventPatternMap["name_list"] = v.(map[string]interface{})["NameList"]
		}

		if v, ok := eventPatternArg["LevelList"]; ok {
			eventPatternMap["level_list"] = v.(map[string]interface{})["LevelList"]
		}

		if v, ok := eventPatternArg["EventTypeList"]; ok {
			eventPatternMap["event_type_list"] = v.(map[string]interface{})["EventTypeList"]
		}

		eventPatternMaps = append(eventPatternMaps, eventPatternMap)

		d.Set("event_pattern", eventPatternMaps)
	}

	targets, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceEventRuleTargets(d.Id())
	if err != nil {
		return WrapError(err)
	}

	if contactParameters, ok := targets["ContactParameters"]; ok {
		if contactParameterList, ok := contactParameters.(map[string]interface{})["ContactParameter"]; ok {
			contactParameterMaps := make([]map[string]interface{}, 0)
			for _, contactParameter := range contactParameterList.([]interface{}) {
				contactParameterArg := contactParameter.(map[string]interface{})
				contactParameterMap := map[string]interface{}{}

				if id, ok := contactParameterArg["Id"]; ok {
					contactParameterMap["contact_parameters_id"] = id
				}

				if contactGroupName, ok := contactParameterArg["ContactGroupName"]; ok {
					contactParameterMap["contact_group_name"] = contactGroupName
				}

				if level, ok := contactParameterArg["Level"]; ok {
					contactParameterMap["level"] = level
				}

				contactParameterMaps = append(contactParameterMaps, contactParameterMap)
			}

			d.Set("contact_parameters", contactParameterMaps)
		}
	}

	if webhookParameters, ok := targets["WebhookParameters"]; ok {
		if webhookParameterList, ok := webhookParameters.(map[string]interface{})["WebhookParameter"]; ok {
			webhookParameterMaps := make([]map[string]interface{}, 0)
			for _, webhookParameter := range webhookParameterList.([]interface{}) {
				webhookParameterArg := webhookParameter.(map[string]interface{})
				webhookParameterMap := map[string]interface{}{}

				if id, ok := webhookParameterArg["Id"]; ok {
					webhookParameterMap["webhook_parameters_id"] = id
				}

				if protocol, ok := webhookParameterArg["Protocol"]; ok {
					webhookParameterMap["protocol"] = protocol
				}

				if method, ok := webhookParameterArg["Method"]; ok {
					webhookParameterMap["method"] = method
				}

				if url, ok := webhookParameterArg["Url"]; ok {
					webhookParameterMap["url"] = url
				}

				webhookParameterMaps = append(webhookParameterMaps, webhookParameterMap)
			}

			d.Set("webhook_parameters", webhookParameterMaps)
		}
	}

	if fcParameters, ok := targets["FcParameters"]; ok {
		if fcParameterList, ok := fcParameters.(map[string]interface{})["FCParameter"]; ok {
			fcParameterMaps := make([]map[string]interface{}, 0)
			for _, fcParameter := range fcParameterList.([]interface{}) {
				fcParameterArg := fcParameter.(map[string]interface{})
				fcParameterMap := map[string]interface{}{}

				if id, ok := fcParameterArg["Id"]; ok {
					fcParameterMap["fc_parameters_id"] = id
				}

				if serviceName, ok := fcParameterArg["ServiceName"]; ok {
					fcParameterMap["service_name"] = serviceName
				}

				if functionName, ok := fcParameterArg["FunctionName"]; ok {
					fcParameterMap["function_name"] = functionName
				}

				if region, ok := fcParameterArg["Region"]; ok {
					fcParameterMap["region"] = region
				}

				if arn, ok := fcParameterArg["Arn"]; ok {
					fcParameterMap["arn"] = arn
				}

				fcParameterMaps = append(fcParameterMaps, fcParameterMap)
			}

			d.Set("fc_parameters", fcParameterMaps)
		}
	}

	if slsParameters, ok := targets["SlsParameters"]; ok {
		if slsParameterList, ok := slsParameters.(map[string]interface{})["SlsParameter"]; ok {
			slsParameterMaps := make([]map[string]interface{}, 0)
			for _, slsParameter := range slsParameterList.([]interface{}) {
				slsParameterArg := slsParameter.(map[string]interface{})
				slsParameterMap := map[string]interface{}{}

				if id, ok := slsParameterArg["Id"]; ok {
					slsParameterMap["sls_parameters_id"] = id
				}

				if project, ok := slsParameterArg["Project"]; ok {
					slsParameterMap["project"] = project
				}

				if logStore, ok := slsParameterArg["LogStore"]; ok {
					slsParameterMap["log_store"] = logStore
				}

				if region, ok := slsParameterArg["Region"]; ok {
					slsParameterMap["region"] = region
				}

				if arn, ok := slsParameterArg["Arn"]; ok {
					slsParameterMap["arn"] = arn
				}

				slsParameterMaps = append(slsParameterMaps, slsParameterMap)
			}

			d.Set("sls_parameters", slsParameterMaps)
		}
	}

	if mnsParameters, ok := targets["MnsParameters"]; ok {
		if mnsParameterList, ok := mnsParameters.(map[string]interface{})["MnsParameter"]; ok {
			mnsParameterMaps := make([]map[string]interface{}, 0)
			for _, mnsParameter := range mnsParameterList.([]interface{}) {
				mnsParameterArg := mnsParameter.(map[string]interface{})
				mnsParameterMap := map[string]interface{}{}

				if id, ok := mnsParameterArg["Id"]; ok {
					mnsParameterMap["mns_parameters_id"] = id
				}

				if queue, ok := mnsParameterArg["Queue"]; ok {
					mnsParameterMap["queue"] = queue
				}

				if topic, ok := mnsParameterArg["Topic"]; ok {
					mnsParameterMap["topic"] = topic
				}

				if region, ok := mnsParameterArg["Region"]; ok {
					mnsParameterMap["region"] = region
				}

				if arn, ok := mnsParameterArg["Arn"]; ok {
					mnsParameterMap["arn"] = arn
				}

				mnsParameterMaps = append(mnsParameterMaps, mnsParameterMap)
			}

			d.Set("mns_parameters", mnsParameterMaps)
		}
	}

	if openApiParameters, ok := targets["OpenApiParameters"]; ok {
		if openApiParameterList, ok := openApiParameters.(map[string]interface{})["OpenApiParameters"]; ok {
			openApiParameterMaps := make([]map[string]interface{}, 0)
			for _, openApiParameter := range openApiParameterList.([]interface{}) {
				openApiParameterArg := openApiParameter.(map[string]interface{})
				openApiParameterMap := map[string]interface{}{}

				if id, ok := openApiParameterArg["Id"]; ok {
					openApiParameterMap["open_api_parameters_id"] = id
				}

				if product, ok := openApiParameterArg["Product"]; ok {
					openApiParameterMap["product"] = product
				}

				if action, ok := openApiParameterArg["Action"]; ok {
					openApiParameterMap["action"] = action
				}

				if version, ok := openApiParameterArg["Version"]; ok {
					openApiParameterMap["version"] = version
				}

				if role, ok := openApiParameterArg["Role"]; ok {
					openApiParameterMap["role"] = role
				}

				if region, ok := openApiParameterArg["Region"]; ok {
					openApiParameterMap["region"] = region
				}

				if arn, ok := openApiParameterArg["Arn"]; ok {
					openApiParameterMap["arn"] = arn
				}

				openApiParameterMaps = append(openApiParameterMaps, openApiParameterMap)
			}

			d.Set("open_api_parameters", openApiParameterMaps)
		}
	}

	return nil
}

func resourceAliCloudCloudMonitorServiceEventRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	var err error
	d.Partial(true)

	update := false
	request := map[string]interface{}{
		"RuleName": d.Id(),
	}

	if !d.IsNewResource() && d.HasChange("group_id") {
		update = true
	}
	if v, ok := d.GetOk("group_id"); ok {
		request["GroupId"] = v
	}

	if !d.IsNewResource() && d.HasChange("silence_time") {
		update = true
	}
	if v, ok := d.GetOkExists("silence_time"); ok {
		request["SilenceTime"] = v
	}

	if !d.IsNewResource() && d.HasChange("description") {
		update = true
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if !d.IsNewResource() && d.HasChange("status") {
		update = true
	}
	if v, ok := d.GetOk("status"); ok {
		request["State"] = v
	}

	if !d.IsNewResource() && d.HasChange("event_pattern") {
		update = true
	}
	eventPatternMaps := make([]map[string]interface{}, 0)
	for _, eventPattern := range d.Get("event_pattern").(*schema.Set).List() {
		eventPatternMap := make(map[string]interface{})
		eventPatternArg := eventPattern.(map[string]interface{})

		eventPatternMap["Product"] = eventPatternArg["product"]

		if v, ok := eventPatternArg["sql_filter"].(string); ok && v != "" {
			eventPatternMap["SQLFilter"] = v
		}

		if v, ok := eventPatternArg["event_type_list"]; ok {
			eventPatternMap["EventTypeList"] = v
		}

		if v, ok := eventPatternArg["level_list"]; ok {
			eventPatternMap["LevelList"] = v
		}

		if v, ok := eventPatternArg["name_list"]; ok {
			eventPatternMap["NameList"] = v
		}

		eventPatternMaps = append(eventPatternMaps, eventPatternMap)
	}

	request["EventPattern"] = eventPatternMaps

	if update {
		action := "PutEventRule"
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

		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}

		d.SetPartial("group_id")
		d.SetPartial("silence_time")
		d.SetPartial("description")
		d.SetPartial("status")
		d.SetPartial("event_pattern")
	}

	update = false
	putEventRuleTargetsReq := map[string]interface{}{
		"RuleName": d.Id(),
	}

	if d.HasChange("contact_parameters") {
		update = true
	}
	if v, ok := d.GetOk("contact_parameters"); ok {
		contactParametersMaps := make([]map[string]interface{}, 0)
		for _, contactParameters := range v.([]interface{}) {
			contactParametersMap := map[string]interface{}{}
			contactParametersArg := contactParameters.(map[string]interface{})

			if id, ok := contactParametersArg["contact_parameters_id"]; ok {
				contactParametersMap["Id"] = id
			}

			if contactGroupName, ok := contactParametersArg["contact_group_name"]; ok {
				contactParametersMap["ContactGroupName"] = contactGroupName
			}

			if level, ok := contactParametersArg["level"]; ok {
				contactParametersMap["Level"] = level
			}

			contactParametersMaps = append(contactParametersMaps, contactParametersMap)
		}

		putEventRuleTargetsReq["ContactParameters"] = contactParametersMaps
	}

	if d.HasChange("webhook_parameters") {
		update = true
	}
	if v, ok := d.GetOk("webhook_parameters"); ok {
		webhookParametersMaps := make([]map[string]interface{}, 0)
		for _, webhookParameters := range v.([]interface{}) {
			webhookParametersMap := map[string]interface{}{}
			webhookParametersArg := webhookParameters.(map[string]interface{})

			if id, ok := webhookParametersArg["webhook_parameters_id"]; ok {
				webhookParametersMap["Id"] = id
			}

			if protocol, ok := webhookParametersArg["protocol"]; ok {
				webhookParametersMap["Protocol"] = protocol
			}

			if method, ok := webhookParametersArg["method"]; ok {
				webhookParametersMap["Method"] = method
			}

			if url, ok := webhookParametersArg["url"]; ok {
				webhookParametersMap["Url"] = url
			}

			webhookParametersMaps = append(webhookParametersMaps, webhookParametersMap)
		}

		putEventRuleTargetsReq["WebhookParameters"] = webhookParametersMaps
	}

	if d.HasChange("fc_parameters") {
		update = true
	}
	if v, ok := d.GetOk("fc_parameters"); ok {
		fcParametersMaps := make([]map[string]interface{}, 0)
		for _, fcParameters := range v.([]interface{}) {
			fcParametersMap := map[string]interface{}{}
			fcParametersArg := fcParameters.(map[string]interface{})

			if id, ok := fcParametersArg["fc_parameters_id"]; ok {
				fcParametersMap["Id"] = id
			}

			if serviceName, ok := fcParametersArg["service_name"]; ok {
				fcParametersMap["ServiceName"] = serviceName
			}

			if functionName, ok := fcParametersArg["function_name"]; ok {
				fcParametersMap["FunctionName"] = functionName
			}

			if region, ok := fcParametersArg["region"]; ok {
				fcParametersMap["Region"] = region
			}

			fcParametersMaps = append(fcParametersMaps, fcParametersMap)
		}

		putEventRuleTargetsReq["FcParameters"] = fcParametersMaps
	}

	if d.HasChange("sls_parameters") {
		update = true
	}
	if v, ok := d.GetOk("sls_parameters"); ok {
		slsParametersMaps := make([]map[string]interface{}, 0)
		for _, slsParameters := range v.([]interface{}) {
			slsParametersMap := map[string]interface{}{}
			slsParametersArg := slsParameters.(map[string]interface{})

			if id, ok := slsParametersArg["sls_parameters_id"]; ok {
				slsParametersMap["Id"] = id
			}

			if project, ok := slsParametersArg["project"]; ok {
				slsParametersMap["Project"] = project
			}

			if logStore, ok := slsParametersArg["log_store"]; ok {
				slsParametersMap["LogStore"] = logStore
			}

			if region, ok := slsParametersArg["region"]; ok {
				slsParametersMap["Region"] = region
			}

			slsParametersMaps = append(slsParametersMaps, slsParametersMap)
		}

		putEventRuleTargetsReq["SlsParameters"] = slsParametersMaps
	}

	if d.HasChange("mns_parameters") {
		update = true
	}
	if v, ok := d.GetOk("mns_parameters"); ok {
		mnsParametersMaps := make([]map[string]interface{}, 0)
		for _, mnsParameters := range v.([]interface{}) {
			mnsParametersMap := map[string]interface{}{}
			mnsParametersArg := mnsParameters.(map[string]interface{})

			if id, ok := mnsParametersArg["mns_parameters_id"]; ok {
				mnsParametersMap["Id"] = id
			}

			if queue, ok := mnsParametersArg["queue"]; ok {
				mnsParametersMap["Queue"] = queue
			}

			if topic, ok := mnsParametersArg["topic"]; ok {
				mnsParametersMap["Topic"] = topic
			}

			if region, ok := mnsParametersArg["region"]; ok {
				mnsParametersMap["Region"] = region
			}

			mnsParametersMaps = append(mnsParametersMaps, mnsParametersMap)
		}

		putEventRuleTargetsReq["MnsParameters"] = mnsParametersMaps
	}

	if d.HasChange("open_api_parameters") {
		update = true
	}
	if v, ok := d.GetOk("open_api_parameters"); ok {
		openApiParametersMaps := make([]map[string]interface{}, 0)
		for _, openApiParameters := range v.([]interface{}) {
			openApiParametersMap := map[string]interface{}{}
			openApiParametersArg := openApiParameters.(map[string]interface{})

			if id, ok := openApiParametersArg["open_api_parameters_id"]; ok {
				openApiParametersMap["Id"] = id
			}

			if product, ok := openApiParametersArg["product"]; ok {
				openApiParametersMap["Product"] = product
			}

			if action, ok := openApiParametersArg["action"]; ok {
				openApiParametersMap["Action"] = action
			}

			if version, ok := openApiParametersArg["version"]; ok {
				openApiParametersMap["Version"] = version
			}

			if role, ok := openApiParametersArg["role"]; ok {
				openApiParametersMap["Role"] = role
			}

			if region, ok := openApiParametersArg["region"]; ok {
				openApiParametersMap["Region"] = region
			}

			openApiParametersMaps = append(openApiParametersMaps, openApiParametersMap)
		}

		putEventRuleTargetsReq["OpenApiParameters"] = openApiParametersMaps
	}

	if update {
		action := "PutEventRuleTargets"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *resource.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, putEventRuleTargetsReq, false)
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, putEventRuleTargetsReq)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}

		d.SetPartial("resource_group_id")
	}

	d.Partial(false)

	return resourceAliCloudCloudMonitorServiceEventRuleRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceEventRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteEventRules"
	var response map[string]interface{}
	var err error
	request := map[string]interface{}{
		"RuleNames": []string{d.Id()},
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
