// Package alikafka. This file is generated automatically. Please do not modify it manually, thank you!
package alikafka

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/convert"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/errs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/helper"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/logs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceScheduledScalingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceScheduledScalingRuleCreate,
		Read:   resourceScheduledScalingRuleRead,
		Update: resourceScheduledScalingRuleUpdate,
		Delete: resourceScheduledScalingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"duration_minutes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"first_scheduled_time": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repeat_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reserved_pub_flow": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"reserved_sub_flow": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"schedule_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"weekly_types": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceScheduledScalingRuleCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)

	action := "CreateScheduledScalingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		request["InstanceId"] = v
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request["RuleName"] = v
	}
	request["RegionId"] = client.RegionId

	request["ReservedPubFlow"] = d.Get("reserved_pub_flow")
	request["ReservedSubFlow"] = d.Get("reserved_sub_flow")
	request["ScheduleType"] = d.Get("schedule_type")
	request["DurationMinutes"] = d.Get("duration_minutes")
	request["TimeZone"] = d.Get("time_zone")
	request["FirstScheduledTime"] = d.Get("first_scheduled_time")
	if v, ok := d.GetOkExists("enable"); ok {
		request["Enable"] = v
	}
	if v, ok := d.GetOk("repeat_type"); ok {
		request["RepeatType"] = v
	}
	if v, ok := d.GetOk("weekly_types"); ok {
		weeklyTypesMapsArray := convert.ConvertToInterfaceArray(v)

		weeklyTypesMapsJson, err := json.Marshal(weeklyTypesMapsArray)
		if err != nil {
			return errs.WrapError(err)
		}
		request["WeeklyTypes"] = string(weeklyTypesMapsJson)
	}

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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
		return errs.WrapErrorf(err, errs.DefaultErrorMsg, "alicloud_alikafka_scheduled_scaling_rule", action, errs.AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprintf("%v:%v", request["InstanceId"], request["RuleName"]))

	return resourceScheduledScalingRuleRead(d, meta)
}

func resourceScheduledScalingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	alikafkaServiceV2 := NewAlikafkaServiceV2(client)

	objectRaw, err := alikafkaServiceV2.DescribeAlikafkaScheduledScalingRule(d.Id())
	if err != nil {
		if !d.IsNewResource() && errs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_alikafka_scheduled_scaling_rule DescribeAlikafkaScheduledScalingRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errs.WrapError(err)
	}

	d.Set("duration_minutes", objectRaw["DurationMinutes"])
	d.Set("enable", objectRaw["Enable"])
	d.Set("first_scheduled_time", objectRaw["FirstScheduledTime"])
	d.Set("repeat_type", objectRaw["RepeatType"])
	d.Set("reserved_pub_flow", objectRaw["ReservedPubFlow"])
	d.Set("reserved_sub_flow", objectRaw["ReservedSubFlow"])
	d.Set("schedule_type", objectRaw["ScheduleType"])
	d.Set("time_zone", objectRaw["TimeZone"])
	d.Set("rule_name", objectRaw["RuleName"])

	weeklyTypesRaw, _ := jsonpath.Get("$.WeeklyTypes.WeeklyTypes", objectRaw)
	d.Set("weekly_types", weeklyTypesRaw)

	parts := strings.Split(d.Id(), ":")
	d.Set("instance_id", parts[0])

	return nil
}

func resourceScheduledScalingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var request map[string]interface{}
	var response map[string]interface{}
	var query map[string]interface{}
	update := false

	var err error
	parts := strings.Split(d.Id(), ":")
	action := "ModifyScheduledScalingRule"
	request = make(map[string]interface{})
	query = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleName"] = parts[1]
	request["RegionId"] = client.RegionId
	if d.HasChange("enable") {
		update = true

		if v, ok := d.GetOkExists("enable"); ok {
			request["Enable"] = v
		}
	}
	if update {
		wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
		err = retry.Retry(d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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
	}

	return resourceScheduledScalingRuleRead(d, meta)
}

func resourceScheduledScalingRuleDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AliyunClient)
	parts := strings.Split(d.Id(), ":")
	action := "DeleteScheduledScalingRule"
	var request map[string]interface{}
	var response map[string]interface{}
	query := make(map[string]interface{})
	var err error
	request = make(map[string]interface{})
	request["InstanceId"] = parts[0]
	request["RuleName"] = parts[1]
	request["RegionId"] = client.RegionId

	wait := helper.IncrementalWait(3*time.Second, 5*time.Second)
	err = retry.Retry(d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		response, err = client.RpcPost("alikafka", "2019-09-16", action, query, request, true)
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

	return nil
}
