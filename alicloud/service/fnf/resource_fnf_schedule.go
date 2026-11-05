package fnf

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceFnfSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceFnfScheduleCreate,
		Read:   resourceFnfScheduleRead,
		Update: resourceFnfScheduleUpdate,
		Delete: resourceFnfScheduleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cron_expression": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"flow_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payload": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFnfScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSchedule"
	request := make(map[string]interface{})
	var err error
	request["CronExpression"] = d.Get("cron_expression")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}

	if v, ok := d.GetOkExists("enable"); ok {
		request["Enable"] = v
	}

	request["FlowName"] = d.Get("flow_name")
	if v, ok := d.GetOk("payload"); ok {
		request["Payload"] = v
	}

	request["ScheduleName"] = d.Get("schedule_name")
	wait := helper.IncrementalWait(3*time.Second, 1*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("fnf", "2019-03-15", action, nil, request, false)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_fnf_schedule", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["FlowName"], ":", response["ScheduleName"]))

	return resourceFnfScheduleRead(d, meta)
}
func resourceFnfScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	fnfService := NewService(client)
	object, err := fnfService.DescribeFnfSchedule(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_fnf_schedule fnfService.DescribeFnfSchedule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	d.Set("flow_name", parts[0])
	d.Set("schedule_name", parts[1])
	d.Set("cron_expression", object["CronExpression"])
	d.Set("description", object["Description"])
	d.Set("enable", object["Enable"])
	d.Set("last_modified_time", object["LastModifiedTime"])
	d.Set("payload", object["Payload"])
	d.Set("schedule_id", object["ScheduleId"])
	return nil
}
func resourceFnfScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"FlowName":     parts[0],
		"ScheduleName": parts[1],
	}
	if d.HasChange("cron_expression") {
		update = true
		request["CronExpression"] = d.Get("cron_expression")
	}
	if d.HasChange("description") {
		update = true
		request["Description"] = d.Get("description")
	}
	if d.HasChange("enable") {
		update = true
		request["Enable"] = d.Get("enable")
	}
	if d.HasChange("payload") {
		update = true
		request["Payload"] = d.Get("payload")
	}
	if update {
		action := "UpdateSchedule"
		wait := helper.IncrementalWait(3*time.Second, 1*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("fnf", "2019-03-15", action, nil, request, false)
			if err != nil {
				if errs.IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || helper.NeedRetry(err) {
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
	return resourceFnfScheduleRead(d, meta)
}
func resourceFnfScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := helper.ParseResourceId(d.Id(), 2)
	if err != nil {
		return errs.WrapError(err)
	}
	action := "DeleteSchedule"
	var response map[string]interface{}
	request := map[string]interface{}{
		"FlowName":     parts[0],
		"ScheduleName": parts[1],
	}

	wait := helper.IncrementalWait(3*time.Second, 1*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcGet("fnf", "2019-03-15", action, request, nil)
		if err != nil {
			if errs.IsExpectedErrors(err, []string{"ConcurrentUpdateError", "InternalServerError"}) || helper.NeedRetry(err) {
				wait()
				return retry.RetryableError(err)
			}
			return retry.NonRetryableError(err)
		}
		logs.AddDebug(action, response, request)
		return nil
	})
	if err != nil {
		if errs.IsExpectedErrors(err, []string{"ExecutionNotExists", "FlowNotExists", "ScheduleNotExists"}) {
			return nil
		}
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, d.Id(), action, errs.AlibabaCloudSdkGoERROR)
	}
	return nil
}
