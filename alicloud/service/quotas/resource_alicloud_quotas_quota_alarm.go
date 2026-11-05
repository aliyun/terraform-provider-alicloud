// Package quotas. This file is generated automatically. Please do not modify it manually, thank you!
package quotas

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/blues/jsonata-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceQuotaAlarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceQuotaAlarmCreate,
		Read:   resourceQuotaAlarmRead,
		Update: resourceQuotaAlarmUpdate,
		Delete: resourceQuotaAlarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_alarm_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"quota_dimensions": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			"threshold_percent": {
				Type:         schema.TypeFloat,
				ExactlyOneOf: []string{"threshold_percent", "threshold"},
				Optional:     true,
			},
			"threshold_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"web_hook": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceQuotaAlarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "CreateQuotaAlarm"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["ProductCode"] = d.Get("product_code")
	request["QuotaActionCode"] = d.Get("quota_action_code")
	if v, ok := d.GetOk("threshold"); ok {
		request["Threshold"] = v
	}
	if v, ok := d.GetOk("threshold_percent"); ok {
		request["ThresholdPercent"] = v
	}
	if v, ok := d.GetOk("web_hook"); ok {
		request["WebHook"] = v
	}
	if v, ok := d.GetOk("quota_dimensions"); ok {
		quotaDimensionsMaps := make([]map[string]interface{}, 0)
		for _, dataLoop := range v.(*schema.Set).List() {
			dataLoopTmp := dataLoop.(map[string]interface{})
			dataLoopMap := make(map[string]interface{})
			dataLoopMap["Key"] = dataLoopTmp["key"]
			dataLoopMap["Value"] = dataLoopTmp["value"]
			quotaDimensionsMaps = append(quotaDimensionsMaps, dataLoopMap)
		}
		request["QuotaDimensions"] = quotaDimensionsMaps
	}

	request["AlarmName"] = d.Get("quota_alarm_name")
	if v, ok := d.GetOk("threshold_type"); ok {
		request["ThresholdType"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)
		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_quotas_quota_alarm", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["AlarmId"]))

	return resourceQuotaAlarmUpdate(d, meta)
}

func resourceQuotaAlarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	quotasServiceV2 := QuotasServiceV2{client}

	objectRaw, err := quotasServiceV2.DescribeQuotasQuotaAlarm(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_quotas_quota_alarm DescribeQuotasQuotaAlarm Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("create_time", objectRaw["CreateTime"])
	d.Set("product_code", objectRaw["ProductCode"])
	d.Set("quota_action_code", objectRaw["QuotaActionCode"])
	d.Set("quota_alarm_name", objectRaw["AlarmName"])
	d.Set("threshold", objectRaw["Threshold"])
	d.Set("threshold_percent", objectRaw["ThresholdPercent"])
	d.Set("threshold_type", objectRaw["ThresholdType"])
	d.Set("web_hook", objectRaw["Webhook"])

	e := jsonata.MustCompile("$each($.QuotaDimension, function($v, $k) {{\"value\":$v, \"key\": $k}})[]")
	evaluation, _ := e.Eval(objectRaw)
	d.Set("quota_dimensions", evaluation)

	return nil
}

func resourceQuotaAlarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	update := false
	action := "UpdateQuotaAlarm"
	var err error
	request = make(map[string]interface{})

	request["AlarmId"] = d.Id()
	if !d.IsNewResource() && d.HasChange("threshold") {
		update = true
		request["Threshold"] = d.Get("threshold")
	}

	if !d.IsNewResource() && d.HasChange("threshold_percent") {
		update = true
		request["ThresholdPercent"] = d.Get("threshold_percent")
	}

	if !d.IsNewResource() && d.HasChange("web_hook") {
		update = true
		request["WebHook"] = d.Get("web_hook")
	}

	if !d.IsNewResource() && d.HasChange("quota_alarm_name") {
		update = true
	}
	request["AlarmName"] = d.Get("quota_alarm_name")
	if !d.IsNewResource() && d.HasChange("threshold_type") {
		update = true
		request["ThresholdType"] = d.Get("threshold_type")
	}

	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)

			if err != nil {
				if helper.NeedRetry(err) {
					wait()
					return retry.RetryableError(err)
				}
				return retry.NonRetryableError(err)
			}
			logs.AddDebug(action, response, request)
			return nil
		})
		if err != nil {
			return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
		}

	}
	return resourceQuotaAlarmRead(d, meta)
}

func resourceQuotaAlarmDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "DeleteQuotaAlarm"
	var request map[string]interface{}
	var response map[string]interface{}
	var err error
	request = make(map[string]interface{})

	request["AlarmId"] = d.Id()

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.Do("quotas", helper.RpcParam("POST", "2020-05-10", action), nil, request, nil, nil, false)

		if err != nil {
			if helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})

	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}

	return nil
}
