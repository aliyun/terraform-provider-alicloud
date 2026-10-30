package dts

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceJobMonitorRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobMonitorRuleCreate,
		Read:   resourceJobMonitorRuleRead,
		Update: resourceJobMonitorRuleUpdate,
		Delete: resourceJobMonitorRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"delay_rule_time": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"dts_job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"phone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Y", "N"}, false),
			},
			"type": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"delay", "error"}, false),
			},
		},
	}
}

func resourceJobMonitorRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateJobMonitorRule"
	request := make(map[string]interface{})
	var err error
	if v, ok := d.GetOk("delay_rule_time"); ok {
		request["DelayRuleTime"] = v
	}
	request["DtsJobId"] = d.Get("dts_job_id")
	if v, ok := d.GetOk("phone"); ok {
		request["Phone"] = v
	}
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("state"); ok {
		request["State"] = v
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_dts_job_monitor_rule", action, errs.AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return errs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	d.SetId(fmt.Sprint(request["DtsJobId"]))

	return resourceJobMonitorRuleRead(d, meta)
}
func resourceJobMonitorRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dtsService := NewService(client)
	object, err := dtsService.DescribeDtsJobMonitorRule(d.Id())
	if err != nil {
		if errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dts_job_monitor_rule dtsService.DescribeDtsJobMonitorRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("dts_job_id", d.Id())
	if v := object["MonitorRules"].([]interface{}); len(v) > 0 {
		d.Set("delay_rule_time", v[0].(map[string]interface{})["DelayRuleTime"])
		d.Set("phone", v[0].(map[string]interface{})["Phone"])
		d.Set("state", v[0].(map[string]interface{})["State"])
		d.Set("type", v[0].(map[string]interface{})["Type"])
	}
	return nil
}
func resourceJobMonitorRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var err error
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"DtsJobId": d.Id(),
	}
	if d.HasChange("delay_rule_time") {
		update = true
		if v, ok := d.GetOk("delay_rule_time"); ok {
			request["DelayRuleTime"] = v
		}
	}
	if d.HasChange("phone") {
		update = true
		if v, ok := d.GetOk("phone"); ok {
			request["Phone"] = v
		}
	}
	request["RegionId"] = client.RegionId
	if d.HasChange("state") {
		update = true
		if v, ok := d.GetOk("state"); ok {
			request["State"] = v
		}
	}
	if v, ok := d.GetOk("type"); ok {
		request["Type"] = v
	}
	if update {
		action := "CreateJobMonitorRule"
		wait := helper.IncrementalWait(3*time.Second, 3*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("Dts", "2020-01-01", action, nil, request, false)
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
	return resourceJobMonitorRuleRead(d, meta)
}
func resourceJobMonitorRuleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudDtsJobMonitorRule. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
