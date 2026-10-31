package cloudmonitorservice

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/validators"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcess() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessCreate,
		Read:   resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessRead,
		Update: resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessUpdate,
		Delete: resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"process_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"match_express_filter_relation": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validators.StringInSlice([]string{"all", "and", "or"}, false),
			},
			"group_monitoring_agent_process_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"match_express": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"function": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Computed:     true,
							ValidateFunc: validators.StringInSlice([]string{"all", "startWith", "endWith", "contains", "notContains", "equals"}, false),
						},
					},
				},
			},
			"alert_config": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"escalations_level": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validators.StringInSlice([]string{"critical", "warn", "info"}, false),
						},
						"comparison_operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validators.StringInSlice([]string{"GreaterThanOrEqualToThreshold", "GreaterThanThreshold", "LessThanOrEqualToThreshold", "LessThanThreshold", "NotEqualToThreshold", "GreaterThanYesterday", "LessThanYesterday", "GreaterThanLastWeek", "LessThanLastWeek", "GreaterThanLastPeriod", "LessThanLastPeriod"}, false),
						},
						"statistics": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validators.StringInSlice([]string{"Average"}, false),
						},
						"threshold": {
							Type:     schema.TypeString,
							Required: true,
						},
						"times": {
							Type:     schema.TypeString,
							Required: true,
						},
						"effective_interval": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"silence_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validators.IntBetween(3600, 86400),
						},
						"webhook": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_list_id": {
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
										ValidateFunc: validators.StringInSlice([]string{"CRITICAL", "WARN", "INFO"}, false),
									},
									"arn": {
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
	}
}

func resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateGroupMonitoringAgentProcess"
	request := make(map[string]interface{})
	var err error

	request["GroupId"] = d.Get("group_id")
	request["ProcessName"] = d.Get("process_name")

	if v, ok := d.GetOk("match_express_filter_relation"); ok {
		request["MatchExpressFilterRelation"] = v
	}

	if v, ok := d.GetOk("match_express"); ok {
		matchExpressMaps := make([]map[string]interface{}, 0)
		for _, matchExpress := range v.([]interface{}) {
			matchExpressMap := map[string]interface{}{}
			matchExpressArg := matchExpress.(map[string]interface{})

			if name, ok := matchExpressArg["name"]; ok {
				matchExpressMap["Name"] = name
			}

			if value, ok := matchExpressArg["value"]; ok {
				matchExpressMap["Value"] = value
			}

			if function, ok := matchExpressArg["function"]; ok {
				matchExpressMap["Function"] = function
			}

			matchExpressMaps = append(matchExpressMaps, matchExpressMap)
		}

		request["MatchExpress"] = matchExpressMaps
	}

	alertConfig := d.Get("alert_config")
	alertConfigMaps := make([]map[string]interface{}, 0)
	for _, alertConfigList := range alertConfig.([]interface{}) {
		alertConfigMap := map[string]interface{}{}
		alertConfigArg := alertConfigList.(map[string]interface{})

		alertConfigMap["EscalationsLevel"] = alertConfigArg["escalations_level"]
		alertConfigMap["ComparisonOperator"] = alertConfigArg["comparison_operator"]
		alertConfigMap["Statistics"] = alertConfigArg["statistics"]
		alertConfigMap["Threshold"] = alertConfigArg["threshold"]
		alertConfigMap["Times"] = alertConfigArg["times"]

		if effectiveInterval, ok := alertConfigArg["effective_interval"]; ok {
			alertConfigMap["EffectiveInterval"] = effectiveInterval
		}

		if silenceTime, ok := alertConfigArg["silence_time"]; ok {
			alertConfigMap["SilenceTime"] = silenceTime
		}

		if webhook, ok := alertConfigArg["webhook"]; ok {
			alertConfigMap["Webhook"] = webhook
		}

		if targetList, ok := alertConfigArg["target_list"]; ok {
			targetListMaps := make([]map[string]interface{}, 0)
			for _, targetListArgList := range targetList.([]interface{}) {
				targetListMap := map[string]interface{}{}
				targetListArg := targetListArgList.(map[string]interface{})

				if id, ok := targetListArg["target_list_id"]; ok {
					targetListMap["Id"] = id
				}

				if jsonParams, ok := targetListArg["json_params"]; ok {
					targetListMap["JsonParams"] = jsonParams
				}

				if level, ok := targetListArg["level"]; ok {
					targetListMap["Level"] = level
				}

				if arn, ok := targetListArg["arn"]; ok {
					targetListMap["Arn"] = arn
				}

				targetListMaps = append(targetListMaps, targetListMap)
			}

			alertConfigMap["TargetList"] = targetListMaps
		}

		alertConfigMaps = append(alertConfigMaps, alertConfigMap)
	}

	request["AlertConfig"] = alertConfigMaps

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutCreate)), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_cloud_monitor_service_group_monitoring_agent_process", action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	if resp, err := jsonpath.Get("$.Resource", response); err != nil || resp == nil {
		return errs.WrapErrorf(err, errs.IdMsg, "alicloud_cloud_monitor_service_group_monitoring_agent_process")
	} else {
		groupProcessId := resp.(map[string]interface{})["GroupProcessId"]
		d.SetId(fmt.Sprintf("%v:%v", request["GroupId"], groupProcessId))
	}

	return resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cloudMonitorServiceServiceV2 := CloudMonitorServiceServiceV2{client}

	object, err := cloudMonitorServiceServiceV2.DescribeCloudMonitorServiceGroupMonitoringAgentProcess(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_monitor_service_group_monitoring_agent_process DescribeCloudMonitorServiceGroupMonitoringAgentProcess Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("group_id", object["GroupId"])
	d.Set("process_name", object["ProcessName"])
	d.Set("match_express_filter_relation", object["MatchExpressFilterRelation"])
	d.Set("group_monitoring_agent_process_id", object["Id"])

	if matchExpress, ok := object["MatchExpress"]; ok {
		if matchExpressList, ok := matchExpress.(map[string]interface{})["MatchExpress"]; ok {
			matchExpressMaps := make([]map[string]interface{}, 0)
			for _, matchExpresses := range matchExpressList.([]interface{}) {
				matchExpressArg := matchExpresses.(map[string]interface{})
				matchExpressMap := map[string]interface{}{}

				if name, ok := matchExpressArg["Name"]; ok {
					matchExpressMap["name"] = name
				}

				if value, ok := matchExpressArg["Value"]; ok {
					matchExpressMap["value"] = value
				}

				if function, ok := matchExpressArg["Function"]; ok {
					matchExpressMap["function"] = function
				}

				matchExpressMaps = append(matchExpressMaps, matchExpressMap)
			}

			d.Set("match_express", matchExpressMaps)
		}
	}

	if alertConfig, ok := object["AlertConfig"]; ok {
		if alertConfigList, ok := alertConfig.(map[string]interface{})["AlertConfig"]; ok {
			alertConfigMaps := make([]map[string]interface{}, 0)
			for _, alertConfigs := range alertConfigList.([]interface{}) {
				alertConfigArg := alertConfigs.(map[string]interface{})
				alertConfigMap := map[string]interface{}{}

				if escalationsLevel, ok := alertConfigArg["EscalationsLevel"]; ok {
					alertConfigMap["escalations_level"] = escalationsLevel
				}

				if comparisonOperator, ok := alertConfigArg["ComparisonOperator"]; ok {
					alertConfigMap["comparison_operator"] = comparisonOperator
				}

				if statistics, ok := alertConfigArg["Statistics"]; ok {
					alertConfigMap["statistics"] = statistics
				}

				if threshold, ok := alertConfigArg["Threshold"]; ok {
					alertConfigMap["threshold"] = threshold
				}

				if times, ok := alertConfigArg["Times"]; ok {
					alertConfigMap["times"] = times
				}

				if effectiveInterval, ok := alertConfigArg["EffectiveInterval"]; ok {
					alertConfigMap["effective_interval"] = effectiveInterval
				}

				if silenceTime, ok := alertConfigArg["SilenceTime"]; ok {
					alertConfigMap["silence_time"] = silenceTime
				}

				if webhook, ok := alertConfigArg["Webhook"]; ok {
					alertConfigMap["webhook"] = webhook
				}

				if targetList, ok := alertConfigArg["TargetList"]; ok {
					if targets, ok := targetList.(map[string]interface{})["Target"]; ok {
						targetsMaps := make([]map[string]interface{}, 0)
						for _, targetsList := range targets.([]interface{}) {
							targetsArg := targetsList.(map[string]interface{})
							targetsMap := map[string]interface{}{}

							if id, ok := targetsArg["Id"]; ok {
								targetsMap["target_list_id"] = id
							}

							if jsonParmas, ok := targetsArg["JsonParmas"]; ok {
								targetsMap["json_params"] = jsonParmas
							}

							if level, ok := targetsArg["Level"]; ok {
								targetsMap["level"] = level
							}

							if arn, ok := targetsArg["Arn"]; ok {
								targetsMap["arn"] = arn
							}

							targetsMaps = append(targetsMaps, targetsMap)
						}

						alertConfigMap["target_list"] = targetsMaps
					}
				}

				alertConfigMaps = append(alertConfigMaps, alertConfigMap)
			}

			d.Set("alert_config", alertConfigMaps)
		}
	}

	return nil
}

func resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"GroupId": parts[0],
		"Id":      parts[1],
	}

	if d.HasChange("alert_config") {
		update = true
	}
	alertConfig := d.Get("alert_config")
	alertConfigMaps := make([]map[string]interface{}, 0)
	for _, alertConfigList := range alertConfig.([]interface{}) {
		alertConfigMap := map[string]interface{}{}
		alertConfigArg := alertConfigList.(map[string]interface{})

		alertConfigMap["EscalationsLevel"] = alertConfigArg["escalations_level"]
		alertConfigMap["ComparisonOperator"] = alertConfigArg["comparison_operator"]
		alertConfigMap["Statistics"] = alertConfigArg["statistics"]
		alertConfigMap["Threshold"] = alertConfigArg["threshold"]
		alertConfigMap["Times"] = alertConfigArg["times"]

		if effectiveInterval, ok := alertConfigArg["effective_interval"]; ok {
			alertConfigMap["EffectiveInterval"] = effectiveInterval
		}

		if silenceTime, ok := alertConfigArg["silence_time"]; ok {
			alertConfigMap["SilenceTime"] = silenceTime
		}

		if webhook, ok := alertConfigArg["webhook"]; ok {
			alertConfigMap["Webhook"] = webhook
		}

		if targetList, ok := alertConfigArg["target_list"]; ok {
			targetListMaps := make([]map[string]interface{}, 0)
			for _, targetListArgList := range targetList.([]interface{}) {
				targetListMap := map[string]interface{}{}
				targetListArg := targetListArgList.(map[string]interface{})

				if id, ok := targetListArg["target_list_id"]; ok {
					targetListMap["Id"] = id
				}

				if jsonParams, ok := targetListArg["json_params"]; ok {
					targetListMap["JsonParams"] = jsonParams
				}

				if level, ok := targetListArg["level"]; ok {
					targetListMap["Level"] = level
				}

				if arn, ok := targetListArg["arn"]; ok {
					targetListMap["Arn"] = arn
				}

				targetListMaps = append(targetListMaps, targetListMap)
			}

			alertConfigMap["TargetList"] = targetListMaps
		}

		alertConfigMaps = append(alertConfigMaps, alertConfigMap)
	}

	request["AlertConfig"] = alertConfigMaps

	if update {
		action := "ModifyGroupMonitoringAgentProcess"

		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutUpdate)), func() *retry.RetryError {
			response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			return nil
		})
		logs.AddDebug(action, response, request)

		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

		if fmt.Sprint(response["Success"]) == "false" {
			return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
	}

	return resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessRead(d, meta)
}

func resourceAliCloudCloudMonitorServiceGroupMonitoringAgentProcessDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteGroupMonitoringAgentProcess"
	var response map[string]interface{}
	var err error

	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}

	request := map[string]interface{}{
		"GroupId": parts[0],
		"Id":      parts[1],
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(client.GetRetryTimeout(d.Timeout(schema.TimeoutDelete)), func() *retry.RetryError {
		response, err = client.RpcPost("Cms", "2019-01-01", action, nil, request, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		return nil
	})
	logs.AddDebug(action, response, request)

	if err != nil {
		if errs.NotFoundError(err) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	return nil
}
